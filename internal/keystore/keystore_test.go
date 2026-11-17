package keystore_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"sort"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/dbtest"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
)

func randomKEK(t *testing.T) []byte {
	kek := make([]byte, keystore.KEKSize)
	_, _ = rand.Read(kek)
	return kek
}

func newStore(t *testing.T) *keystore.Store {
	t.Helper()
	kek := randomKEK(t)
	s, err := keystore.New(kek, chain.ChainID1337)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestCreateAndSign(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	ks := newStore(t)

	k, err := ks.Create(ctx, pool, keystore.BoxOffice)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ks.Signer(ctx, pool, k.ID)
	if err != nil {
		t.Fatal(err)
	}
	if signer.Address() != k.Address {
		t.Fatalf("signer address %s, key address %s", signer.Address(), k.Address)
	}
	got, err := ks.Get(ctx, pool, k.ID)
	if err != nil || got.Purpose != keystore.BoxOffice || got.NextNonce != 0 {
		t.Fatalf("Get = %+v, %v", got, err)
	}

	// What is stored is not the raw key.
	var sealed []byte
	if err := pool.QueryRow(ctx, "SELECT encrypted_private_key FROM signing_keys WHERE id = $1", k.ID).Scan(&sealed); err != nil {
		t.Fatal(err)
	}
	if len(sealed) != 1+12+32+16 {
		t.Errorf("sealed length %d, want 61 (version + nonce + key + tag)", len(sealed))
	}
}

func TestSignerFailsWithWrongKEK(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	k, err := newStore(t).Create(ctx, pool, keystore.Wallet)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := newStore(t).Signer(ctx, pool, k.ID); err == nil {
		t.Fatal("decrypted with a different KEY_ENCRYPTION_KEY")
	}
}

func TestSignerRejectsSwappedCiphertext(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	ks := newStore(t)
	a, _ := ks.Create(ctx, pool, keystore.Wallet)
	b, _ := ks.Create(ctx, pool, keystore.Wallet)

	// Someone with database access copies b's ciphertext onto a's row.
	if _, err := pool.Exec(ctx, `UPDATE signing_keys SET encrypted_private_key =
		(SELECT encrypted_private_key FROM signing_keys WHERE id = $2) WHERE id = $1`, a.ID, b.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := ks.Signer(ctx, pool, a.ID); err == nil {
		t.Fatal("decrypted a ciphertext that belongs to another address")
	}
}

func TestUnknownKey(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	ks := newStore(t)
	const missing = "00000000-0000-0000-0000-000000000000"
	if _, err := ks.Signer(ctx, pool, missing); !errors.Is(err, keystore.ErrNotFound) {
		t.Errorf("Signer: %v, want ErrNotFound", err)
	}
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		_, err := keystore.ReserveNonce(ctx, tx, missing)
		return err
	})
	if !errors.Is(err, keystore.ErrNotFound) {
		t.Errorf("ReserveNonce: %v, want ErrNotFound", err)
	}
}

func TestEnsureDeployerConcurrent(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	kek := randomKEK(t)

	const n = 10 // instances starting at the same moment
	addrs := make([]common.Address, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			ks, _ := keystore.New(kek, chain.ChainID1337)
			k, err := ks.EnsureDeployer(ctx, pool)
			if err != nil {
				t.Error(err)
				return
			}
			addrs[i] = k.Address
		})
	}
	wg.Wait()

	for _, a := range addrs {
		if a != addrs[0] {
			t.Fatalf("instances got different deployers: %v", addrs)
		}
	}
	var count int
	_ = pool.QueryRow(ctx, "SELECT count(*) FROM signing_keys WHERE purpose = 'deployer'").Scan(&count)
	if count != 1 {
		t.Fatalf("%d deployer keys, want 1", count)
	}
}

// The design doc's concurrency test: many goroutines (standing in for many
// instances) reserve nonces on one key; every nonce is distinct, there are no
// gaps, and the counter ends at the total.
func TestReserveNonceConcurrent(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	k, err := newStore(t).Create(ctx, pool, keystore.BoxOffice)
	if err != nil {
		t.Fatal(err)
	}

	const n = 50
	got := make([]uint64, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
				nonce, err := keystore.ReserveNonce(ctx, tx, k.ID)
				got[i] = nonce
				return err
			})
			if err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()

	sort.Slice(got, func(a, b int) bool { return got[a] < got[b] })
	for i, nonce := range got {
		if nonce != uint64(i) {
			t.Fatalf("sorted nonces %v: position %d is %d (duplicate or gap)", got, i, nonce)
		}
	}
	if next := nextNonce(t, pool, k.ID); next != n {
		t.Fatalf("next_nonce = %d, want %d", next, n)
	}
}

// A rolled-back transaction gives its nonce back, so a failed request never
// leaves a gap that would block every later transaction from the key.
func TestReserveNonceRollback(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	k, _ := newStore(t).Create(ctx, pool, keystore.Gate)

	errBoom := errors.New("signing failed")
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		if n, err := keystore.ReserveNonce(ctx, tx, k.ID); err != nil || n != 0 {
			t.Fatalf("first reserve = %d, %v", n, err)
		}
		return errBoom // e.g. signing or the outbox insert failed
	})
	if !errors.Is(err, errBoom) {
		t.Fatal(err)
	}

	var n uint64
	_ = pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		n, err = keystore.ReserveNonce(ctx, tx, k.ID)
		return err
	})
	if err != nil || n != 0 {
		t.Fatalf("after rollback reserve = %d, %v; want 0", n, err)
	}
}

// Two nonces for one key in the same outbox would violate the unique
// (from_address, nonce) safety net.
func TestOutboxRejectsDuplicateNonce(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	insert := func(hash byte) error {
		_, err := pool.Exec(ctx, `INSERT INTO transactions (from_address, nonce, calldata, raw_tx, tx_hash)
			VALUES ('0x00000000000000000000000000000000000000aa', 7, '\x', '\x', $1)`,
			"0x"+string(bytes.Repeat([]byte{hash}, 64)))
		return err
	}
	if err := insert('a'); err != nil {
		t.Fatal(err)
	}
	if err := insert('b'); err == nil {
		t.Fatal("second transaction with the same (from_address, nonce) was accepted")
	}
}

func nextNonce(t *testing.T, pool *pgxpool.Pool, id string) uint64 {
	t.Helper()
	var n int64
	if err := pool.QueryRow(context.Background(), "SELECT next_nonce FROM signing_keys WHERE id = $1", id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return uint64(n)
}
