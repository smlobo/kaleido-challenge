package chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client is a thin JSON-RPC client for the Besu node.
type Client struct {
	rpc     *rpc.Client
	eth     *ethclient.Client
	chainID *big.Int
}

// Dial connects to the node and checks that it serves the expected chain.
func Dial(ctx context.Context, url string, wantChainID *big.Int) (*Client, error) {
	rc, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", url, err)
	}
	ec := ethclient.NewClient(rc)
	id, err := ec.ChainID(ctx)
	if err != nil {
		rc.Close()
		return nil, fmt.Errorf("eth_chainId: %w", err)
	}
	if wantChainID != nil && id.Cmp(wantChainID) != 0 {
		rc.Close()
		return nil, fmt.Errorf("chain id mismatch: node reports %s, want %s", id, wantChainID)
	}
	return &Client{rpc: rc, eth: ec, chainID: id}, nil
}

// Close releases the connection.
func (c *Client) Close() { c.rpc.Close() }

// ChainID is the chain ID the node reported at dial time.
func (c *Client) ChainID() *big.Int { return new(big.Int).Set(c.chainID) }

// Eth exposes the go-ethereum client, e.g. for abigen read-only bindings.
func (c *Client) Eth() *ethclient.Client { return c.eth }

// PendingNonce returns the next nonce the node expects from addr, counting
// transactions already in its pool.
func (c *Client) PendingNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return c.eth.PendingNonceAt(ctx, addr)
}

// EstimateGas asks the node how much gas a call or deployment needs, and adds
// headroom. Gas is free on this chain, so the headroom costs nothing.
func (c *Client) EstimateGas(ctx context.Context, from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error) {
	gas, err := c.eth.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	})
	if err != nil {
		return 0, fmt.Errorf("eth_estimateGas: %w", err)
	}
	return gas + gas/5, nil
}

// SendRaw submits an already-signed transaction with eth_sendRawTransaction and
// returns the hash the node reports.
func (c *Client) SendRaw(ctx context.Context, tx *types.Transaction) (common.Hash, error) {
	raw, err := tx.MarshalBinary()
	if err != nil {
		return common.Hash{}, fmt.Errorf("encode tx: %w", err)
	}
	var hash common.Hash
	if err := c.rpc.CallContext(ctx, &hash, "eth_sendRawTransaction", hexutil.Encode(raw)); err != nil {
		return common.Hash{}, fmt.Errorf("eth_sendRawTransaction: %w", err)
	}
	if hash != tx.Hash() {
		return hash, fmt.Errorf("node returned hash %s, expected %s", hash, tx.Hash())
	}
	return hash, nil
}

// ErrReverted is returned by WaitMined when the transaction was mined but failed.
var ErrReverted = errors.New("transaction reverted")

// WaitMined polls for the receipt until it appears or ctx ends. QBFT has
// immediate finality, so a single receipt is final.
func (c *Client) WaitMined(ctx context.Context, hash common.Hash, poll time.Duration) (*types.Receipt, error) {
	t := time.NewTicker(poll)
	defer t.Stop()
	for {
		rcpt, err := c.eth.TransactionReceipt(ctx, hash)
		switch {
		case err == nil:
			if rcpt.Status != types.ReceiptStatusSuccessful {
				return rcpt, fmt.Errorf("%w: %s", ErrReverted, hash)
			}
			return rcpt, nil
		case !errors.Is(err, ethereum.NotFound):
			return nil, fmt.Errorf("eth_getTransactionReceipt: %w", err)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.C:
		}
	}
}
