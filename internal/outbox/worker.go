package outbox

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jackc/pgx/v5"
)

// rpcTimeout bounds each node call, so a hung node cannot hold row locks.
const rpcTimeout = 10 * time.Second

// fillerGas is the gas of a plain value transfer, which a filler is.
const fillerGas = 21_000

// Run runs the submitter and receipt loops until ctx ends. Run it on every
// instance: FOR UPDATE SKIP LOCKED gives each row to one worker at a time.
func (o *Outbox) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for name, pass := range map[string]func(context.Context) (int, error){
		"submitter": o.submitPass,
		"receipts":  o.receiptPass,
	} {
		wg.Go(func() {
			o.loop(ctx, name, pass)
		})
	}
	wg.Wait()
}

func (o *Outbox) loop(ctx context.Context, name string, pass func(context.Context) (int, error)) {
	t := time.NewTicker(o.cfg.PollInterval)
	defer t.Stop()
	for {
		// Keep going while passes come back full, so a backlog drains at once.
		for {
			n, err := pass(ctx)
			if err != nil && ctx.Err() == nil {
				o.log.Error("pass failed", "loop", name, "err", err)
			}
			if err != nil || n < o.cfg.BatchSize {
				break
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// submitPass sends 'signed' rows and resends 'submitted' rows whose receipt is
// overdue, oldest nonce first per key.
func (o *Outbox) submitPass(ctx context.Context) (int, error) {
	return o.claim(ctx, `WHERE status IN ('signed', 'submitted') AND next_attempt_at <= now()
		ORDER BY from_address, nonce`, o.attempt)
}

// receiptPass looks up the receipt of each 'submitted' row, least recently
// checked first.
func (o *Outbox) receiptPass(ctx context.Context) (int, error) {
	return o.claim(ctx, `WHERE status = 'submitted' ORDER BY updated_at`, o.checkReceipt)
}

// claim locks up to BatchSize rows that no other worker holds and handles
// each one in the same database transaction.
func (o *Outbox) claim(ctx context.Context, where string, handle func(context.Context, pgx.Tx, Op) error) (int, error) {
	var n int
	err := pgx.BeginFunc(ctx, o.db, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, selectOp+" "+where+" LIMIT $1 FOR UPDATE SKIP LOCKED", o.cfg.BatchSize)
		if err != nil {
			return err
		}
		ops, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Op, error) { return scanOp(r) })
		if err != nil {
			return err
		}
		n = len(ops)
		for _, op := range ops {
			if err := o.handleOne(ctx, tx, op, handle); err != nil {
				return fmt.Errorf("op %s: %w", op.ID, err)
			}
		}
		return nil
	})
	return n, err
}

// handleOne handles one row inside a savepoint. If it fails (a final hook
// error, say), only that row's changes roll back: it keeps its status, moves
// to the back of the receipt queue, and the rest of the batch still commits.
// An error is returned only if the database transaction itself is broken.
func (o *Outbox) handleOne(ctx context.Context, tx pgx.Tx, op Op, handle func(context.Context, pgx.Tx, Op) error) error {
	err := pgx.BeginFunc(ctx, tx, func(sp pgx.Tx) error { return handle(ctx, sp, op) })
	if err == nil {
		return nil
	}
	if ctx.Err() == nil {
		o.log.Error("row rolled back; will retry", "op", op.ID, "nonce", op.Nonce, "err", err)
	}
	_, err = tx.Exec(ctx, `UPDATE transactions SET updated_at = now() WHERE id = $1`, op.ID)
	return err
}

// attempt sends a row's saved bytes and records what happened. Only database
// errors are returned; send errors become the row's status.
func (o *Outbox) attempt(ctx context.Context, tx pgx.Tx, op Op) error {
	signed := new(types.Transaction)
	if err := signed.UnmarshalBinary(op.raw); err != nil {
		return o.fail(ctx, tx, op, fmt.Sprintf("stored raw_tx does not decode: %v", err))
	}
	sendCtx, cancel := context.WithTimeout(ctx, rpcTimeout)
	_, sendErr := o.node.SendRaw(sendCtx, signed)
	cancel()
	log := o.log.With("op", op.ID, "from", op.From, "nonce", op.Nonce, "tx", op.TxHash)

	switch classify(sendErr) {
	case sendOK, sendKnown:
		_, err := tx.Exec(ctx, `UPDATE transactions SET status = 'submitted', attempts = attempts + 1,
			last_error = NULL, next_attempt_at = now() + $2 * interval '1 millisecond', updated_at = now()
			WHERE id = $1`, op.ID, o.cfg.ResendAfter.Milliseconds())
		if err == nil && op.Status == Submitted {
			log.Info("resent: no receipt yet", "attempt", op.Attempts+1)
		}
		return err

	case sendNonceUsed:
		// Either these bytes were already mined (a resend), or another
		// transaction took the nonce. The receipt tells which.
		rcpt, reason, err := o.receipt(ctx, op.TxHash)
		if err != nil {
			return o.retryLater(ctx, tx, op, err)
		}
		if rcpt != nil {
			return o.finalize(ctx, tx, op, rcpt, reason)
		}
		log.Error("nonce used by another transaction; the key was used outside the outbox")
		return o.fail(ctx, tx, op, fmt.Sprintf("nonce %d was used by another transaction (%v)", op.Nonce, sendErr))

	case sendRejected:
		log.Warn("node rejected transaction", "err", sendErr)
		if err := o.fail(ctx, tx, op, sendErr.Error()); err != nil {
			return err
		}
		if isFiller(op) {
			return nil
		}
		return o.fillGap(ctx, tx, op)

	default:
		log.Warn("send failed; will retry", "err", sendErr, "attempt", op.Attempts+1)
		return o.retryLater(ctx, tx, op, sendErr)
	}
}

