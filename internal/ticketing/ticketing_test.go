package ticketing

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
	"github.com/smlobo/kaleido-challenge/internal/dbtest"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
	"github.com/smlobo/kaleido-challenge/internal/outbox"
)

// These tests need Postgres only. The node is a fake that accepts every
// transaction and never mines; tests finish operations by hand with finish,
// which runs the real final-status hook.

const adminToken = "test-admin-token-0123456789"

type fakeNode struct {
	mu        sync.Mutex
	estimates func(to *common.Address, data []byte) error
}

func (f *fakeNode) SendRaw(_ context.Context, tx *types.Transaction) (common.Hash, error) {
	return tx.Hash(), nil
}

func (f *fakeNode) Receipt(context.Context, common.Hash) (*types.Receipt, string, error) {
	return nil, "", nil
}

func (f *fakeNode) EstimateGas(_ context.Context, _ common.Address, to *common.Address, _ *big.Int, data []byte) (uint64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.estimates != nil {
		if err := f.estimates(to, data); err != nil {
			return 0, err
		}
	}
	return 100_000, nil
}

// revertErr is how go-ethereum's rpc client reports a revert from
// eth_estimateGas: code 3, with the revert data as hex.
type revertErr struct{ data string }

func (e revertErr) Error() string          { return "execution reverted" }
func (e revertErr) ErrorCode() int         { return 3 }
func (e revertErr) ErrorData() interface{} { return e.data }

type fixture struct {
	pool  *pgxpool.Pool
	svc   *Service
	node  *fakeNode
	admin Principal
}

