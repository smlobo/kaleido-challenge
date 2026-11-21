// Package outbox is the service's write path to the chain. A transaction is
// signed and saved to Postgres before it is sent, so any instance can finish
// sending it after a crash:
//
//  1. Enqueue, inside the caller's database transaction: reserve the key's
//     next nonce, sign, insert the row as 'signed'. The caller inserts its
//     own rows (a ticket, an event) in the same transaction and commits.
//  2. Send (fast path), right after the commit: eth_sendRawTransaction, then
//     'submitted'. Optional: the submitter loop does it anyway.
//  3. The submitter loop sends 'signed' rows, and resends 'submitted' rows
//     with no receipt after ResendAfter, using FOR UPDATE SKIP LOCKED so
//     instances share the work. Resending is safe: the same bytes have the
//     same hash, so the node never mines them twice.
//  4. The receipt loop records each receipt: 'mined' or 'reverted'.
//
// See the design doc's "Transaction pipeline" section.
package outbox

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
	"github.com/smlobo/kaleido-challenge/internal/telemetry"
)

// Status is a transactions row's status column.
type Status string

const (
	Signed    Status = "signed"    // saved, not yet accepted by a node
	Submitted Status = "submitted" // a node accepted it; waiting for a receipt
	Mined     Status = "mined"     // receipt status 1
	Reverted  Status = "reverted"  // receipt status 0; the nonce is used
	Failed    Status = "failed"    // never mined; see LastError
)

// Final reports whether no loop will touch the row again.
func (s Status) Final() bool { return s == Mined || s == Reverted || s == Failed }

// Node is the part of the chain client the outbox uses. *chain.Client
// satisfies it; tests use a fake.
type Node interface {
	SendRaw(ctx context.Context, tx *types.Transaction) (common.Hash, error)
	Receipt(ctx context.Context, hash common.Hash) (*types.Receipt, string, error)
}

var _ Node = (*chain.Client)(nil)

// Request is one call or deployment to sign and queue.
type Request struct {
	KeyID string          // signing_keys.id of the sender
	To    *common.Address // nil deploys a contract
	Data  []byte          // calldata, or bytecode + constructor args
	Gas   uint64          // gas limit; required (estimate it, or a fixed limit per method)
	// IdempotencyKey, if set, makes a repeated request return the first
	// operation instead of signing a second transaction.
	IdempotencyKey string
}

// Op is a transactions row: one operation, as the API reports it.
type Op struct {
	ID              string
	From            common.Address
	Nonce           uint64
	To              *common.Address
	Calldata        []byte
	TxHash          common.Hash
	Status          Status
	BlockNumber     *uint64
	ContractAddress *common.Address // set for a mined deployment
	RevertReason    string          // hex revert data from the receipt
	LastError       string
	Attempts        int
	// Existing is true when Enqueue found an operation with the same
	// idempotency key instead of creating one.
	Existing bool

	raw         []byte
	traceparent string // the queuing request's trace context; "" if untraced
}

// FinalHook runs in the same database transaction that gives a row its final
// status, so the caller's own rows (a ticket becoming active or failed, an
// event's contract address) change atomically with it.
type FinalHook func(ctx context.Context, tx pgx.Tx, op Op) error

// Config tunes the loops. Zero values take the defaults.
type Config struct {
	PollInterval time.Duration // how often each loop looks for work (1 s)
	ResendAfter  time.Duration // resend a submitted row with no receipt after this (10 s)
	MaxBackoff   time.Duration // longest wait between retries after a connection error (1 min)
	BatchSize    int           // rows claimed per pass (50)
	OnFinal      FinalHook
	Logger       *slog.Logger
	// ABIs name the method a transaction calls (eth.method on its spans),
	// from the first 4 bytes of its calldata. A selector none of them knows
	// shows as hex.
	ABIs []*abi.ABI
}

// Outbox signs, saves, sends and confirms transactions.
type Outbox struct {
	db   DB
	keys *keystore.Store
	node Node
	cfg  Config
	log  *slog.Logger
}

