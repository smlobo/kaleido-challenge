package ticketing

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/smlobo/kaleido-challenge/internal/outbox"
)

// Ticket is a ticket as the API shows it.
type Ticket struct {
	ID, EventID, TokenID string
	Status               string // pending_mint, active, checked_in, failed
	OwnerAddress         *string
	ListingPrice         *string         // set while listed for resale
	Record               json.RawMessage // canonical form
	RecordHash           string
	MetadataURL          string
	MintOperationID      string
	CheckInOperationID   *string
}

// Metadata is the ERC-721 metadata the contract's tokenURI points at.
type Metadata struct {
	Name, Description string
	Attributes        [][2]string // trait, value
	Record            json.RawMessage
	RecordHash        string
}

// BuyTicket hands the caller the event's next ticket number and signs a
// mint of it to their wallet, sent by the organizer's box office key.
func (s *Service) BuyTicket(ctx context.Context, p Principal, eventID, key string) (Ticket, error) {
	if err := need(p, User); err != nil {
		return Ticket{}, err
	}
	if uuid.Validate(eventID) != nil {
		return Ticket{}, errNotFound("event")
	}
	idem := idemKey(p, "POST /events/"+eventID+"/tickets", key)
	if id, found, err := s.opByIdemKey(ctx, idem); err != nil || found {
		return s.ticketByOp(ctx, "mint_tx_id", id, err)
	}

	var (
		status, faceValue, boxOfficeID, boxOffice string
		contract                                  *string
		startsAt                                  time.Time
		issued, maxSupply                         int64
		wallet                                    string
	)
	err := s.pool.QueryRow(ctx,
		`SELECT e.status, e.contract_address, e.starts_at, e.tickets_issued, e.max_supply, e.face_value::text,
		        o.box_office_key_id, k.address
		 FROM events e JOIN organizers o ON o.id = e.organizer_id JOIN signing_keys k ON k.id = o.box_office_key_id
		 WHERE e.id = $1`, eventID).Scan(&status, &contract, &startsAt, &issued, &maxSupply, &faceValue, &boxOfficeID, &boxOffice)
	if isNoRows(err) {
		return Ticket{}, errNotFound("event")
	}
	if err != nil {
		return Ticket{}, err
	}
	switch {
	case status == "deploying":
		return Ticket{}, fail(Conflict, "event_not_ready", "Event not ready", "The event's contract is still being deployed; try again in a few seconds.")
	case status == "failed":
		return Ticket{}, fail(Rejected, "event_failed", "Event failed", "The event's contract deployment failed, so it has no tickets.")
	case !s.now().Before(startsAt):
		return Ticket{}, fail(Rejected, "event_started", "Event started", "Tickets are not sold after the event starts.")
	case issued >= maxSupply:
		return Ticket{}, fail(Rejected, "sold_out", "Sold out", "Every ticket for this event has been issued.")
	}
	if err := s.pool.QueryRow(ctx,
		`SELECT k.address FROM users u JOIN signing_keys k ON k.id = u.wallet_key_id WHERE u.id = $1`, p.ID).Scan(&wallet); err != nil {
		return Ticket{}, err
	}
	to := common.HexToAddress(*contract)
	face, _ := new(big.Int).SetString(faceValue, 10)

	// Dry-run a mint of the next number before taking it, so the event row
	// is not locked during a call to the node.
	probe := newRecord(eventID, issued+1, faceValue, s.now())
	if err := s.precheck(ctx, common.HexToAddress(boxOffice), &to,
		s.nft.PackMintTicket(common.HexToAddress(wallet), big.NewInt(issued+1), face, keccak(probe))); err != nil {
		return Ticket{}, err
	}

	var op outbox.Op
	var ticketID string
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		// Take the next number; this also enforces the supply under concurrency.
		var number int64
		err := tx.QueryRow(ctx,
			`UPDATE events SET tickets_issued = tickets_issued + 1
			 WHERE id = $1 AND tickets_issued < max_supply RETURNING tickets_issued`, eventID).Scan(&number)
		if isNoRows(err) {
			return fail(Rejected, "sold_out", "Sold out", "Every ticket for this event has been issued.")
		}
		if err != nil {
			return err
		}
		record := newRecord(eventID, number, faceValue, s.now())
		hash := keccak(record)
		op, err = s.ob.Enqueue(ctx, tx, outbox.Request{KeyID: boxOfficeID, To: &to, Gas: gasMint, IdempotencyKey: idem,
			Data: s.nft.PackMintTicket(common.HexToAddress(wallet), big.NewInt(number), face, hash)})
		if err != nil {
			return err
		}
		if op.Existing {
			return errRepeat
		}
		if _, err := tx.Exec(ctx, `UPDATE transactions SET requested_by = $2 WHERE id = $1`, op.ID, p.String()); err != nil {
			return err
		}
		return tx.QueryRow(ctx,
			`INSERT INTO tickets (event_id, token_id, buyer_user_id, record, record_hash, mint_tx_id)
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			eventID, number, p.ID, string(record), hash.Hex(), op.ID).Scan(&ticketID)
	})
	if isRepeat(err) {
		id, _, err := s.opByIdemKey(ctx, idem)
		return s.ticketByOp(ctx, "mint_tx_id", id, err)
	}
	if err != nil {
		return Ticket{}, err
	}
	s.send(ctx, op.ID)
	return s.GetTicket(ctx, ticketID)
}

// CheckIn signs a checkIn (burn) of the ticket with the organizer's gate key,
// carrying the holder's pass (from CreatePass). Only the event's organizer
// may, and only once; the contract also requires the holder's signature.
func (s *Service) CheckIn(ctx context.Context, p Principal, ticketID string, pass PassInput, key string) (Ticket, error) {
	if err := need(p, Organizer); err != nil {
		return Ticket{}, err
	}
	if uuid.Validate(ticketID) != nil {
		return Ticket{}, errNotFound("ticket")
	}
	idem := idemKey(p, "POST /tickets/"+ticketID+"/check-in", key)
	if id, found, err := s.opByIdemKey(ctx, idem); err != nil || found {
		return s.ticketByOp(ctx, "checkin_tx_id", id, err)
	}
	// After the repeat check: a retry of a check-in already saved must not
	// fail just because its pass has since expired.
	deadline, sig, err := s.parsePass(ticketID, pass)
	if err != nil {
		return Ticket{}, err
	}

	var op outbox.Op
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		// Lock the ticket, so two check-ins of it cannot both pass the checks.
		var (
			organizerID, status, tokenID, gateID, gate, contract string
			checkinStatus                                        *string
			reselling                                            bool
		)
		err := tx.QueryRow(ctx,
			`SELECT e.organizer_id, t.status, t.token_id::text, o.gate_key_id, k.address, e.contract_address, c.status,
			        EXISTS (SELECT 1 FROM resale_actions a JOIN transactions x ON x.id = a.tx_id
			                WHERE a.ticket_id = t.id AND x.status IN ('signed', 'submitted'))
			 FROM tickets t
			 JOIN events e ON e.id = t.event_id
			 JOIN organizers o ON o.id = e.organizer_id
			 JOIN signing_keys k ON k.id = o.gate_key_id
			 LEFT JOIN transactions c ON c.id = t.checkin_tx_id
			 WHERE t.id = $1 FOR UPDATE OF t`, ticketID).Scan(&organizerID, &status, &tokenID, &gateID, &gate, &contract,
			&checkinStatus, &reselling)
		if isNoRows(err) {
			return errNotFound("ticket")
		}
		if err != nil {
			return err
		}
		switch {
		case organizerID != p.ID:
			return fail(Forbidden, "not_your_event", "Not your event", "Only the event's organizer can check its tickets in.")
		case status == "pending_mint":
			return errNotMintedYet
		case status == "failed":
			return errTicketFailed
		case status == "checked_in":
			return errAlreadyCheckedIn
		case checkinStatus != nil && *checkinStatus != string(outbox.Reverted) && *checkinStatus != string(outbox.Failed):
			return fail(Conflict, "already_checked_in", "Already checked in", "A check-in of this ticket is already on its way.")
		case reselling:
			return errActionInFlight
		}

		to := common.HexToAddress(contract)
		token, _ := new(big.Int).SetString(tokenID, 10)
		data := s.nft.PackCheckIn(token, deadline, sig)
		if err := s.precheck(ctx, common.HexToAddress(gate), &to, data); err != nil {
			return err
		}
		op, err = s.ob.Enqueue(ctx, tx, outbox.Request{KeyID: gateID, To: &to, Gas: gasCheckIn, Data: data, IdempotencyKey: idem})
		if err != nil {
			return err
		}
		if op.Existing {
			return errRepeat
		}
		if _, err := tx.Exec(ctx, `UPDATE transactions SET requested_by = $2 WHERE id = $1`, op.ID, p.String()); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE tickets SET checkin_tx_id = $2 WHERE id = $1`, ticketID, op.ID)
		return err
	})
	if isRepeat(err) {
		id, _, err := s.opByIdemKey(ctx, idem)
		return s.ticketByOp(ctx, "checkin_tx_id", id, err)
	}
	if err != nil {
		return Ticket{}, err
	}
	s.send(ctx, op.ID)
	return s.GetTicket(ctx, ticketID)
}