func setup(t *testing.T) *fixture {
	t.Helper()
	pool := dbtest.Pool(t)
	kek := make([]byte, keystore.KEKSize)
	_, _ = rand.Read(kek)
	keys, err := keystore.New(kek, chain.ChainID1337)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	deployer, err := keys.EnsureDeployer(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	node := &fakeNode{}
	ob := outbox.New(pool, keys, node, outbox.Config{OnFinal: OnFinal})
	svc := New(Config{Pool: pool, Keys: keys, Outbox: ob, Node: node, Deployer: deployer, ChainID: chain.ChainID1337,
		PublicBaseURL: "http://tickets.test", AdminToken: adminToken})
	admin, err := svc.Authenticate(ctx, adminToken)
	if err != nil {
		t.Fatal(err)
	}
	return &fixture{pool: pool, svc: svc, node: node, admin: admin}
}

func (f *fixture) organizer(t *testing.T) Principal {
	t.Helper()
	o, err := f.svc.CreateOrganizer(context.Background(), f.admin, "Org "+t.Name())
	if err != nil {
		t.Fatal(err)
	}
	return f.login(t, o.Token)
}

func (f *fixture) user(t *testing.T, subject string) Principal {
	t.Helper()
	u, err := f.svc.CreateUser(context.Background(), subject)
	if err != nil {
		t.Fatal(err)
	}
	return f.login(t, u.Token)
}

func (f *fixture) login(t *testing.T, token string) Principal {
	t.Helper()
	p, err := f.svc.Authenticate(context.Background(), token)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func newEvent(supply int64) NewEvent {
	return NewEvent{Name: "Test Fest", Symbol: "TF", StartsAt: time.Now().Add(24 * time.Hour),
		MaxSupply: supply, FaceValue: "100000000", MaxResaleBps: 11_000}
}

// finish records an operation as final and runs the hook, as the receipt
// loop does when it sees a receipt.
func (f *fixture) finish(t *testing.T, opID string, status outbox.Status, contract *common.Address) {
	t.Helper()
	err := pgx.BeginFunc(context.Background(), f.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(context.Background(), `UPDATE transactions SET status = $2 WHERE id = $1`, opID, string(status)); err != nil {
			return err
		}
		return OnFinal(context.Background(), tx, outbox.Op{ID: opID, Status: status, ContractAddress: contract})
	})
	if err != nil {
		t.Fatal(err)
	}
}

// onSale creates an event and marks its deploy mined.
func (f *fixture) onSale(t *testing.T, org Principal, supply int64) Event {
	t.Helper()
	ctx := context.Background()
	e, err := f.svc.CreateEvent(ctx, org, newEvent(supply), "deploy-"+t.Name())
	if err != nil {
		t.Fatal(err)
	}
	addr := common.BytesToAddress(crypto.Keccak256([]byte(e.ID))[:20])
	f.finish(t, e.DeployOperationID, outbox.Mined, &addr)
	e, err = f.svc.GetEvent(ctx, e.ID)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func wantRefusal(t *testing.T, err error, kind Kind, code string) {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) || e.Kind != kind || e.Code != code {
		t.Fatalf("err = %v, want kind %d code %s", err, kind, code)
	}
}

func TestTokensResolveToTheirOwner(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	org := f.organizer(t)
	user := f.user(t, "alice")
	if f.admin.Role != Admin || org.Role != Organizer || user.Role != User || org.ID == "" || user.ID == "" {
		t.Fatalf("admin %+v, organizer %+v, user %+v", f.admin, org, user)
	}
	if _, err := f.svc.Authenticate(ctx, "tk_not-a-real-token"); !errors.Is(err, errBadToken) {
		t.Fatalf("unknown token: %v", err)
	}
	if p, err := f.svc.Authenticate(ctx, ""); err != nil || p.Role != Anonymous {
		t.Fatalf("no token: %+v %v", p, err)
	}
	_, err := f.svc.CreateOrganizer(ctx, user, "Sneaky Ltd")
	wantRefusal(t, err, Forbidden, "forbidden")
	_, err = f.svc.CreateUser(ctx, "alice")
	wantRefusal(t, err, Conflict, "subject_taken")

	// Only hashes are stored.
	var n int
	if err := f.pool.QueryRow(ctx, `SELECT count(*) FROM api_tokens WHERE token_hash::text LIKE '%tk_%'`).Scan(&n); err != nil || n != 0 {
		t.Fatalf("plain tokens stored: %d %v", n, err)
	}
}

func TestNoAdminTokenMeansNoOrganizers(t *testing.T) {
	f := setup(t)
	f.svc.adminHash = nil
	_, err := f.svc.CreateOrganizer(context.Background(), Principal{Role: User, ID: "x"}, "Org")
	wantRefusal(t, err, Forbidden, "admin_disabled")
}

func TestCreateEventReportsEveryProblem(t *testing.T) {
	f := setup(t)
	_, err := f.svc.CreateEvent(context.Background(), f.organizer(t), NewEvent{
		Name: "", Symbol: "rock fest", StartsAt: time.Now(), MaxSupply: 0, FaceValue: "-5", MaxResaleBps: -1,
	}, "k")
	wantRefusal(t, err, Invalid, "invalid_event")
	for _, field := range []string{"name", "symbol", "starts_at", "max_supply", "face_value", "max_resale_bps"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error does not mention %s: %v", field, err)
		}
	}
}

func TestEventDeployOutcomes(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	org := f.organizer(t)

	e, err := f.svc.CreateEvent(ctx, org, newEvent(10), "a")
	if err != nil {
		t.Fatal(err)
	}
	if e.Status != "deploying" || e.ContractAddress != nil || e.OrganizerID != org.ID {
		t.Fatalf("new event: %+v", e)
	}
	op, err := f.svc.GetOperation(ctx, org, e.DeployOperationID)
	if err != nil || op.To != nil {
		t.Fatalf("deploy operation: %+v %v", op, err)
	}

	// A repeat with the same key returns the same event and signs nothing.
	again, err := f.svc.CreateEvent(ctx, org, newEvent(10), "a")
	if err != nil || again.ID != e.ID {
		t.Fatalf("repeat: %+v %v", again, err)
	}

	addr := common.HexToAddress("0x00000000000000000000000000000000000000e1")
	f.finish(t, e.DeployOperationID, outbox.Mined, &addr)
	f.finish(t, e.DeployOperationID, outbox.Mined, &addr) // the hook runs twice harmlessly
	if e, _ = f.svc.GetEvent(ctx, e.ID); e.Status != "on_sale" || e.ContractAddress == nil || *e.ContractAddress != strings.ToLower(addr.Hex()) {
		t.Fatalf("after deploy: %+v", e)
	}

	b, err := f.svc.CreateEvent(ctx, org, newEvent(10), "b")
	if err != nil {
		t.Fatal(err)
	}
	f.finish(t, b.DeployOperationID, outbox.Reverted, nil)
	if b, _ = f.svc.GetEvent(ctx, b.ID); b.Status != "failed" {
		t.Fatalf("after revert: %+v", b)
	}
	_, err = f.svc.BuyTicket(ctx, f.user(t, "bob"), b.ID, "k")
	wantRefusal(t, err, Rejected, "event_failed")
}

