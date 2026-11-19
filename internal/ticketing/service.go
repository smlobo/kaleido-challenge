// Package ticketing holds the business rules behind the API: organizers,
// users, events, tickets and check-in. It knows nothing about HTTP. Every
// write that touches the chain follows the same path: check the rules
// off-chain, dry-run the call on the node, then sign and save it through the
// outbox in one database transaction with the row it belongs to.
package ticketing

import (
	"context"
	"crypto/sha256"
	"errors"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
	"github.com/smlobo/kaleido-challenge/internal/outbox"
)

// Fixed gas limits per call, as the design doc says. Gas is free on this
// chain; these leave room above what each call uses.
const (
	gasDeploy  = 6_000_000
	gasMint    = 300_000
	gasCheckIn = 200_000
)

// Estimator dry-runs a call on the node. *chain.Client satisfies it.
type Estimator interface {
	EstimateGas(ctx context.Context, from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error)
}

// Config is what the service needs.
type Config struct {
	Pool          *pgxpool.Pool
	Keys          *keystore.Store
	Outbox        *outbox.Outbox
	Node          Estimator
	Deployer      keystore.Key
	ChainID       *big.Int // the EIP-712 domain of check-in passes
	PublicBaseURL string   // no trailing slash
	AdminToken    string   // empty: POST /organizers is refused
	Logger        *slog.Logger
	Now           func() time.Time // tests only
}

// Service carries out the API's requests.
type Service struct {
	pool      *pgxpool.Pool
	keys      *keystore.Store
	ob        *outbox.Outbox
	node      Estimator
	deployer  keystore.Key
	chainID   *big.Int
	baseURL   string
	adminHash []byte // SHA-256 of ADMIN_TOKEN, or nil
	log       *slog.Logger
	now       func() time.Time
	nft       *ticketnft.TicketNFT
}

// New returns a Service.
func New(cfg Config) *Service {
	s := &Service{
		pool:     cfg.Pool,
		keys:     cfg.Keys,
		ob:       cfg.Outbox,
		node:     cfg.Node,
		deployer: cfg.Deployer,
		chainID:  cfg.ChainID,
		baseURL:  strings.TrimRight(cfg.PublicBaseURL, "/"),
		log:      cfg.Logger,
		now:      cfg.Now,
		nft:      ticketnft.NewTicketNFT(),
	}
	if cfg.AdminToken != "" {
		h := sha256.Sum256([]byte(cfg.AdminToken))
		s.adminHash = h[:]
	}
	if s.log == nil {
		s.log = slog.Default()
	}
	if s.now == nil {
		s.now = time.Now
	}
	return s
}

// Role is the kind of caller.
type Role string

const (
	Anonymous Role = ""
	Admin     Role = "admin"
	Organizer Role = "organizer"
	User      Role = "user"
)

// Principal is who is calling, resolved from their token. Handlers only
// ever see this, never the token, so the login method can change.
type Principal struct {
	Role Role
	ID   string // organizers.id or users.id; empty for Admin
}

func (p Principal) String() string { return string(p.Role) + ":" + p.ID }

// Kind says what went wrong, so the API can pick a status code.
type Kind int

const (
	Invalid      Kind = iota // the request breaks a rule: 400 or 422
	Unauthorized             // no valid token: 401
	Forbidden                // valid token, wrong caller: 403
	NotFound                 // 404
	Conflict                 // the current state does not allow it, maybe yet: 409
	Rejected                 // a business rule or the contract says no: 422
)

// Error is a refusal the caller can act on. Code is stable; Detail is for
// people.
type Error struct {
	Kind   Kind
	Code   string
	Title  string
	Detail string
}

func (e *Error) Error() string { return e.Code + ": " + e.Detail }

func fail(kind Kind, code, title, detail string) *Error {
	return &Error{Kind: kind, Code: code, Title: title, Detail: detail}
}