// DB is satisfied by *pgxpool.Pool.
type DB interface {
	keystore.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

// New returns an Outbox. Call Run to start its loops.
func New(db DB, keys *keystore.Store, node Node, cfg Config) *Outbox {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = time.Second
	}
	if cfg.ResendAfter <= 0 {
		cfg.ResendAfter = 10 * time.Second
	}
	if cfg.MaxBackoff <= 0 {
		cfg.MaxBackoff = time.Minute
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 50
	}
	log := cfg.Logger
	if log == nil {
		log = slog.Default()
	}
	return &Outbox{db: db, keys: keys, node: node, cfg: cfg, log: log.With("component", "outbox")}
}

// Enqueue reserves the key's next nonce, signs the transaction and inserts it
// as 'signed', all inside tx. Nothing is sent: the caller commits, then calls
// Send or leaves it to the submitter. If tx rolls back, the nonce goes back
// with it.
//
// With an IdempotencyKey that already exists, Enqueue returns that operation
// with Existing set and signs nothing. Two requests racing with the same new
// key both sign; the second insert fails on the unique index and its
// transaction must be rolled back.
func (o *Outbox) Enqueue(ctx context.Context, tx pgx.Tx, req Request) (Op, error) {
	if req.Gas == 0 {
		return Op{}, errors.New("outbox: Request.Gas is required")
	}
	// The loops' later spans hang off the request's span, not this one.
	traceparent := telemetry.Inject(ctx)
	ctx, span := telemetry.Child(ctx, "outbox enqueue")
	op, err := o.enqueue(ctx, tx, req, traceparent)
	span.SetAttributes(o.opAttrs(op)...)
	span.SetAttributes(attribute.Bool("outbox.existing", op.Existing))
	telemetry.EndSpan(span, err)
	return op, err
}

func (o *Outbox) enqueue(ctx context.Context, tx pgx.Tx, req Request, traceparent string) (Op, error) {
	if req.IdempotencyKey != "" {
		op, err := scanOp(tx.QueryRow(ctx, selectOp+` WHERE idempotency_key = $1`, req.IdempotencyKey))
		if err == nil {
			op.Existing = true
			return op, nil
		}
		if !errors.Is(err, ErrNotFound) {
			return Op{}, err
		}
	}
	nonce, err := keystore.ReserveNonce(ctx, tx, req.KeyID)
	if err != nil {
		return Op{}, err
	}
	signer, err := o.keys.Signer(ctx, tx, req.KeyID)
	if err != nil {
		return Op{}, err
	}
	return insert(ctx, tx, signer, nonce, req.To, req.Data, req.Gas, req.IdempotencyKey, traceparent)
}

// Send is the fast path: send a 'signed' row right after its commit. It does
// nothing if the row is no longer 'signed' or another worker holds it. Send
// errors are recorded on the row, not returned; only database errors are.
func (o *Outbox) Send(ctx context.Context, id string) error {
	return pgx.BeginFunc(ctx, o.db, func(tx pgx.Tx) error {
		op, err := scanOp(tx.QueryRow(ctx,
			selectOp+` WHERE id = $1 AND status = 'signed' FOR UPDATE SKIP LOCKED`, id))
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		ctx, span := telemetry.Child(ctx, "outbox send", o.opAttrs(op)...)
		err = o.attempt(ctx, tx, op)
		telemetry.EndSpan(span, err)
		return err
	})
}

// opAttrs describes an operation on a span. Besu's own spans show only the
// JSON-RPC method: what the transaction does is known here, not there.
func (o *Outbox) opAttrs(op Op) []attribute.KeyValue {
	if op.ID == "" {
		return nil
	}
	attrs := []attribute.KeyValue{
		attribute.String("outbox.op_id", op.ID),
		attribute.String("outbox.status", string(op.Status)),
		attribute.String("eth.from", op.From.Hex()),
		attribute.Int64("eth.nonce", int64(op.Nonce)),
		attribute.String("eth.tx_hash", op.TxHash.Hex()),
		attribute.String("eth.method", methodName(o.cfg.ABIs, op)),
	}
	if op.To != nil { // none for a contract creation
		attrs = append(attrs, attribute.String("eth.to", op.To.Hex()))
	}
	return attrs
}

