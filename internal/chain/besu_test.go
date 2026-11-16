package chain_test

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/contracts/demonft"
)

// dialBesu connects to the dev chain from docker/docker-compose.yaml, or skips
// the test if it is not running. Set BESU_RPC_URL to point elsewhere.
func dialBesu(t *testing.T) *chain.Client {
	t.Helper()
	url := os.Getenv("BESU_RPC_URL")
	if url == "" {
		url = "http://localhost:8545"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	c, err := chain.Dial(ctx, url, chain.ChainID1337)
	if err != nil {
		t.Skipf("Besu not reachable at %s (%v); run: docker compose -f docker/docker-compose.yaml up -d --wait", url, err)
	}
	t.Cleanup(c.Close)
	return c
}

func TestDeployMintAndReadBack(t *testing.T) {
	c := dialBesu(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// A brand-new key with zero balance: gas is free on this chain.
	owner, err := chain.GenerateSigner(c.ChainID())
	if err != nil {
		t.Fatal(err)
	}
	nft := demonft.NewDemoNFT()
	addr, rcpt, err := chain.Deploy(ctx, c, owner, common.FromHex(demonft.DemoNFTMetaData.Bin),
		nft.PackConstructor("Test", "TST", owner.Address()))
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	if addr == (common.Address{}) {
		t.Fatalf("no contract address in receipt %s", rcpt.TxHash)
	}

	holder := common.HexToAddress("0x000000000000000000000000000000000000bEEF")
	if _, err := chain.SendAndWait(ctx, c, owner, &addr, nft.PackSafeMint(holder, big.NewInt(42))); err != nil {
		t.Fatalf("safeMint: %v", err)
	}

	got, err := bind.Call(nft.Instance(c.Eth(), addr), &bind.CallOpts{Context: ctx},
		nft.PackOwnerOf(big.NewInt(42)), nft.UnpackOwnerOf)
	if err != nil {
		t.Fatalf("ownerOf: %v", err)
	}
	if got != holder {
		t.Fatalf("ownerOf(42) = %s, want %s", got, holder)
	}
}

func TestNonOwnerMintIsRejected(t *testing.T) {
	c := dialBesu(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	owner, _ := chain.GenerateSigner(c.ChainID())
	stranger, _ := chain.GenerateSigner(c.ChainID())
	nft := demonft.NewDemoNFT()
	addr, _, err := chain.Deploy(ctx, c, owner, common.FromHex(demonft.DemoNFTMetaData.Bin),
		nft.PackConstructor("Test", "TST", owner.Address()))
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}

	// Gas estimation runs the call, so Ownable's revert surfaces before anything is sent.
	_, err = chain.SendAndWait(ctx, c, stranger, &addr, nft.PackSafeMint(stranger.Address(), big.NewInt(1)))
	if err == nil {
		t.Fatal("expected stranger's safeMint to fail")
	}
	t.Logf("rejected as expected: %v", err)
}
