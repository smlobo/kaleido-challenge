package ticketing

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/smlobo/kaleido-challenge/internal/outbox"
)

// Resale is milestone 7. A holder lists a ticket at no more than the event's
// resale cap, can cancel the listing, and another user buys it. Each is one
// transaction from the user's own custodial wallet: the contract checks that
// the lister is the holder (msg.sender) and moves the ticket to the buyer.
// No money moves yet; the contract records the price (TicketResold).
//
// A ticket has at most one chain action in flight at a time (a listing,
// cancellation, purchase or check-in). The final-status hook then applies
// mined actions to the ticket in the order they happened on-chain.

// ResaleAction is the answer to a listing, cancellation or purchase: the
// ticket as it is now, and the operation to follow.
type ResaleAction struct {
	OperationID string
	Ticket      Ticket
}

// Listing is a ticket on the resale market.
type Listing struct {
	TicketID, TokenID, Price, Seller string
}

// resaleTicket is what the resale checks need, read with the ticket locked.
type resaleTicket struct {
	status, tokenID, faceValue    string
	owner, listingPrice, contract *string
	startsAt                      time.Time
	maxResaleBps                  int64
	busy                          bool // a chain action of the ticket is in flight
}

// lockForResale locks the ticket's row for the rest of tx and reads it.
func lockForResale(ctx context.Context, tx pgx.Tx, ticketID string) (resaleTicket, error) {
	var r resaleTicket
	err := tx.QueryRow(ctx,
		`SELECT t.status, t.token_id::text, t.owner_address, t.listing_price::text, e.contract_address,
		        e.starts_at, e.face_value::text, e.max_resale_bps,
		        EXISTS (SELECT 1 FROM resale_actions a JOIN transactions x ON x.id = a.tx_id
		                WHERE a.ticket_id = t.id AND x.status IN ('signed', 'submitted'))
		        OR EXISTS (SELECT 1 FROM transactions c
		                WHERE c.id = t.checkin_tx_id AND c.status IN ('signed', 'submitted'))
		 FROM tickets t JOIN events e ON e.id = t.event_id
		 WHERE t.id = $1 FOR UPDATE OF t`, ticketID).Scan(&r.status, &r.tokenID, &r.owner, &r.listingPrice,
		&r.contract, &r.startsAt, &r.faceValue, &r.maxResaleBps, &r.busy)
	if isNoRows(err) {
		return r, errNotFound("ticket")
	}
	return r, err
}

// tradable refuses a ticket that cannot be listed, unlisted or bought now.
func (s *Service) tradable(r resaleTicket) error {
	switch {
	case r.status == "pending_mint":
		return errNotMintedYet
	case r.status == "failed":
		return errTicketFailed
	case r.status == "checked_in":
		return errAlreadyCheckedIn
	case !s.now().Before(r.startsAt):
		return fail(Rejected, "event_started", "Event started", "Tickets cannot be listed or resold after the event starts.")
	case r.busy:
		return errActionInFlight
	}
	return nil
}

// cap is the highest price the ticket may be listed at: face value times
// the event's resale cap, rounded down as the contract's Math.mulDiv does.
func (r resaleTicket) cap() *big.Int {
	face, _ := new(big.Int).SetString(r.faceValue, 10)
	c := new(big.Int).Mul(face, big.NewInt(r.maxResaleBps))
	return c.Quo(c, big.NewInt(10_000))
}

var (
	errActionInFlight = fail(Conflict, "action_in_progress", "Action in progress",
		"Another transaction for this ticket is still on its way; try again once it is mined.")
	errNotYourTicket = fail(Forbidden, "not_your_ticket", "Not your ticket", "Only the ticket's holder can do that.")
	errNotListed     = fail(Conflict, "not_listed", "Not listed", "The ticket is not listed for resale.")
)

