package ticketing

import (
	"context"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
	"github.com/smlobo/kaleido-challenge/internal/outbox"
)

// NewEvent is what an organizer asks for.
type NewEvent struct {
	Name         string
	Symbol       string
	StartsAt     time.Time
	MaxSupply    int64
	FaceValue    string // base units, decimal
	MaxResaleBps int64
}

// Event is an event as the API shows it.
type Event struct {
	ID, OrganizerID, Name, Symbol string
	StartsAt                      time.Time
	MaxSupply                     int64
	FaceValue                     string
	MaxResaleBps                  int64
	Status                        string // deploying, on_sale, failed
	ContractAddress               *string
	TicketsIssued, TicketsSold    int64
	DeployOperationID             string
}

const maxUint32 = 1<<32 - 1

var (
	symbolRE   = regexp.MustCompile(`^[A-Z0-9]{1,11}$`)
	uint256RE  = regexp.MustCompile(`^[0-9]{1,78}$`)
	maxUint256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
)

// MinLeadTime is how far ahead an event must start: the contract refuses a
// start in the past, and the deploy takes a block or two to be mined.
const MinLeadTime = time.Minute

func (s *Service) validate(in NewEvent) (*big.Int, error) {
	var bad []string
	if n := len(strings.TrimSpace(in.Name)); n == 0 || n > 100 {
		bad = append(bad, "name must be 1 to 100 characters")
	}
	if !symbolRE.MatchString(in.Symbol) {
		bad = append(bad, "symbol must be 1 to 11 capital letters or digits")
	}
	if in.StartsAt.Before(s.now().Add(MinLeadTime)) {
		bad = append(bad, "starts_at must be at least a minute from now")
	}
	if in.MaxSupply < 1 || in.MaxSupply > maxUint32 {
		bad = append(bad, "max_supply must be 1 to 4294967295")
	}
	if in.MaxResaleBps < 0 || in.MaxResaleBps > maxUint32 {
		bad = append(bad, "max_resale_bps must be 0 to 4294967295")
	}
	face, ok := new(big.Int).SetString(in.FaceValue, 10)
	if !uint256RE.MatchString(in.FaceValue) || !ok || face.Sign() == 0 || face.Cmp(maxUint256) > 0 {
		bad = append(bad, "face_value must be a positive whole number of base units, as a string")
	}
	if len(bad) > 0 {
		return nil, fail(Invalid, "invalid_event", "Invalid event", strings.Join(bad, "; ")+".")
	}
	return face, nil
}

// CreateEvent saves the event and signs its TicketNFT deployment. The
// caller's organizer keys get the contract's roles: box office as minter and
// admin, gate as gate. The event is on_sale once the deploy is mined.
func (s *Service) CreateEvent(ctx context.Context, p Principal, in NewEvent, key string) (Event, error) {
	if err := need(p, Organizer); err != nil {
		return Event{}, err
	}
	face, err := s.validate(in)
	if err != nil {
		return Event{}, err
	}
	idem := idemKey(p, "POST /events", key)
	if id, found, err := s.opByIdemKey(ctx, idem); err != nil || found {
		return s.eventByDeploy(ctx, id, err)
	}

	var boxOffice, gate string
	if err := s.pool.QueryRow(ctx,
		`SELECT b.address, g.address FROM organizers o
		 JOIN signing_keys b ON b.id = o.box_office_key_id
		 JOIN signing_keys g ON g.id = o.gate_key_id
		 WHERE o.id = $1`, p.ID).Scan(&boxOffice, &gate); err != nil {
		return Event{}, err
	}

	eventID := uuid.NewString()
	baseURI := s.baseURL + "/metadata/" + eventID + "/"
	code := append(common.FromHex(ticketnft.TicketNFTMetaData.Bin), s.nft.PackConstructor(
		strings.TrimSpace(in.Name), in.Symbol, uint32(in.MaxSupply), uint32(in.MaxResaleBps),
		uint64(in.StartsAt.Unix()), baseURI, common.HexToAddress(boxOffice), common.HexToAddress(gate))...)
	if err := s.precheck(ctx, s.deployer.Address, nil, code); err != nil {
		return Event{}, err
	}

	var op outbox.Op
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx,
			`INSERT INTO events (id, organizer_id, name, symbol, starts_at, max_supply, face_value, max_resale_bps)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			eventID, p.ID, strings.TrimSpace(in.Name), in.Symbol, in.StartsAt, in.MaxSupply, face.String(), in.MaxResaleBps); err != nil {
			return err
		}
		op, err = s.ob.Enqueue(ctx, tx, outbox.Request{KeyID: s.deployer.ID, Data: code, Gas: gasDeploy, IdempotencyKey: idem})
		if err != nil {
			return err
		}
		if op.Existing {
			return errRepeat
		}
		if _, err := tx.Exec(ctx, `UPDATE transactions SET requested_by = $2 WHERE id = $1`, op.ID, p.String()); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE events SET deploy_tx_id = $2 WHERE id = $1`, eventID, op.ID)
		return err
	})
	if isRepeat(err) {
		id, _, err := s.opByIdemKey(ctx, idem)
		return s.eventByDeploy(ctx, id, err)
	}
	if err != nil {
		return Event{}, err
	}
	s.send(ctx, op.ID)
	return s.GetEvent(ctx, eventID)
}

func (s *Service) eventByDeploy(ctx context.Context, opID string, err error) (Event, error) {
	if err != nil {
		return Event{}, err
	}
	var id string
	if err := s.pool.QueryRow(ctx, `SELECT id FROM events WHERE deploy_tx_id = $1`, opID).Scan(&id); err != nil {
		return Event{}, err
	}
	return s.GetEvent(ctx, id)
}

// GetEvent returns one event.
func (s *Service) GetEvent(ctx context.Context, id string) (Event, error) {
	if uuid.Validate(id) != nil {
		return Event{}, errNotFound("event")
	}
	var e Event
	err := s.pool.QueryRow(ctx,
		`SELECT e.id, e.organizer_id, e.name, e.symbol, e.starts_at, e.max_supply, e.face_value::text,
		        e.max_resale_bps, e.status, e.contract_address, e.tickets_issued, e.deploy_tx_id,
		        (SELECT count(*) FROM tickets t WHERE t.event_id = e.id AND t.status <> 'failed')
		 FROM events e WHERE e.id = $1`, id).Scan(
		&e.ID, &e.OrganizerID, &e.Name, &e.Symbol, &e.StartsAt, &e.MaxSupply, &e.FaceValue,
		&e.MaxResaleBps, &e.Status, &e.ContractAddress, &e.TicketsIssued, &e.DeployOperationID, &e.TicketsSold)
	if isNoRows(err) {
		return Event{}, errNotFound("event")
	}
	return e, err
}