var (
	errNotFound  = func(what string) *Error { return fail(NotFound, "not_found", "Not found", "No "+what+" with that ID.") }
	errNeedsRole = func(r Role) *Error {
		return fail(Forbidden, "forbidden", "Forbidden", "This needs a "+string(r)+" token.")
	}
	errNoToken    = fail(Unauthorized, "unauthorized", "Unauthorized", "Send Authorization: Bearer <token>.")
	errBadToken   = fail(Unauthorized, "unauthorized", "Unauthorized", "The token is not valid.")
	errNoAdmin    = fail(Forbidden, "forbidden", "Forbidden", "Creating organizers needs the platform's ADMIN_TOKEN.")
	errAdminUnset = fail(Forbidden, "admin_disabled", "Admin disabled", "ADMIN_TOKEN is not set on this service, so organizers cannot be created.")
)

func need(p Principal, r Role) error {
	switch {
	case p.Role == Anonymous:
		return errNoToken
	case p.Role != r:
		return errNeedsRole(r)
	}
	return nil
}

// idemKey scopes a client's Idempotency-Key to the caller and the request,
// so two callers, or one caller on two endpoints, never share a key.
func idemKey(p Principal, request, key string) string {
	return p.String() + " " + request + " " + key
}

// opByIdemKey returns the operation an earlier request with this key created.
func (s *Service) opByIdemKey(ctx context.Context, key string) (string, bool, error) {
	var id string
	err := s.pool.QueryRow(ctx, `SELECT id FROM transactions WHERE idempotency_key = $1`, key).Scan(&id)
	if err != nil {
		if isNoRows(err) {
			return "", false, nil
		}
		return "", false, err
	}
	return id, true, nil
}

// errRepeat aborts a database transaction when Enqueue finds the request's
// idempotency key already used, so the caller can return the first result.
var errRepeat = errors.New("repeated request")

// isRepeat reports whether err means another request with the same
// idempotency key got there first.
func isRepeat(err error) bool {
	return errors.Is(err, errRepeat) || isUnique(err, "transactions_idempotency_key_key")
}

// precheck dry-runs a call with eth_estimateGas. A revert becomes a refusal
// with the decoded contract error, before anything is signed. If the node
// cannot be reached, the call is signed anyway: the outbox sends it once the
// node is back.
func (s *Service) precheck(ctx context.Context, from common.Address, to *common.Address, data []byte) error {
	if s.node == nil {
		return nil
	}
	_, err := s.node.EstimateGas(ctx, from, to, nil, data)
	if err == nil {
		return nil
	}
	if raw, ok := ethclient.RevertErrorData(err); ok {
		return contractRefusal(DecodeRevert(raw))
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	s.log.Warn("could not dry-run the call; saving it for the outbox to send", "err", err)
	return nil
}

func contractRefusal(reason string) *Error {
	switch {
	case strings.HasPrefix(reason, "SoldOut"):
		return fail(Rejected, "sold_out", "Sold out", "Every ticket for this event has been issued.")
	case strings.HasPrefix(reason, "EventStarted"):
		return fail(Rejected, "event_started", "Event started", "The event has started.")
	case strings.HasPrefix(reason, "PassExpired"):
		return errPassExpired
	case strings.HasPrefix(reason, "NotSignedByHolder"):
		return fail(Rejected, "pass_not_from_holder", "Pass not from the holder",
			"The pass was not signed by the ticket's current holder, or it was made for another ticket or event. Ask the holder for a new pass.")
	}
	return fail(Rejected, "rejected_by_contract", "Rejected by the contract", "The contract would revert with "+reason+".")
}

// send is the fast path after a commit. A failure is logged, not returned:
// the row is saved, and the submitter loop sends it.
func (s *Service) send(ctx context.Context, opID string) {
	if s.ob == nil {
		return
	}
	if err := s.ob.Send(context.WithoutCancel(ctx), opID); err != nil {
		s.log.Warn("send after commit failed; the submitter will retry", "op", opID, "err", err)
	}
}

func isNoRows(err error) bool { return errors.Is(err, pgx.ErrNoRows) }

func isUnique(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}

func hexAddr(a common.Address) string { return strings.ToLower(a.Hex()) }