func TestBuyTicketRules(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	org := f.organizer(t)
	alice := f.user(t, "alice")

	deploying, err := f.svc.CreateEvent(ctx, org, newEvent(2), "d")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.svc.BuyTicket(ctx, alice, deploying.ID, "k")
	wantRefusal(t, err, Conflict, "event_not_ready")

	e := f.onSale(t, org, 2)
	_, err = f.svc.BuyTicket(ctx, org, e.ID, "k")
	wantRefusal(t, err, Forbidden, "forbidden")
	_, err = f.svc.BuyTicket(ctx, Principal{}, e.ID, "k")
	wantRefusal(t, err, Unauthorized, "unauthorized")

	t1, err := f.svc.BuyTicket(ctx, alice, e.ID, "first")
	if err != nil {
		t.Fatal(err)
	}
	if t1.Status != "pending_mint" || t1.TokenID != "1" || t1.OwnerAddress != nil ||
		t1.MetadataURL != "http://tickets.test/metadata/"+e.ID+"/1" {
		t.Fatalf("ticket: %+v", t1)
	}
	again, err := f.svc.BuyTicket(ctx, alice, e.ID, "first")
	if err != nil || again.ID != t1.ID {
		t.Fatalf("repeat: %+v %v", again, err)
	}
	if _, err := f.svc.BuyTicket(ctx, alice, e.ID, "second"); err != nil {
		t.Fatal(err)
	}
	_, err = f.svc.BuyTicket(ctx, alice, e.ID, "third")
	wantRefusal(t, err, Rejected, "sold_out")

	// The mint is signed by the box office, for the buyer's wallet.
	op, err := f.svc.GetOperation(ctx, alice, t1.MintOperationID)
	if err != nil {
		t.Fatal(err)
	}
	var boxOffice string
	_ = f.pool.QueryRow(ctx, `SELECT k.address FROM organizers o JOIN signing_keys k ON k.id = o.box_office_key_id WHERE o.id = $1`, org.ID).Scan(&boxOffice)
	if strings.ToLower(op.From.Hex()) != boxOffice || op.To == nil || strings.ToLower(op.To.Hex()) != *e.ContractAddress {
		t.Fatalf("mint from %s to %v; box office %s, contract %s", op.From.Hex(), op.To, boxOffice, *e.ContractAddress)
	}

	// Mined: active and owned by the buyer's wallet.
	f.finish(t, t1.MintOperationID, outbox.Mined, nil)
	t1, _ = f.svc.GetTicket(ctx, t1.ID)
	var wallet string
	_ = f.pool.QueryRow(ctx, `SELECT k.address FROM users u JOIN signing_keys k ON k.id = u.wallet_key_id WHERE u.id = $1`, alice.ID).Scan(&wallet)
	if t1.Status != "active" || t1.OwnerAddress == nil || *t1.OwnerAddress != wallet {
		t.Fatalf("after mint: %+v, wallet %s", t1, wallet)
	}

	// Started events sell nothing.
	if _, err := f.pool.Exec(ctx, `UPDATE events SET starts_at = now() - interval '1 minute', tickets_issued = 0 WHERE id = $1`, e.ID); err != nil {
		t.Fatal(err)
	}
	_, err = f.svc.BuyTicket(ctx, alice, e.ID, "late")
	wantRefusal(t, err, Rejected, "event_started")
}