var (
	errNotMintedYet     = fail(Conflict, "not_minted_yet", "Not minted yet", "The ticket's mint is not mined yet; try again in a few seconds.")
	errTicketFailed     = fail(Rejected, "ticket_failed", "Ticket failed", "The ticket's mint failed, so there is nothing to check in.")
	errAlreadyCheckedIn = fail(Conflict, "already_checked_in", "Already checked in", "This ticket has already been checked in.")
)

func (s *Service) ticketByOp(ctx context.Context, column, opID string, err error) (Ticket, error) {
	if err != nil {
		return Ticket{}, err
	}
	var id string
	// column is one of two constants, never input.
	if err := s.pool.QueryRow(ctx, `SELECT id FROM tickets WHERE `+column+` = $1`, opID).Scan(&id); err != nil {
		return Ticket{}, err
	}
	return s.GetTicket(ctx, id)
}

const selectTicket = `SELECT id, event_id, token_id::text, status, owner_address, listing_price::text,
	record::text, record_hash, mint_tx_id, checkin_tx_id FROM tickets`

func (s *Service) scanTicket(row pgx.Row) (Ticket, error) {
	var t Ticket
	var record string
	if err := row.Scan(&t.ID, &t.EventID, &t.TokenID, &t.Status, &t.OwnerAddress, &t.ListingPrice, &record, &t.RecordHash,
		&t.MintOperationID, &t.CheckInOperationID); err != nil {
		if isNoRows(err) {
			return Ticket{}, errNotFound("ticket")
		}
		return Ticket{}, err
	}
	c, err := Canonical([]byte(record))
	if err != nil {
		return Ticket{}, err
	}
	t.Record = c
	t.MetadataURL = s.baseURL + "/metadata/" + t.EventID + "/" + t.TokenID
	return t, nil
}

