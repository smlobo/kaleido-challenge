// Package keystore holds the service's signing keys in Postgres, encrypted
// under a key-encryption key that only ever comes from the environment, and
// hands out nonces for them.
package keystore

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/smlobo/kaleido-challenge/internal/chain"
)

// Purpose is what a key is for; see the design doc's key table.
type Purpose string

const (
	Deployer  Purpose = "deployer"
	BoxOffice Purpose = "box_office"
	Payout    Purpose = "payout"
	Gate      Purpose = "gate"
	Wallet    Purpose = "wallet"
)

// DBTX is satisfied by *pgxpool.Pool, *pgx.Conn and pgx.Tx, so callers choose
// whether an operation joins their transaction.
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Key is a stored key's public data. The private key never leaves the store
// except inside a *chain.Signer.
type Key struct {
	ID        string
	Address   common.Address
	Purpose   Purpose
	NextNonce uint64
}

// ErrNotFound is returned when no key matches.
var ErrNotFound = errors.New("keystore: key not found")

// Store creates, loads and decrypts keys.
type Store struct {
	sealer  *sealer
	chainID *big.Int
}

// New returns a Store that encrypts with kek (32 bytes, see ParseKEK) and
// builds signers for chainID.
func New(kek []byte, chainID *big.Int) (*Store, error) {
	s, err := newSealer(kek)
	if err != nil {
		return nil, err
	}
	return &Store{sealer: s, chainID: new(big.Int).Set(chainID)}, nil
}

// Create generates a new key, encrypts it and inserts it. Pass a pgx.Tx to
// create it atomically with the row that references it (an organizer, a user).
func (s *Store) Create(ctx context.Context, db DBTX, purpose Purpose) (Key, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return Key{}, err
	}
	addr := crypto.PubkeyToAddress(priv.PublicKey)
	sealed, err := s.sealer.seal(crypto.FromECDSA(priv), hexAddr(addr))
	if err != nil {
		return Key{}, err
	}
	k := Key{Address: addr, Purpose: purpose}
	err = db.QueryRow(ctx,
		`INSERT INTO signing_keys (address, purpose, encrypted_private_key)
		 VALUES ($1, $2, $3) RETURNING id`,
		hexAddr(addr), string(purpose), sealed).Scan(&k.ID)
	if err != nil {
		return Key{}, fmt.Errorf("insert signing key: %w", err)
	}
	return k, nil
}

// EnsureDeployer returns the single deployer key, creating it if none exists.
// Safe when several instances start at once: the partial unique index lets
// exactly one insert win, and the others read the winner's row.
func (s *Store) EnsureDeployer(ctx context.Context, db DBTX) (Key, error) {
	k, err := s.byPurpose(ctx, db, Deployer)
	if err == nil || !errors.Is(err, ErrNotFound) {
		return k, err
	}
	priv, err := crypto.GenerateKey()
	if err != nil {
		return Key{}, err
	}
	addr := crypto.PubkeyToAddress(priv.PublicKey)
	sealed, err := s.sealer.seal(crypto.FromECDSA(priv), hexAddr(addr))
	if err != nil {
		return Key{}, err
	}
	if _, err := db.Exec(ctx,
		`INSERT INTO signing_keys (address, purpose, encrypted_private_key)
		 VALUES ($1, 'deployer', $2)
		 ON CONFLICT (purpose) WHERE purpose = 'deployer' DO NOTHING`,
		hexAddr(addr), sealed); err != nil {
		return Key{}, fmt.Errorf("insert deployer key: %w", err)
	}
	return s.byPurpose(ctx, db, Deployer)
}

// Get loads a key's public data by ID.
func (s *Store) Get(ctx context.Context, db DBTX, id string) (Key, error) {
	return scanKey(db.QueryRow(ctx,
		`SELECT id, address, purpose, next_nonce FROM signing_keys WHERE id = $1`, id))
}

func (s *Store) byPurpose(ctx context.Context, db DBTX, p Purpose) (Key, error) {
	return scanKey(db.QueryRow(ctx,
		`SELECT id, address, purpose, next_nonce FROM signing_keys WHERE purpose = $1
		 ORDER BY created_at LIMIT 1`, string(p)))
}

// Signer decrypts a key and returns a signer for it. Decrypt as late as
// possible and do not keep the signer around longer than the request.
func (s *Store) Signer(ctx context.Context, db DBTX, id string) (*chain.Signer, error) {
	var addr string
	var sealed []byte
	err := db.QueryRow(ctx,
		`SELECT address, encrypted_private_key FROM signing_keys WHERE id = $1`, id).Scan(&addr, &sealed)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	raw, err := s.sealer.open(sealed, addr)
	if err != nil {
		return nil, err
	}
	priv, err := crypto.ToECDSA(raw)
	clear(raw)
	if err != nil {
		return nil, fmt.Errorf("keystore: stored key is invalid: %w", err)
	}
	signer := chain.NewSigner(priv, s.chainID)
	if !strings.EqualFold(signer.Address().Hex(), addr) {
		return nil, fmt.Errorf("keystore: key for %s decrypts to %s", addr, signer.Address().Hex())
	}
	return signer, nil
}

// ReserveNonce takes the next nonce for a key and advances the counter.
//
// It must run inside the database transaction that also stores the signed
// transaction (the outbox row). The UPDATE locks the key's row until that
// transaction ends, so concurrent callers on any instance get distinct,
// consecutive nonces; and if the transaction rolls back, the nonce is
// returned with it, so no gap is left behind.
func ReserveNonce(ctx context.Context, tx pgx.Tx, keyID string) (uint64, error) {
	var nonce int64
	err := tx.QueryRow(ctx,
		`UPDATE signing_keys SET next_nonce = next_nonce + 1
		 WHERE id = $1 RETURNING next_nonce - 1`, keyID).Scan(&nonce)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("reserve nonce: %w", err)
	}
	return uint64(nonce), nil
}

func scanKey(row pgx.Row) (Key, error) {
	var k Key
	var addr, purpose string
	var next int64
	err := row.Scan(&k.ID, &addr, &purpose, &next)
	if errors.Is(err, pgx.ErrNoRows) {
		return Key{}, ErrNotFound
	}
	if err != nil {
		return Key{}, err
	}
	k.Address = common.HexToAddress(addr)
	k.Purpose = Purpose(purpose)
	k.NextNonce = uint64(next)
	return k, nil
}

// hexAddr is the lowercase form the schema's eth_address domain requires.
func hexAddr(a common.Address) string { return strings.ToLower(a.Hex()) }