// checkReceipt records a submitted row's receipt if there is one.
func (o *Outbox) checkReceipt(ctx context.Context, tx pgx.Tx, op Op) error {
	rcpt, reason, err := o.receipt(ctx, op.TxHash)
	if err != nil {
		o.log.Warn("receipt lookup failed", "op", op.ID, "err", err)
		return nil // try again next pass
	}
	if rcpt != nil {
		return o.finalize(ctx, tx, op, rcpt, reason)
	}
	// Not mined yet. Moving updated_at sends this row to the back of the queue.
	_, err = tx.Exec(ctx, `UPDATE transactions SET updated_at = now() WHERE id = $1`, op.ID)
	return err
}

func (o *Outbox) receipt(ctx context.Context, hash common.Hash) (*types.Receipt, string, error) {
	ctx, cancel := context.WithTimeout(ctx, rpcTimeout)
	defer cancel()
	return o.node.Receipt(ctx, hash)
}

// finalize records a receipt. QBFT never reorganizes, so one receipt is final.
func (o *Outbox) finalize(ctx context.Context, tx pgx.Tx, op Op, rcpt *types.Receipt, reason string) error {
	status := Mined
	if rcpt.Status != types.ReceiptStatusSuccessful {
		status = Reverted
	}
	var contract *string
	if rcpt.ContractAddress != (common.Address{}) {
		contract = hexAddrPtr(&rcpt.ContractAddress)
	}
	var revert *string
	if status == Reverted && reason != "" {
		revert = &reason
	}
	if _, err := tx.Exec(ctx, `UPDATE transactions SET status = $2, block_number = $3,
		contract_address = $4, revert_reason = $5, updated_at = now() WHERE id = $1`,
		op.ID, string(status), rcpt.BlockNumber.Int64(), contract, revert); err != nil {
		return err
	}
	o.log.Info("transaction "+string(status), "op", op.ID, "nonce", op.Nonce, "block", rcpt.BlockNumber)
	return o.final(ctx, tx, op.ID)
}

// fail marks a row as never to be mined.
func (o *Outbox) fail(ctx context.Context, tx pgx.Tx, op Op, reason string) error {
	if _, err := tx.Exec(ctx, `UPDATE transactions SET status = 'failed', attempts = attempts + 1,
		last_error = $2, updated_at = now() WHERE id = $1`, op.ID, reason); err != nil {
		return err
	}
	return o.final(ctx, tx, op.ID)
}

// retryLater keeps the row's status and backs off: 1 s, 2 s, 4 s ... up to MaxBackoff.
func (o *Outbox) retryLater(ctx context.Context, tx pgx.Tx, op Op, cause error) error {
	backoff := o.cfg.MaxBackoff
	if op.Attempts < 16 {
		backoff = min(time.Second<<op.Attempts, o.cfg.MaxBackoff)
	}
	_, err := tx.Exec(ctx, `UPDATE transactions SET attempts = attempts + 1, last_error = $2,
		next_attempt_at = now() + $3 * interval '1 millisecond', updated_at = now() WHERE id = $1`,
		op.ID, cause.Error(), backoff.Milliseconds())
	return err
}

// fillGap signs a zero-value transfer to self at a failed row's nonce. Without
// it, every later transaction from the key would wait forever for that nonce.
func (o *Outbox) fillGap(ctx context.Context, tx pgx.Tx, op Op) error {
	var keyID string
	if err := tx.QueryRow(ctx, `SELECT id FROM signing_keys WHERE address = $1`, hexAddr(op.From)).Scan(&keyID); err != nil {
		return fmt.Errorf("find key for %s: %w", op.From, err)
	}
	signer, err := o.keys.Signer(ctx, tx, keyID)
	if err != nil {
		return err
	}
	filler, err := insert(ctx, tx, signer, op.Nonce, &op.From, nil, fillerGas, "")
	if err != nil {
		return err
	}
	o.log.Warn("filling nonce gap with a transfer to self", "failed", op.ID, "filler", filler.ID, "nonce", op.Nonce)
	return nil
}

func isFiller(op Op) bool { return op.To != nil && *op.To == op.From && len(op.Calldata) == 0 }

func (o *Outbox) final(ctx context.Context, tx pgx.Tx, id string) error {
	if o.cfg.OnFinal == nil {
		return nil
	}
	op, err := Get(ctx, tx, id)
	if err != nil {
		return err
	}
	return o.cfg.OnFinal(ctx, tx, op)
}

type sendResult int

const (
	sendOK        sendResult = iota
	sendKnown                // the node already has these bytes
	sendNonceUsed            // the nonce is taken: by these bytes (mined) or by another transaction
	sendRejected             // the node will never accept these bytes
	sendTransient            // no answer, or a temporary refusal: retry
)

// classify sorts a send error by the action it needs. A JSON-RPC error is
// the node's answer about the transaction; anything else (refused connection,
// timeout, HTTP error) means the outcome is unknown. Messages are Besu's,
// with geth's equivalents.
func classify(err error) sendResult {
	if err == nil {
		return sendOK
	}
	var rpcErr rpc.Error
	if !errors.As(err, &rpcErr) {
		return sendTransient
	}
	msg := strings.ToLower(rpcErr.Error())
	has := func(subs ...string) bool {
		for _, s := range subs {
			if strings.Contains(msg, s) {
				return true
			}
		}
		return false
	}
	switch {
	case has("known transaction", "already known"):
		return sendKnown
	case has("nonce too low"):
		return sendNonceUsed
	case has("pool is full", "txpool is full", "underpriced", "too many"):
		return sendTransient
	default:
		return sendRejected
	}
}