// Twenty buyers at once for five tickets: exactly five get one, numbered 1-5.
func TestConcurrentBuyersNeverOversell(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	e := f.onSale(t, f.organizer(t), 5)
	buyers := make([]Principal, 20)
	for i := range buyers {
		buyers[i] = f.user(t, fmt.Sprintf("buyer%d", i))
	}
	var wg sync.WaitGroup
	tickets := make([]Ticket, len(buyers))
	errs := make([]error, len(buyers))
	for i := range buyers {
		wg.Go(func() { tickets[i], errs[i] = f.svc.BuyTicket(ctx, buyers[i], e.ID, "k") })
	}
	wg.Wait()
	got := map[string]bool{}
	for i, err := range errs {
		if err == nil {
			got[tickets[i].TokenID] = true
			continue
		}
		wantRefusal(t, err, Rejected, "sold_out")
	}
	if len(got) != 5 || !got["1"] || !got["5"] {
		t.Fatalf("token IDs sold: %v", got)
	}
}

func TestCheckInRules(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	org := f.organizer(t)
	e := f.onSale(t, org, 5)
	alice := f.user(t, "alice")
	tk, err := f.svc.BuyTicket(ctx, alice, e.ID, "k")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.svc.CreatePass(ctx, alice, tk.ID)
	wantRefusal(t, err, Conflict, "not_minted_yet")
	f.finish(t, tk.MintOperationID, outbox.Mined, nil)

	// Only the holder can make a pass, and only a user.
	_, err = f.svc.CreatePass(ctx, f.user(t, "bob"), tk.ID)
	wantRefusal(t, err, Forbidden, "not_your_ticket")
	_, err = f.svc.CreatePass(ctx, org, tk.ID)
	wantRefusal(t, err, Forbidden, "forbidden")

	pass, err := f.svc.CreatePass(ctx, alice, tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if pass.TokenID != tk.TokenID || pass.Contract != *e.ContractAddress || pass.ChainID.Cmp(chain.ChainID1337) != 0 {
		t.Fatalf("pass: %+v", pass)
	}
	if left := time.Until(pass.Deadline); left <= 4*time.Minute || left > PassLifetime {
		t.Fatalf("pass valid for %s, want about %s", left, PassLifetime)
	}
	in := PassInput{TicketID: tk.ID, Deadline: pass.Deadline.Unix(), Signature: pass.Signature}

	// Refused before anything is signed: wrong organizer, a pass that is
	// expired, malformed, or for another ticket.
	other := f.organizer(t)
	_, err = f.svc.CheckIn(ctx, other, tk.ID, in, "c1")
	wantRefusal(t, err, Forbidden, "not_your_event")
	_, err = f.svc.CheckIn(ctx, org, tk.ID, PassInput{Deadline: time.Now().Add(-time.Second).Unix(), Signature: pass.Signature}, "c1")
	wantRefusal(t, err, Rejected, "pass_expired")
	_, err = f.svc.CheckIn(ctx, org, tk.ID, PassInput{Deadline: in.Deadline, Signature: "0x1234"}, "c1")
	wantRefusal(t, err, Invalid, "invalid_pass")
	_, err = f.svc.CheckIn(ctx, org, tk.ID, PassInput{TicketID: e.ID, Deadline: in.Deadline, Signature: pass.Signature}, "c1")
	wantRefusal(t, err, Invalid, "pass_for_another_ticket")

	// The check-in carries the pass, signed by alice's wallet over the
	// contract's EIP-712 digest.
	var sent []byte
	f.node.estimates = func(_ *common.Address, data []byte) error { sent = data; return nil }
	c, err := f.svc.CheckIn(ctx, org, tk.ID, in, "c1")
	if err != nil {
		t.Fatal(err)
	}
	if c.CheckInOperationID == nil || c.Status != "active" {
		t.Fatalf("after check-in request: %+v", c)
	}
	args, err := ticketABI().Methods["checkIn"].Inputs.Unpack(sent[4:])
	if err != nil {
		t.Fatal(err)
	}
	token, deadline, sig := args[0].(*big.Int), args[1].(*big.Int), args[2].([]byte)
	if token.String() != tk.TokenID || deadline.Int64() != in.Deadline || "0x"+common.Bytes2Hex(sig) != pass.Signature {
		t.Fatalf("checkIn(%s, %s, %x) does not carry the pass %+v", token, deadline, sig, pass)
	}
	digest, err := ticketnft.CheckInDigest(chain.ChainID1337, common.HexToAddress(pass.Contract), token, deadline)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] -= 27
	pub, err := crypto.SigToPub(digest[:], sig)
	if err != nil {
		t.Fatal(err)
	}
	if owner := strings.ToLower(crypto.PubkeyToAddress(*pub).Hex()); c.OwnerAddress == nil || owner != *c.OwnerAddress {
		t.Fatalf("pass signed by %s, holder is %v", owner, c.OwnerAddress)
	}

	again, err := f.svc.CheckIn(ctx, org, tk.ID, in, "c1")
	if err != nil || *again.CheckInOperationID != *c.CheckInOperationID {
		t.Fatalf("repeat with the same key: %+v %v", again, err)
	}
	_, err = f.svc.CheckIn(ctx, org, tk.ID, in, "c2")
	wantRefusal(t, err, Conflict, "already_checked_in")

	// A reverted check-in can be tried again, with a new pass.
	f.finish(t, *c.CheckInOperationID, outbox.Reverted, nil)
	if pass, err = f.svc.CreatePass(ctx, alice, tk.ID); err != nil {
		t.Fatal(err)
	}
	c2, err := f.svc.CheckIn(ctx, org, tk.ID, PassInput{Deadline: pass.Deadline.Unix(), Signature: pass.Signature}, "c3")
	if err != nil {
		t.Fatal(err)
	}
	f.finish(t, *c2.CheckInOperationID, outbox.Mined, nil)
	if c2, _ = f.svc.GetTicket(ctx, tk.ID); c2.Status != "checked_in" {
		t.Fatalf("after check-in mined: %+v", c2)
	}
	_, err = f.svc.CheckIn(ctx, org, tk.ID, in, "c4")
	wantRefusal(t, err, Conflict, "already_checked_in")
	_, err = f.svc.CreatePass(ctx, alice, tk.ID)
	wantRefusal(t, err, Conflict, "already_checked_in")
}

