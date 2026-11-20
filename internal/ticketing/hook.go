package ticketing

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"

	"github.com/smlobo/kaleido-challenge/internal/outbox"
)

// OnFinal is the outbox's final-status hook. It runs in the same database
// transaction that records the operation as mined, reverted or failed, and
// moves the event or ticket the operation belongs to:
//
//	deploy   mined → event on_sale with its contract address; else failed
//	mint     mined → ticket active, owned by the buyer's wallet; else failed
//	check-in mined → ticket checked_in, listing cleared; else stays active
//	list     mined → ticket's listing_price set
//	cancel   mined → listing_price cleared
//	buy      mined → ticket owned by the buyer's wallet, listing cleared
//
// A reverted or failed listing, cancellation or purchase changes nothing:
// the ticket is as the chain left it.
//
// Each update only moves a row out of the state it expects, so running it
// twice changes nothing, and operations that are neither (fillers) match no
// row. Until the indexer (milestone 8) exists, this is where owner_address
// comes from: for a mined mint it is exactly what the chain says.
func OnFinal(ctx context.Context, tx pgx.Tx, op outbox.Op) error {
	mined := op.Status == outbox.Mined
	var contract *string
	if op.ContractAddress != nil {
		a := strings.ToLower(op.ContractAddress.Hex())
		contract = &a
	}
	if _, err := tx.Exec(ctx,
		`UPDATE events SET status = CASE WHEN $2 THEN 'on_sale' ELSE 'failed' END,
		        contract_address = CASE WHEN $2 THEN $3 END
		 WHERE deploy_tx_id = $1 AND status = 'deploying'`, op.ID, mined, contract); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`UPDATE tickets t SET status = CASE WHEN $2 THEN 'active' ELSE 'failed' END,
		        owner_address = CASE WHEN $2 THEN (
		            SELECT k.address FROM users u JOIN signing_keys k ON k.id = u.wallet_key_id
		            WHERE u.id = t.buyer_user_id) END
		 WHERE mint_tx_id = $1 AND status = 'pending_mint'`, op.ID, mined); err != nil {
		return err
	}
	if mined {
		if _, err := tx.Exec(ctx,
			`UPDATE tickets SET status = 'checked_in', listing_price = NULL
			 WHERE checkin_tx_id = $1 AND status = 'active'`, op.ID); err != nil {
			return err
		}
		// The contract clears a listing when the ticket is sold, so a sale
		// does both: new owner, no listing.
		if _, err := tx.Exec(ctx,
			`UPDATE tickets t SET
			     listing_price = CASE a.kind WHEN 'list' THEN a.price END,
			     owner_address = CASE a.kind WHEN 'buy' THEN (
			         SELECT k.address FROM users u JOIN signing_keys k ON k.id = u.wallet_key_id
			         WHERE u.id = a.user_id) ELSE t.owner_address END
			 FROM resale_actions a
			 WHERE a.tx_id = $1 AND t.id = a.ticket_id AND t.status = 'active'`, op.ID); err != nil {
			return err
		}
	}
	return nil
}

func keccak(b []byte) common.Hash { return crypto.Keccak256Hash(b) }
