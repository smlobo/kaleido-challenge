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
//	check-in mined → ticket checked_in; else the ticket stays active
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
			`UPDATE tickets SET status = 'checked_in' WHERE checkin_tx_id = $1 AND status = 'active'`, op.ID); err != nil {
			return err
		}
	}
	return nil
}

func keccak(b []byte) common.Hash { return crypto.Keccak256Hash(b) }
