package chain

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// receiptPoll is half the 2 s block time.
const receiptPoll = time.Second

// SendAndWait estimates gas, signs, sends with eth_sendRawTransaction and waits
// for the receipt. It reads the nonce from the node, which only suits a single
// sender (the milestone 1 CLI and tests). The service will use the Postgres
// nonce reservation and the outbox instead.
func SendAndWait(ctx context.Context, c *Client, s *Signer, to *common.Address, data []byte) (*types.Receipt, error) {
	nonce, err := c.PendingNonce(ctx, s.Address())
	if err != nil {
		return nil, err
	}
	gas, err := c.EstimateGas(ctx, s.Address(), to, nil, data)
	if err != nil {
		return nil, err
	}
	tx, err := s.Sign(TxRequest{Nonce: nonce, To: to, Data: data, Gas: gas})
	if err != nil {
		return nil, err
	}
	hash, err := c.SendRaw(ctx, tx)
	if err != nil {
		return nil, err
	}
	return c.WaitMined(ctx, hash, receiptPoll)
}

// Deploy sends creation bytecode followed by the packed constructor args
// (from an abigen v2 binding's PackConstructor) and returns the new address.
func Deploy(ctx context.Context, c *Client, s *Signer, bytecode, packedArgs []byte) (common.Address, *types.Receipt, error) {
	data := append(append([]byte{}, bytecode...), packedArgs...)
	rcpt, err := SendAndWait(ctx, c, s, nil, data)
	if err != nil {
		return common.Address{}, rcpt, err
	}
	return rcpt.ContractAddress, rcpt, nil
}

// ChainID1337 is the local Besu development chain.
var ChainID1337 = big.NewInt(1337)