// The contract's refusal of a pass reaches the caller as a 422, unsigned.
func TestCheckInPassRefusedByContract(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	org := f.organizer(t)
	e := f.onSale(t, org, 5)
	alice := f.user(t, "alice")
	tk, err := f.svc.BuyTicket(ctx, alice, e.ID, "k")
	if err != nil {
		t.Fatal(err)
	}
	f.finish(t, tk.MintOperationID, outbox.Mined, nil)
	pass, err := f.svc.CreatePass(ctx, alice, tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	in := PassInput{Deadline: pass.Deadline.Unix(), Signature: pass.Signature}

	notHolder := ticketABI().Errors["NotSignedByHolder"]
	data, _ := notHolder.Inputs.Pack(big.NewInt(1))
	f.node.estimates = func(*common.Address, []byte) error {
		return revertErr{data: "0x" + common.Bytes2Hex(append(notHolder.ID[:4:4], data...))}
	}
	_, err = f.svc.CheckIn(ctx, org, tk.ID, in, "c1")
	wantRefusal(t, err, Rejected, "pass_not_from_holder")
	if tk, _ = f.svc.GetTicket(ctx, tk.ID); tk.CheckInOperationID != nil {
		t.Fatalf("a refused check-in was saved: %+v", tk)
	}
}

// A call the contract would revert is refused before anything is signed.
func TestDryRunRevertIsRefusedUnsigned(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	e := f.onSale(t, f.organizer(t), 5)
	soldOut := ticketABI().Errors["SoldOut"].ID
	f.node.estimates = func(*common.Address, []byte) error {
		return fmt.Errorf("eth_estimateGas: %w", revertErr{data: "0x" + common.Bytes2Hex(soldOut[:4])})
	}
	var before, after int
	_ = f.pool.QueryRow(ctx, `SELECT count(*) FROM transactions`).Scan(&before)
	_, err := f.svc.BuyTicket(ctx, f.user(t, "alice"), e.ID, "k")
	wantRefusal(t, err, Rejected, "sold_out")
	_ = f.pool.QueryRow(ctx, `SELECT count(*) FROM transactions`).Scan(&after)
	if after != before {
		t.Fatalf("signed %d transactions for a refused request", after-before)
	}

	// A node that cannot be reached does not stop the request: the outbox sends it later.
	f.node.estimates = func(*common.Address, []byte) error { return errors.New("connection refused") }
	if _, err := f.svc.BuyTicket(ctx, f.user(t, "bob"), e.ID, "k"); err != nil {
		t.Fatalf("with the node down: %v", err)
	}
}

func TestOperationsAreVisibleOnlyToTheirRequester(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	e := f.onSale(t, f.organizer(t), 5)
	alice, bob := f.user(t, "alice"), f.user(t, "bob")
	tk, err := f.svc.BuyTicket(ctx, alice, e.ID, "k")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.GetOperation(ctx, alice, tk.MintOperationID); err != nil {
		t.Fatal(err)
	}
	_, err = f.svc.GetOperation(ctx, bob, tk.MintOperationID)
	wantRefusal(t, err, NotFound, "not_found")
	_, err = f.svc.GetOperation(ctx, Principal{}, tk.MintOperationID)
	wantRefusal(t, err, Unauthorized, "unauthorized")
}

// The hash on-chain must match the record however Postgres stores it.
func TestRecordHashSurvivesJSONB(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	e := f.onSale(t, f.organizer(t), 5)
	tk, err := f.svc.BuyTicket(ctx, f.user(t, "alice"), e.ID, "k")
	if err != nil {
		t.Fatal(err)
	}
	var stored string
	_ = f.pool.QueryRow(ctx, `SELECT record::text FROM tickets WHERE id = $1`, tk.ID).Scan(&stored)
	if !strings.Contains(stored, ": ") {
		t.Fatalf("expected jsonb's spaced text, got %s", stored)
	}
	h, err := RecordHash([]byte(stored))
	if err != nil || h.Hex() != tk.RecordHash {
		t.Fatalf("hash of stored record %s, ticket says %s (%v)", h.Hex(), tk.RecordHash, err)
	}
	rec := string(tk.Record)
	if !strings.HasPrefix(rec, `{"event_id":"`+e.ID+`","face_value":"100000000","issued_at":"`) ||
		!strings.HasSuffix(rec, `Z","ticket_number":1,"tier":"general"}`) || strings.Contains(rec, " ") {
		t.Fatalf("record not canonical: %s", rec)
	}

	m, err := f.svc.GetMetadata(ctx, e.ID, "1")
	if err != nil || m.RecordHash != tk.RecordHash || m.Name != "Test Fest #1" {
		t.Fatalf("metadata: %+v %v", m, err)
	}
	_, err = f.svc.GetMetadata(ctx, e.ID, "2")
	wantRefusal(t, err, NotFound, "not_found")
}

func TestDecodeRevert(t *testing.T) {
	capErr, soldOut := ticketABI().Errors["PriceAboveCap"], ticketABI().Errors["SoldOut"]
	args, _ := capErr.Inputs.Pack(big.NewInt(120), big.NewInt(110))
	for _, c := range []struct {
		data []byte
		want string
	}{
		{append(capErr.ID[:4:4], args...), "PriceAboveCap(120, 110)"},
		{soldOut.ID[:4], "SoldOut()"},
		{common.FromHex("0x08c379a0" +
			"0000000000000000000000000000000000000000000000000000000000000020" +
			"0000000000000000000000000000000000000000000000000000000000000004" +
			"6e6f706500000000000000000000000000000000000000000000000000000000"), `Error("nope")`},
		{[]byte{0xde, 0xad, 0xbe, 0xef}, "0xdeadbeef"},
		{nil, "a revert with no reason"},
	} {
		if got := DecodeRevert(c.data); got != c.want {
			t.Errorf("DecodeRevert(%x) = %s, want %s", c.data, got, c.want)
		}
	}
}