// GetTicket returns one ticket.
func (s *Service) GetTicket(ctx context.Context, id string) (Ticket, error) {
	if uuid.Validate(id) != nil {
		return Ticket{}, errNotFound("ticket")
	}
	return s.scanTicket(s.pool.QueryRow(ctx, selectTicket+` WHERE id = $1`, id))
}

// GetMetadata returns a ticket's ERC-721 metadata, by event and token ID.
func (s *Service) GetMetadata(ctx context.Context, eventID, tokenID string) (Metadata, error) {
	if uuid.Validate(eventID) != nil || !uint256RE.MatchString(tokenID) {
		return Metadata{}, errNotFound("ticket")
	}
	t, err := s.scanTicket(s.pool.QueryRow(ctx, selectTicket+` WHERE event_id = $1 AND token_id = $2::numeric`, eventID, tokenID))
	if err != nil {
		return Metadata{}, err
	}
	e, err := s.GetEvent(ctx, eventID)
	if err != nil {
		return Metadata{}, err
	}
	return Metadata{
		Name:        e.Name + " #" + t.TokenID,
		Description: "Ticket " + t.TokenID + " for " + e.Name + ", starting " + e.StartsAt.UTC().Format(time.RFC1123) + ".",
		Attributes: [][2]string{
			{"Event", e.Name},
			{"Starts", e.StartsAt.UTC().Format(time.RFC3339)},
			{"Face value", e.FaceValue},
			{"Status", t.Status},
		},
		Record:     t.Record,
		RecordHash: t.RecordHash,
	}, nil
}
