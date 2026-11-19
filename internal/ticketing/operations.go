package ticketing

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/smlobo/kaleido-challenge/internal/outbox"
)

// Operation is an outbox row with its revert reason decoded.
type Operation struct {
	outbox.Op
	Revert string // e.g. "SoldOut()"; empty unless reverted
}

// GetOperation returns an operation, only to the caller who started it.
// Anyone else gets "not found", so operation IDs reveal nothing.
func (s *Service) GetOperation(ctx context.Context, p Principal, id string) (Operation, error) {
	if p.Role == Anonymous {
		return Operation{}, errNoToken
	}
	if uuid.Validate(id) != nil {
		return Operation{}, errNotFound("operation")
	}
	var by *string
	err := s.pool.QueryRow(ctx, `SELECT requested_by FROM transactions WHERE id = $1`, id).Scan(&by)
	if isNoRows(err) || (err == nil && (by == nil || *by != p.String())) {
		return Operation{}, errNotFound("operation")
	}
	if err != nil {
		return Operation{}, err
	}
	op, err := outbox.Get(ctx, s.pool, id)
	if errors.Is(err, outbox.ErrNotFound) {
		return Operation{}, errNotFound("operation")
	}
	if err != nil {
		return Operation{}, err
	}
	out := Operation{Op: op}
	if op.Status == outbox.Reverted {
		out.Revert = DecodeRevertHex(op.RevertReason)
	}
	return out, nil
}
