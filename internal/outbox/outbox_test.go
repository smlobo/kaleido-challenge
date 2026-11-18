package outbox

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/dbtest"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
)

// fakeNode records every send and answers with scripted errors and receipts.
type fakeNode struct {
	mu       sync.Mutex
	sends    map[common.Hash]int
	sendErr  func(tx *types.Transaction) error
	receipts map[common.Hash]*types.Receipt
	reasons  map[common.Hash]string
}

func newFakeNode() *fakeNode {
	return &fakeNode{sends: map[common.Hash]int{}, receipts: map[common.Hash]*types.Receipt{}, reasons: map[common.Hash]string{}}
}

func (f *fakeNode) SendRaw(_ context.Context, tx *types.Transaction) (common.Hash, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sends[tx.Hash()]++
	if f.sendErr != nil {
		if err := f.sendErr(tx); err != nil {
			return common.Hash{}, err
		}
	}
	return tx.Hash(), nil
}

func (f *fakeNode) Receipt(_ context.Context, h common.Hash) (*types.Receipt, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.receipts[h], f.reasons[h], nil
}

func (f *fakeNode) mine(h common.Hash, status uint64, block int64, contract common.Address, reason string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.receipts[h] = &types.Receipt{Status: status, BlockNumber: big.NewInt(block), ContractAddress: contract}
	f.reasons[h] = reason
}

func (f *fakeNode) sendCount(h common.Hash) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.sends[h]
}

// nodeError is a JSON-RPC error answer, as go-ethereum's rpc client returns it.
type nodeError struct {
	code int
	msg  string
}

func (e nodeError) Error() string  { return e.msg }
func (e nodeError) ErrorCode() int { return e.code }

// Besu's answers, as seen from eth_sendRawTransaction.
var (
	errKnown      = nodeError{-32000, "Known transaction"}
	errNonceLow   = nodeError{-32001, "Nonce too low"}
	errIntrinsic  = nodeError{-32003, "Intrinsic gas exceeds gas limit"}
	errConnRefuse = &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}
)

type fixture struct {
	pool *pgxpool.Pool
	keys *keystore.Store
	node *fakeNode
	ob   *Outbox
	key  keystore.Key
}

func newFixture(t *testing.T, cfg Config) *fixture {
	t.Helper()
	pool := dbtest.Pool(t)
	kek := make([]byte, 32)
	_, _ = rand.Read(kek)
	ks, err := keystore.New(kek, chain.ChainID1337)
	if err != nil {
		t.Fatal(err)
	}
	key, err := ks.Create(context.Background(), pool, keystore.Wallet)
	if err != nil {
		t.Fatal(err)
	}
	node := newFakeNode()
	return &fixture{pool: pool, keys: ks, node: node, ob: New(pool, ks, node, cfg), key: key}
}

var target = common.HexToAddress("0x000000000000000000000000000000000000c001")

// enqueue signs and commits one call to target, without sending it.
func (f *fixture) enqueue(t *testing.T, req Request) Op {
	t.Helper()
	if req.To == nil {
		req.To = &target
	}
	return f.commit(t, req)
}

// enqueueDeploy signs and commits a contract creation.
func (f *fixture) enqueueDeploy(t *testing.T) Op {
	t.Helper()
	return f.commit(t, Request{Data: []byte{0x60, 0x00}})
}