// methodName is "deploy" for a contract creation, the ABI name of the called
// method (mintTicket), its selector in hex if no ABI has it, or "transfer"
// for a plain value transfer with no calldata.
func methodName(abis []*abi.ABI, op Op) string {
	switch {
	case op.To == nil:
		return "deploy"
	case len(op.Calldata) < 4:
		return "transfer"
	}
	for _, a := range abis {
		if m, err := a.MethodById(op.Calldata[:4]); err == nil {
			return m.Name
		}
	}
	return hexutil.Encode(op.Calldata[:4])
}

// Get loads one operation.
func Get(ctx context.Context, db keystore.DBTX, id string) (Op, error) {
	return scanOp(db.QueryRow(ctx, selectOp+` WHERE id = $1`, id))
}

// Wait polls until the operation reaches a final status or ctx ends.
func (o *Outbox) Wait(ctx context.Context, id string) (Op, error) {
	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()
	for {
		op, err := Get(ctx, o.db, id)
		if err != nil || op.Status.Final() {
			return op, err
		}
		select {
		case <-ctx.Done():
			return op, ctx.Err()
		case <-t.C:
		}
	}
}

// ErrNotFound is returned when no operation matches.
var ErrNotFound = errors.New("outbox: operation not found")

const selectOp = `SELECT id, from_address, nonce, to_address, calldata, raw_tx, tx_hash, status,
	block_number, contract_address, coalesce(revert_reason, ''), coalesce(last_error, ''), attempts,
	coalesce(traceparent, '')
	FROM transactions`

func scanOp(row pgx.Row) (Op, error) {
	var op Op
	var from, hash, status string
	var to, contract *string
	var nonce int64
	var block *int64
	err := row.Scan(&op.ID, &from, &nonce, &to, &op.Calldata, &op.raw, &hash, &status,
		&block, &contract, &op.RevertReason, &op.LastError, &op.Attempts, &op.traceparent)
	if errors.Is(err, pgx.ErrNoRows) {
		return Op{}, ErrNotFound
	}
	if err != nil {
		return Op{}, err
	}
	op.From = common.HexToAddress(from)
	op.Nonce = uint64(nonce)
	op.TxHash = common.HexToHash(hash)
	op.Status = Status(status)
	if to != nil {
		a := common.HexToAddress(*to)
		op.To = &a
	}
	if contract != nil {
		a := common.HexToAddress(*contract)
		op.ContractAddress = &a
	}
	if block != nil {
		b := uint64(*block)
		op.BlockNumber = &b
	}
	return op, nil
}

func insert(ctx context.Context, tx pgx.Tx, signer *chain.Signer, nonce uint64, to *common.Address,
	data []byte, gas uint64, idemKey, traceparent string) (Op, error) {
	signed, err := signer.Sign(chain.TxRequest{Nonce: nonce, To: to, Data: data, Gas: gas})
	if err != nil {
		return Op{}, err
	}
	raw, err := signed.MarshalBinary()
	if err != nil {
		return Op{}, err
	}
	var key *string
	if idemKey != "" {
		key = &idemKey
	}
	if data == nil {
		data = []byte{}
	}
	var tp *string
	if traceparent != "" {
		tp = &traceparent
	}
	var id string
	if err := tx.QueryRow(ctx,
		`INSERT INTO transactions (idempotency_key, from_address, nonce, to_address, calldata, raw_tx, tx_hash, traceparent)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		key, hexAddr(signer.Address()), int64(nonce), hexAddrPtr(to), data, raw, hexHash(signed.Hash()), tp,
	).Scan(&id); err != nil {
		return Op{}, fmt.Errorf("insert transaction: %w", err)
	}
	return Op{ID: id, From: signer.Address(), Nonce: nonce, To: to, Calldata: data,
		TxHash: signed.Hash(), Status: Signed, raw: raw, traceparent: traceparent}, nil
}

// The schema's eth_address and eth_hash domains require lowercase hex.
func hexAddr(a common.Address) string { return strings.ToLower(a.Hex()) }
func hexHash(h common.Hash) string    { return strings.ToLower(h.Hex()) }

func hexAddrPtr(a *common.Address) *string {
	if a == nil {
		return nil
	}
	s := hexAddr(*a)
	return &s
}
