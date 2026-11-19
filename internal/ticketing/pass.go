package ticketing

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"

	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
)

// PassLifetime is how long a check-in pass stays valid. Short, so a
// screenshot passed around soon stops working; the burn at check-in makes
// every pass single-use anyway.
const PassLifetime = 5 * time.Minute

// Pass is a holder's signed permission to check one ticket in: the holder
// shows it at the door (a QR code) and the gate sends it with the check-in.
// The contract checks the signature against the ticket's current holder.
type Pass struct {
	TicketID, TokenID string
	Contract          string // the event's contract address
	ChainID           *big.Int
	Deadline          time.Time // valid until, whole seconds
	Signature         string    // 0x-prefixed 65 bytes: r ‖ s ‖ v
}

// PassInput is the part of a pass the gate sends back.
type PassInput struct {
	TicketID  string // optional; when given, must be the ticket being checked in
	Deadline  int64  // unix seconds
	Signature string
}

var errPassExpired = fail(Rejected, "pass_expired", "Pass expired",
	"The check-in pass has expired. Ask the holder for a new one.")

// CreatePass signs a check-in pass for one of the caller's tickets with
// their custodial wallet key. It writes nothing and sends nothing: the pass
// is only a signature, which the gate later carries into its transaction.
func (s *Service) CreatePass(ctx context.Context, p Principal, ticketID string) (Pass, error) {
	if err := need(p, User); err != nil {
		return Pass{}, err
	}
	if uuid.Validate(ticketID) != nil {
		return Pass{}, errNotFound("ticket")
	}
	var (
		status, tokenID, walletID, wallet string
		owner, contract                   *string
	)
	err := s.pool.QueryRow(ctx,
		`SELECT t.status, t.token_id::text, t.owner_address, e.contract_address, u.wallet_key_id, k.address
		 FROM tickets t
		 JOIN events e ON e.id = t.event_id
		 JOIN users u ON u.id = $2
		 JOIN signing_keys k ON k.id = u.wallet_key_id
		 WHERE t.id = $1`, ticketID, p.ID).Scan(&status, &tokenID, &owner, &contract, &walletID, &wallet)
	if isNoRows(err) {
		return Pass{}, errNotFound("ticket")
	}
	if err != nil {
		return Pass{}, err
	}
	switch {
	case status == "pending_mint":
		return Pass{}, errNotMintedYet
	case status == "failed":
		return Pass{}, errTicketFailed
	case status == "checked_in":
		return Pass{}, errAlreadyCheckedIn
	case owner == nil || !strings.EqualFold(*owner, wallet):
		return Pass{}, fail(Forbidden, "not_your_ticket", "Not your ticket", "Only the ticket's holder can make a check-in pass for it.")
	}

	deadline := s.now().Add(PassLifetime).Truncate(time.Second)
	token, _ := new(big.Int).SetString(tokenID, 10)
	digest, err := ticketnft.CheckInDigest(s.chainID, common.HexToAddress(*contract), token, big.NewInt(deadline.Unix()))
	if err != nil {
		return Pass{}, err
	}
	signer, err := s.keys.Signer(ctx, s.pool, walletID)
	if err != nil {
		return Pass{}, err
	}
	sig, err := signer.SignHash(digest)
	if err != nil {
		return Pass{}, err
	}
	return Pass{
		TicketID: ticketID, TokenID: tokenID, Contract: *contract, ChainID: s.chainID,
		Deadline: deadline, Signature: hexutil.Encode(sig),
	}, nil
}

// parsePass checks a pass's form before anything is dry-run or signed. Whether
// the holder really signed it is the contract's call.
func (s *Service) parsePass(ticketID string, in PassInput) (deadline *big.Int, sig []byte, err error) {
	if in.TicketID != "" && !strings.EqualFold(in.TicketID, ticketID) {
		return nil, nil, fail(Invalid, "pass_for_another_ticket", "Pass for another ticket",
			"This pass was made for ticket "+in.TicketID+".")
	}
	sig, decErr := hexutil.Decode(in.Signature)
	if decErr != nil || len(sig) != 65 {
		return nil, nil, fail(Invalid, "invalid_pass", "Invalid pass",
			"signature must be 65 bytes as 0x-prefixed hex, as POST /tickets/{id}/pass returns it.")
	}
	if in.Deadline <= 0 {
		return nil, nil, fail(Invalid, "invalid_pass", "Invalid pass", "deadline must be a unix time in seconds.")
	}
	if s.now().Unix() > in.Deadline {
		return nil, nil, errPassExpired
	}
	return big.NewInt(in.Deadline), sig, nil
}