// ListTicket lists one of the caller's tickets for resale at price, or
// changes the price of an existing listing.
func (s *Service) ListTicket(ctx context.Context, p Principal, ticketID, price, key string) (ResaleAction, error) {
	if err := need(p, User); err != nil {
		return ResaleAction{}, err
	}
	amount, err := parsePrice("price", price)
	if err != nil {
		return ResaleAction{}, err
	}
	return s.resale(ctx, p, ticketID, "PUT /tickets/"+ticketID+"/listing", key, "list",
		func(r resaleTicket, wallet string) ([]byte, *big.Int, error) {
			if r.owner == nil || !strings.EqualFold(*r.owner, wallet) {
				return nil, nil, errNotYourTicket
			}
			if limit := r.cap(); amount.Cmp(limit) > 0 {
				return nil, nil, fail(Rejected, "price_above_cap", "Price above the cap",
					"This event's tickets can be listed for at most "+limit.String()+".")
			}
			return s.nft.PackListForResale(tokenNumber(r), amount), amount, nil
		})
}

// CancelListing takes one of the caller's tickets off the resale market.
func (s *Service) CancelListing(ctx context.Context, p Principal, ticketID, key string) (ResaleAction, error) {
	if err := need(p, User); err != nil {
		return ResaleAction{}, err
	}
	return s.resale(ctx, p, ticketID, "DELETE /tickets/"+ticketID+"/listing", key, "cancel",
		func(r resaleTicket, wallet string) ([]byte, *big.Int, error) {
			if r.owner == nil || !strings.EqualFold(*r.owner, wallet) {
				return nil, nil, errNotYourTicket
			}
			if r.listingPrice == nil {
				return nil, nil, errNotListed
			}
			return s.nft.PackCancelListing(tokenNumber(r)), nil, nil
		})
}

// BuyResale buys a listed ticket for the caller. maxPrice is the most the
// buyer agrees to pay; empty means the price listed now. If the seller
// raises the price before the purchase is mined, the contract refuses it.
func (s *Service) BuyResale(ctx context.Context, p Principal, ticketID, maxPrice, key string) (ResaleAction, error) {
	if err := need(p, User); err != nil {
		return ResaleAction{}, err
	}
	var limit *big.Int
	if maxPrice != "" {
		var err error
		if limit, err = parsePrice("max_price", maxPrice); err != nil {
			return ResaleAction{}, err
		}
	}
	return s.resale(ctx, p, ticketID, "POST /tickets/"+ticketID+"/purchase", key, "buy",
		func(r resaleTicket, wallet string) ([]byte, *big.Int, error) {
			if r.listingPrice == nil {
				return nil, nil, errNotListed
			}
			if r.owner != nil && strings.EqualFold(*r.owner, wallet) {
				return nil, nil, fail(Rejected, "buyer_is_seller", "Buyer is the seller", "You already hold this ticket.")
			}
			listed, _ := new(big.Int).SetString(*r.listingPrice, 10)
			max := limit
			if max == nil {
				max = listed
			}
			if listed.Cmp(max) > 0 {
				return nil, nil, fail(Rejected, "price_changed", "Price changed",
					"The ticket is listed at "+listed.String()+", more than max_price.")
			}
			return s.nft.PackBuyResale(tokenNumber(r), max), max, nil
		})
}

