// Command m1deploy is the milestone 1 smoke test. It deploys a plain ERC-721
// to the Besu chain using a transaction the service signs itself and sends with
// eth_sendRawTransaction, mints one token and reads its owner back.
//
//	RPC_URL      JSON-RPC endpoint (default http://localhost:8545)
//	DEPLOYER_KEY hex private key; if unset a throwaway key is generated
//	             (gas is free on this chain, so it needs no balance)
package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/contracts/demonft"
)

func main() {
	if err := run(); err != nil {
		slog.Error("m1deploy failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rpcURL := envOr("RPC_URL", "http://localhost:8545")
	client, err := chain.Dial(ctx, rpcURL, chain.ChainID1337)
	if err != nil {
		return err
	}
	defer client.Close()

	var signer *chain.Signer
	if k := os.Getenv("DEPLOYER_KEY"); k != "" {
		signer, err = chain.NewSignerFromHex(k, client.ChainID())
	} else {
		signer, err = chain.GenerateSigner(client.ChainID())
		slog.Info("DEPLOYER_KEY not set; using a throwaway key")
	}
	if err != nil {
		return err
	}
	slog.Info("connected", "rpc", rpcURL, "chainId", client.ChainID(), "deployer", signer.Address())

	nft := demonft.NewDemoNFT()
	addr, rcpt, err := chain.Deploy(ctx, client, signer, common.FromHex(demonft.DemoNFTMetaData.Bin),
		nft.PackConstructor("Demo Ticket", "DEMO", signer.Address()))
	if err != nil {
		return fmt.Errorf("deploy: %w", err)
	}
	slog.Info("deployed DemoNFT", "address", addr, "tx", rcpt.TxHash, "block", rcpt.BlockNumber)

	holder, err := chain.GenerateSigner(client.ChainID())
	if err != nil {
		return err
	}
	tokenID := big.NewInt(1)
	rcpt, err = chain.SendAndWait(ctx, client, signer, &addr, nft.PackSafeMint(holder.Address(), tokenID))
	if err != nil {
		return fmt.Errorf("safeMint: %w", err)
	}
	slog.Info("minted", "tokenId", tokenID, "to", holder.Address(), "tx", rcpt.TxHash, "block", rcpt.BlockNumber)

	// Read back through the abigen binding (eth_call, no signing).
	owner, err := bind.Call(nft.Instance(client.Eth(), addr),
		&bind.CallOpts{Context: ctx}, nft.PackOwnerOf(tokenID), nft.UnpackOwnerOf)
	if err != nil {
		return fmt.Errorf("ownerOf: %w", err)
	}
	if owner != holder.Address() {
		return fmt.Errorf("ownerOf(1) = %s, want %s", owner, holder.Address())
	}
	slog.Info("verified on-chain", "ownerOf(1)", owner)
	return nil
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