func (f *fixture) commit(t *testing.T, req Request) Op {
	t.Helper()
	if req.KeyID == "" {
		req.KeyID = f.key.ID
	}
	if req.Gas == 0 {
		req.Gas = 50_000
	}
	var op Op
	err := pgx.BeginFunc(context.Background(), f.pool, func(tx pgx.Tx) error {
		var err error
		op, err = f.ob.Enqueue(context.Background(), tx, req)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return op
}

func (f *fixture) get(t *testing.T, id string) Op {
	t.Helper()
	op, err := Get(context.Background(), f.pool, id)
	if err != nil {
		t.Fatal(err)
	}
	return op
}

func (f *fixture) pass(t *testing.T, pass func(context.Context) (int, error)) {
	t.Helper()
	if _, err := pass(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestEnqueueSignsWithConsecutiveNonces(t *testing.T) {
	f := newFixture(t, Config{})
	for want := range uint64(3) {
		op := f.enqueue(t, Request{Data: []byte{byte(want)}})
		if op.Status != Signed || op.Nonce != want || op.From != f.key.Address {
			t.Fatalf("op %d = %+v", want, op)
		}
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(f.get(t, op.ID).raw); err != nil {
			t.Fatal(err)
		}
		sender, err := types.Sender(types.LatestSignerForChainID(chain.ChainID1337), tx)
		if err != nil || sender != f.key.Address || tx.Hash() != op.TxHash || tx.Nonce() != want {
			t.Fatalf("stored tx: sender %s (%v), hash %s, nonce %d", sender, err, tx.Hash(), tx.Nonce())
		}
	}
	if n := f.node.sendCount(common.Hash{}); n != 0 || len(f.node.sends) != 0 {
		t.Fatal("Enqueue sent something")
	}
}

func TestEnqueueRequiresGas(t *testing.T) {
	f := newFixture(t, Config{})
	err := pgx.BeginFunc(context.Background(), f.pool, func(tx pgx.Tx) error {
		_, err := f.ob.Enqueue(context.Background(), tx, Request{KeyID: f.key.ID, To: &target})
		return err
	})
	if err == nil {
		t.Fatal("Enqueue without Gas succeeded")
	}
}

func TestEnqueueIdempotencyKeyReturnsFirstOp(t *testing.T) {
	f := newFixture(t, Config{})
	first := f.enqueue(t, Request{IdempotencyKey: "mint-42"})
	again := f.enqueue(t, Request{IdempotencyKey: "mint-42", Data: []byte("different")})
	if again.ID != first.ID || !again.Existing || first.Existing {
		t.Fatalf("second enqueue = %s (existing %v), want %s", again.ID, again.Existing, first.ID)
	}
	if k, _ := f.keys.Get(context.Background(), f.pool, f.key.ID); k.NextNonce != 1 {
		t.Fatalf("next_nonce = %d, want 1: the repeat reserved a nonce", k.NextNonce)
	}
}

func TestSendMarksSubmitted(t *testing.T) {
	f := newFixture(t, Config{})
	op := f.enqueue(t, Request{})
	if err := f.ob.Send(context.Background(), op.ID); err != nil {
		t.Fatal(err)
	}
	got := f.get(t, op.ID)
	if got.Status != Submitted || got.Attempts != 1 || f.node.sendCount(op.TxHash) != 1 {
		t.Fatalf("after Send: %s, %d attempts, %d sends", got.Status, got.Attempts, f.node.sendCount(op.TxHash))
	}
	// A second Send is a no-op: the row is no longer 'signed'.
	if err := f.ob.Send(context.Background(), op.ID); err != nil {
		t.Fatal(err)
	}
	if n := f.node.sendCount(op.TxHash); n != 1 {
		t.Fatalf("sent %d times, want 1", n)
	}
}

// The crash case: committed but never sent. The submitter finds it.
func TestSubmitterSendsCommittedButUnsentRows(t *testing.T) {
	f := newFixture(t, Config{})
	var ops []Op
	for range 3 {
		ops = append(ops, f.enqueue(t, Request{}))
	}
	f.pass(t, f.ob.submitPass)
	for _, op := range ops {
		if got := f.get(t, op.ID); got.Status != Submitted || f.node.sendCount(op.TxHash) != 1 {
			t.Fatalf("nonce %d: %s after %d sends", op.Nonce, got.Status, f.node.sendCount(op.TxHash))
		}
	}
	// Submitted rows are not resent before ResendAfter.
	f.pass(t, f.ob.submitPass)
	if n := f.node.sendCount(ops[0].TxHash); n != 1 {
		t.Fatalf("resent too early: %d sends", n)
	}
}

func TestSubmitterResendsWhenReceiptIsOverdue(t *testing.T) {
	f := newFixture(t, Config{ResendAfter: time.Millisecond})
	op := f.enqueue(t, Request{})
	f.pass(t, f.ob.submitPass)
	time.Sleep(5 * time.Millisecond)
	f.node.sendErr = func(*types.Transaction) error { return errKnown } // still in the node's pool
	f.pass(t, f.ob.submitPass)
	got := f.get(t, op.ID)
	if got.Status != Submitted || got.Attempts != 2 || f.node.sendCount(op.TxHash) != 2 {
		t.Fatalf("after resend: %s, %d attempts", got.Status, got.Attempts)
	}
}

func TestKnownTransactionCountsAsSent(t *testing.T) {
	f := newFixture(t, Config{})
	f.node.sendErr = func(*types.Transaction) error { return errKnown }
	op := f.enqueue(t, Request{})
	f.pass(t, f.ob.submitPass)
	if got := f.get(t, op.ID); got.Status != Submitted || got.LastError != "" {
		t.Fatalf("got %s, last error %q", got.Status, got.LastError)
	}
}

func TestConnectionErrorBacksOffAndKeepsRow(t *testing.T) {
	f := newFixture(t, Config{})
	f.node.sendErr = func(*types.Transaction) error { return errConnRefuse }
	op := f.enqueue(t, Request{})
	f.pass(t, f.ob.submitPass)
	got := f.get(t, op.ID)
	if got.Status != Signed || got.Attempts != 1 || got.LastError == "" {
		t.Fatalf("after refused connection: %s, %d attempts, last error %q", got.Status, got.Attempts, got.LastError)
	}
	// The row now waits out its backoff (1 s), so an immediate pass skips it.
	f.node.sendErr = nil
	f.pass(t, f.ob.submitPass)
	if n := f.node.sendCount(op.TxHash); n != 1 {
		t.Fatalf("retried during backoff: %d sends", n)
	}
	time.Sleep(1100 * time.Millisecond)
	f.pass(t, f.ob.submitPass)
	if got := f.get(t, op.ID); got.Status != Submitted || got.LastError != "" {
		t.Fatalf("after backoff: %s, last error %q", got.Status, got.LastError)
	}
}

// A resend of bytes that were already mined gets "nonce too low"; the
// receipt shows they are ours, so the row is recorded, not failed.
func TestNonceTooLowWithOurReceiptIsRecorded(t *testing.T) {
	f := newFixture(t, Config{})
	op := f.enqueue(t, Request{})
	f.node.mine(op.TxHash, types.ReceiptStatusSuccessful, 12, common.Address{}, "")
	f.node.sendErr = func(*types.Transaction) error { return errNonceLow }
	f.pass(t, f.ob.submitPass)
	got := f.get(t, op.ID)
	if got.Status != Mined || got.BlockNumber == nil || *got.BlockNumber != 12 {
		t.Fatalf("got %+v", got)
	}
}

func TestNonceTooLowWithoutReceiptFails(t *testing.T) {
	f := newFixture(t, Config{})
	f.node.sendErr = func(*types.Transaction) error { return errNonceLow }
	op := f.enqueue(t, Request{})
	f.pass(t, f.ob.submitPass)
	got := f.get(t, op.ID)
	if got.Status != Failed || got.LastError == "" {
		t.Fatalf("got %s, last error %q", got.Status, got.LastError)
	}
	if n := countRows(t, f.pool); n != 1 {
		t.Fatalf("%d rows, want 1: no filler, since the nonce is already used", n)
	}
}

// A rejected transaction fails, and a transfer to self takes its nonce so
// the next nonce is not stuck behind the gap.
func TestRejectedTransactionFailsAndFillsNonceGap(t *testing.T) {
	f := newFixture(t, Config{})
	bad := f.enqueue(t, Request{Gas: 1})
	next := f.enqueue(t, Request{})
	f.node.sendErr = func(tx *types.Transaction) error {
		if tx.Gas() == 1 {
			return errIntrinsic
		}
		return nil
	}
	f.pass(t, f.ob.submitPass)
	if got := f.get(t, bad.ID); got.Status != Failed || got.LastError != errIntrinsic.msg {
		t.Fatalf("bad tx: %s, last error %q", got.Status, got.LastError)
	}
	var fillerID string
	var fillerTo string
	if err := f.pool.QueryRow(context.Background(),
		`SELECT id, to_address FROM transactions WHERE nonce = $1 AND status <> 'failed'`, bad.Nonce,
	).Scan(&fillerID, &fillerTo); err != nil {
		t.Fatalf("no filler at nonce %d: %v", bad.Nonce, err)
	}
	if common.HexToAddress(fillerTo) != f.key.Address {
		t.Fatalf("filler goes to %s, want the sender itself", fillerTo)
	}
	f.pass(t, f.ob.submitPass) // sends the filler
	if got := f.get(t, fillerID); got.Status != Submitted {
		t.Fatalf("filler: %s", got.Status)
	}
	if got := f.get(t, next.ID); got.Status != Submitted {
		t.Fatalf("next nonce: %s", got.Status)
	}
}

// A filler the node rejects is not filled again, so there is no loop.
func TestRejectedFillerIsNotRefilled(t *testing.T) {
	f := newFixture(t, Config{})
	f.node.sendErr = func(*types.Transaction) error { return errIntrinsic }
	f.enqueue(t, Request{})
	f.pass(t, f.ob.submitPass) // fails the call, adds a filler
	f.pass(t, f.ob.submitPass) // fails the filler
	f.pass(t, f.ob.submitPass)
	if n := countRows(t, f.pool); n != 2 {
		t.Fatalf("%d rows, want the call and one filler", n)
	}
}

func TestReceiptLoopRecordsMinedAndReverted(t *testing.T) {
	var mu sync.Mutex
	var finals []Op
	f := newFixture(t, Config{OnFinal: func(_ context.Context, _ pgx.Tx, op Op) error {
		mu.Lock()
		defer mu.Unlock()
		finals = append(finals, op)
		return nil
	}})
	deploy := f.enqueueDeploy(t)
	f.ob.Send(context.Background(), deploy.ID)
	mint := f.enqueue(t, Request{})
	f.ob.Send(context.Background(), mint.ID)
	pending := f.enqueue(t, Request{})
	f.ob.Send(context.Background(), pending.ID)

	contract := common.HexToAddress("0x00000000000000000000000000000000000000cc")
	f.node.mine(deploy.TxHash, types.ReceiptStatusSuccessful, 20, contract, "")
	f.node.mine(mint.TxHash, types.ReceiptStatusFailed, 21, common.Address{}, "0x8574adcf")
	f.pass(t, f.ob.receiptPass)

	if got := f.get(t, deploy.ID); got.Status != Mined || got.ContractAddress == nil || *got.ContractAddress != contract {
		t.Fatalf("deploy: %+v", got)
	}
	if got := f.get(t, mint.ID); got.Status != Reverted || got.RevertReason != "0x8574adcf" || *got.BlockNumber != 21 {
		t.Fatalf("mint: %+v", got)
	}
	if got := f.get(t, pending.ID); got.Status != Submitted {
		t.Fatalf("pending: %s", got.Status)
	}
	if len(finals) != 2 {
		t.Fatalf("OnFinal ran %d times, want 2", len(finals))
	}
}

// A hook error rolls back the status change with it, so the row is retried.
func TestFinalHookErrorRollsBackOnlyThatRow(t *testing.T) {
	var badID string
	var hooked []string
	f := newFixture(t, Config{OnFinal: func(_ context.Context, _ pgx.Tx, op Op) error {
		if op.ID == badID {
			return errors.New("domain update failed")
		}
		hooked = append(hooked, op.ID)
		return nil
	}})
	var ops []Op
	for range 3 {
		op := f.enqueue(t, Request{})
		if err := f.ob.Send(context.Background(), op.ID); err != nil {
			t.Fatal(err)
		}
		f.node.mine(op.TxHash, types.ReceiptStatusSuccessful, 5, common.Address{}, "")
		ops = append(ops, op)
	}
	badID = ops[1].ID

	// One pass: the bad row rolls back alone, the other two are finalized.
	f.pass(t, f.ob.receiptPass)
	for i, want := range []Status{Mined, Submitted, Mined} {
		if got := f.get(t, ops[i].ID); got.Status != want {
			t.Fatalf("op %d: status = %s, want %s", i, got.Status, want)
		}
	}
	if len(hooked) != 2 {
		t.Fatalf("hook committed for %d rows, want 2", len(hooked))
	}

	// Once the hook works, the next pass finalizes the bad row too.
	badID = ""
	f.pass(t, f.ob.receiptPass)
	if got := f.get(t, ops[1].ID); got.Status != Mined {
		t.Fatalf("status = %s, want mined", got.Status)
	}
}

// Several workers on the same rows: each row is sent exactly once.
func TestConcurrentSubmittersSendEachRowOnce(t *testing.T) {
	f := newFixture(t, Config{BatchSize: 5})
	var ops []Op
	for range 20 {
		ops = append(ops, f.enqueue(t, Request{}))
	}
	workers := []*Outbox{f.ob, New(f.pool, f.keys, f.node, Config{BatchSize: 5}), New(f.pool, f.keys, f.node, Config{BatchSize: 5})}
	var wg sync.WaitGroup
	for _, w := range workers {
		for range 4 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range 3 {
					if _, err := w.submitPass(context.Background()); err != nil {
						t.Error(err)
					}
				}
			}()
		}
	}
	wg.Wait()
	for _, op := range ops {
		if n := f.node.sendCount(op.TxHash); n != 1 {
			t.Fatalf("nonce %d sent %d times", op.Nonce, n)
		}
	}
}

func TestClassify(t *testing.T) {
	for err, want := range map[error]sendResult{
		nil:                                sendOK,
		errKnown:                           sendKnown,
		nodeError{-32000, "already known"}: sendKnown, // geth
		errNonceLow:                        sendNonceUsed,
		errIntrinsic:                       sendRejected,
		nodeError{-32005, "Transaction gas limit exceeds block gas limit"}: sendRejected,
		nodeError{-32000, "Wrong chainId"}:                                 sendRejected,
		nodeError{-32000, "Transaction pool is full"}:                      sendTransient,
		errConnRefuse:            sendTransient,
		context.DeadlineExceeded: sendTransient,
	} {
		if got := classify(err); got != want {
			t.Errorf("classify(%v) = %d, want %d", err, got, want)
		}
	}
}

func countRows(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM transactions`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