// resale is what listing, cancelling and buying share: the repeat check,
// the locked ticket, the dry run, and saving the transaction with its
// resale_actions row. call checks the request against the ticket and
// returns the contract call to send from the caller's wallet, and the price
// to store with the action (nil for a cancellation).
func (s *Service) resale(ctx context.Context, p Principal, ticketID, request, key, kind string,
	call func(r resaleTicket, wallet string) ([]byte, *big.Int, error)) (ResaleAction, error) {
	if uuid.Validate(ticketID) != nil {
		return ResaleAction{}, errNotFound("ticket")
	}
	idem := idemKey(p, request, key)
	if id, found, err := s.opByIdemKey(ctx, idem); err != nil || found {
		return s.resaleByOp(ctx, id, err)
	}
	var walletID, wallet string
	if err := s.pool.QueryRow(ctx,
		`SELECT k.id, k.address FROM users u JOIN signing_keys k ON k.id = u.wallet_key_id WHERE u.id = $1`,
		p.ID).Scan(&walletID, &wallet); err != nil {
		return ResaleAction{}, err
	}

	var op outbox.Op
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		r, err := lockForResale(ctx, tx, ticketID)
		if err != nil {
			return err
		}
		if err := s.tradable(r); err != nil {
			return err
		}
		data, price, err := call(r, wallet)
		if err != nil {
			return err
		}
		to := common.HexToAddress(*r.contract)
		if err := s.precheck(ctx, common.HexToAddress(wallet), &to, data); err != nil {
			return err
		}
		op, err = s.ob.Enqueue(ctx, tx, outbox.Request{KeyID: walletID, To: &to, Gas: gasResale, Data: data, IdempotencyKey: idem})
		if err != nil {
			return err
		}
		if op.Existing {
			return errRepeat
		}
		if _, err := tx.Exec(ctx, `UPDATE transactions SET requested_by = $2 WHERE id = $1`, op.ID, p.String()); err != nil {
			return err
		}
		var stored *string
		if price != nil {
			v := price.String()
			stored = &v
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO resale_actions (tx_id, ticket_id, kind, user_id, price) VALUES ($1, $2, $3, $4, $5)`,
			op.ID, ticketID, kind, p.ID, stored)
		return err
	})
	if isRepeat(err) {
		id, _, err := s.opByIdemKey(ctx, idem)
		return s.resaleByOp(ctx, id, err)
	}
	if err != nil {
		return ResaleAction{}, err
	}
	s.send(ctx, op.ID)
	t, err := s.GetTicket(ctx, ticketID)
	return ResaleAction{OperationID: op.ID, Ticket: t}, err
}

func (s *Service) resaleByOp(ctx context.Context, opID string, err error) (ResaleAction, error) {
	if err != nil {
		return ResaleAction{}, err
	}
	var ticketID string
	if err := s.pool.QueryRow(ctx, `SELECT ticket_id FROM resale_actions WHERE tx_id = $1`, opID).Scan(&ticketID); err != nil {
		return ResaleAction{}, err
	}
	t, err := s.GetTicket(ctx, ticketID)
	return ResaleAction{OperationID: opID, Ticket: t}, err
}

// Listings returns an event's tickets on the resale market, cheapest first.
// A listing shows once its transaction is mined.
func (s *Service) Listings(ctx context.Context, eventID string) ([]Listing, error) {
	if uuid.Validate(eventID) != nil {
		return nil, errNotFound("event")
	}
	var exists bool
	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, errNotFound("event")
	}
	rows, err := s.pool.Query(ctx,
		`SELECT id, token_id::text, listing_price::text, owner_address FROM tickets
		 WHERE event_id = $1 AND listing_price IS NOT NULL AND status = 'active'
		 ORDER BY listing_price, token_id`, eventID)
	if err != nil {
		return nil, err
	}
	listings, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Listing, error) {
		var l Listing
		err := row.Scan(&l.TicketID, &l.TokenID, &l.Price, &l.Seller)
		return l, err
	})
	if listings == nil {
		listings = []Listing{}
	}
	return listings, err
}

// parsePrice reads a positive whole number of base units, at most 2^256 - 1.
func parsePrice(field, v string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(v, 10)
	if !uint256RE.MatchString(v) || !ok || n.Sign() == 0 || n.Cmp(maxUint256) > 0 {
		return nil, fail(Invalid, "invalid_price", "Invalid price",
			field+" must be a positive whole number of base units, as a string.")
	}
	return n, nil
}

func tokenNumber(r resaleTicket) *big.Int {
	n, _ := new(big.Int).SetString(r.tokenID, 10)
	return n
}
