package outbox

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
)

// These tests run the outbox against the real Besu node from
// docker/docker-compose.yaml, and skip if it is down.

func dialBesu(t *testing.T) *chain.Client {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	c, err := chain.Dial(ctx, "http://localhost:8545", chain.ChainID1337)
	if err != nil {
		t.Skipf("Besu not reachable (%v); run: docker compose -f docker/docker-compose.yaml up -d --wait besu postgres", err)
	}
	if _, err := c.PendingNonce(ctx, common.Address{}); err != nil {
		c.Close()
		t.Skipf("Besu not answering: %v", err)
	}
	t.Cleanup(c.Close)
	return c
}

// besuFixture is a fixture whose outbox talks to Besu. With run set, its
// loops run in the background until the test ends.
func besuFixture(t *testing.T, run bool) (*fixture, *chain.Client, context.Context) {
	t.Helper()
	c := dialBesu(t)
	f := newFixture(t, Config{PollInterval: 200 * time.Millisecond, ResendAfter: 3 * time.Second})
	f.ob = New(f.pool, f.keys, c, f.ob.cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	t.Cleanup(cancel)
	if run {
		done := make(chan struct{})
		go func() { f.ob.Run(ctx); close(done) }()
		t.Cleanup(func() { cancel(); <-done })
	}
	return f, c, ctx
}

func (f *fixture) newKey(t *testing.T, p keystore.Purpose) keystore.Key {
	t.Helper()
	k, err := f.keys.Create(context.Background(), f.pool, p)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func (f *fixture) wait(t *testing.T, ctx context.Context, id string) Op {
	t.Helper()
	op, err := f.ob.Wait(ctx, id)
	if err != nil {
		t.Fatalf("wait for %s: %v (status %s, last error %q)", id, err, op.Status, op.LastError)
	}
	return op
}

// Deploy a TicketNFT, then mint five tickets that are committed but never
// sent (as if the instance crashed after commit), plus one duplicate that
// reverts. The submitter sends them all and the receipt loop records them.
func TestBesuDeployMintAndRevertThroughLoops(t *testing.T) {
	f, c, ctx := besuFixture(t, true)
	nft := ticketnft.NewTicketNFT()
	deployer := f.newKey(t, keystore.Deployer)
	boxOffice := f.newKey(t, keystore.BoxOffice)
	gate := f.newKey(t, keystore.Gate)

	code := append(common.FromHex(ticketnft.TicketNFTMetaData.Bin), nft.PackConstructor("Outbox Fest", "OBX", 100, 11_000,
		uint64(time.Now().Add(time.Hour).Unix()), "", boxOffice.Address, gate.Address)...)
	gas, err := c.EstimateGas(ctx, deployer.Address, nil, nil, code)
	if err != nil {
		t.Fatal(err)
	}
	deploy := f.commit(t, Request{KeyID: deployer.ID, Data: code, Gas: gas})
	if err := f.ob.Send(ctx, deploy.ID); err != nil {
		t.Fatal(err)
	}
	deployed := f.wait(t, ctx, deploy.ID)
	if deployed.Status != Mined || deployed.ContractAddress == nil {
		t.Fatalf("deploy: %+v", deployed)
	}
	addr := *deployed.ContractAddress

	holder := common.HexToAddress("0x000000000000000000000000000000000000abcd")
	var mints []Op
	for id := int64(1); id <= 5; id++ {
		mints = append(mints, f.commit(t, Request{KeyID: boxOffice.ID, To: &addr, Gas: 200_000,
			Data: nft.PackMintTicket(holder, big.NewInt(id), big.NewInt(100_000_000), common.Hash{1})}))
	}
	// Token 3 again, with a fixed gas limit: estimating would catch the revert.
	dup := f.commit(t, Request{KeyID: boxOffice.ID, To: &addr, Gas: 200_000,
		Data: nft.PackMintTicket(holder, big.NewInt(3), big.NewInt(100_000_000), common.Hash{1})})

	for _, m := range mints {
		if got := f.wait(t, ctx, m.ID); got.Status != Mined {
			t.Fatalf("mint nonce %d: %s %q", m.Nonce, got.Status, got.LastError)
		}
	}
	reverted := f.wait(t, ctx, dup.ID)
	if reverted.Status != Reverted || reverted.RevertReason == "" {
		t.Fatalf("duplicate mint: %+v", reverted)
	}
	decoded, err := nft.UnpackError(common.FromHex(reverted.RevertReason))
	if err != nil {
		t.Fatalf("decode revert reason %s: %v", reverted.RevertReason, err)
	}
	if e, ok := decoded.(*ticketnft.TicketNFTERC721InvalidSender); !ok || e.Sender != (common.Address{}) {
		t.Fatalf("revert reason = %#v, want ERC721InvalidSender(0x0)", decoded)
	}

	owner, err := bind.Call(nft.Instance(c.Eth(), addr), &bind.CallOpts{Context: ctx},
		nft.PackOwnerOf(big.NewInt(5)), nft.UnpackOwnerOf)
	if err != nil || owner != holder {
		t.Fatalf("ownerOf(5) = %s, %v", owner, err)
	}
	if n, _ := c.PendingNonce(ctx, boxOffice.Address); n != 6 {
		t.Fatalf("box office chain nonce = %d, want 6 (five mints and the revert)", n)
	}
}

// Besu rejects a gas limit below the intrinsic cost. The row fails, a
// transfer to self takes its nonce, and the next transaction still mines.
func TestBesuRejectedTransactionFillsNonceGap(t *testing.T) {
	f, c, ctx := besuFixture(t, true)
	bad := f.enqueue(t, Request{Gas: 1})
	next := f.enqueue(t, Request{Gas: 21_000})

	if got := f.wait(t, ctx, bad.ID); got.Status != Failed || got.LastError == "" {
		t.Fatalf("bad: %s %q", got.Status, got.LastError)
	}
	if got := f.wait(t, ctx, next.ID); got.Status != Mined {
		t.Fatalf("next: %s %q", got.Status, got.LastError)
	}
	var filler Op
	err := pgx.BeginFunc(ctx, f.pool, func(tx pgx.Tx) error {
		var id string
		if err := tx.QueryRow(ctx, `SELECT id FROM transactions WHERE nonce = 0 AND status <> 'failed'`).Scan(&id); err != nil {
			return err
		}
		var err error
		filler, err = Get(ctx, tx, id)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if filler.Status != Mined || filler.To == nil || *filler.To != f.key.Address {
		t.Fatalf("filler: %+v", filler)
	}
	if n, _ := c.PendingNonce(ctx, f.key.Address); n != 2 {
		t.Fatalf("chain nonce = %d, want 2", n)
	}
}

// The same signed bytes sent again, before and after mining, are mined once.
// The loops are not running, so each step is the test's own.
func TestBesuResendingIsSafe(t *testing.T) {
	f, c, ctx := besuFixture(t, false)
	op := f.enqueue(t, Request{Gas: 21_000})
	if err := f.ob.Send(ctx, op.ID); err != nil {
		t.Fatal(err)
	}
	// Pretend the first send was lost: back to 'signed' and send again.
	if _, err := f.pool.Exec(ctx, `UPDATE transactions SET status = 'signed' WHERE id = $1`, op.ID); err != nil {
		t.Fatal(err)
	}
	if err := f.ob.Send(ctx, op.ID); err != nil {
		t.Fatal(err)
	}
	if got := f.get(t, op.ID); got.Status != Submitted && got.Status != Mined {
		t.Fatalf("after resend: %s %q", got.Status, got.LastError)
	}
	var mined Op
	for range 30 {
		f.pass(t, f.ob.receiptPass)
		if mined = f.get(t, op.ID); mined.Status.Final() {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if mined.Status != Mined {
		t.Fatalf("got %s", mined.Status)
	}
	// Resend after mining: Besu says "Nonce too low", and the receipt shows
	// the bytes are ours, so the row stays mined.
	if _, err := f.pool.Exec(ctx, `UPDATE transactions SET status = 'signed' WHERE id = $1`, op.ID); err != nil {
		t.Fatal(err)
	}
	if err := f.ob.Send(ctx, op.ID); err != nil {
		t.Fatal(err)
	}
	if got := f.get(t, op.ID); got.Status != Mined || *got.BlockNumber != *mined.BlockNumber {
		t.Fatalf("after resend of mined tx: %s at block %v", got.Status, got.BlockNumber)
	}
	if n, _ := c.PendingNonce(ctx, f.key.Address); n != 1 {
		t.Fatalf("chain nonce = %d, want 1", n)
	}
}

// A key used outside the outbox takes the nonce the database hands out next.
// The row fails with a clear error instead of retrying forever.
func TestBesuNonceUsedOutsideOutboxFails(t *testing.T) {
	f, c, ctx := besuFixture(t, true)
	signer, err := f.keys.Signer(ctx, f.pool, f.key.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := chain.SendAndWait(ctx, c, signer, &target, nil); err != nil {
		t.Fatal(err)
	}
	op := f.enqueue(t, Request{Gas: 21_000})
	got := f.wait(t, ctx, op.ID)
	if got.Status != Failed || got.LastError == "" {
		t.Fatalf("got %s %q", got.Status, got.LastError)
	}
	t.Logf("last_error: %s", got.LastError)
}
