package chain_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
)

var (
	faceValue = big.NewInt(100_000_000) // 100.00 with 6 decimals
	resaleCap = big.NewInt(110_000_000) // 110% of face value
)

// ticketEvent is one deployed TicketNFT plus the keys that act on it.
type ticketEvent struct {
	c               *chain.Client
	nft             *ticketnft.TicketNFT
	addr            common.Address
	deployer        *chain.Signer
	boxOffice, gate *chain.Signer
	alice, bob      *chain.Signer
}

func deployTicketEvent(ctx context.Context, t *testing.T, c *chain.Client) *ticketEvent {
	t.Helper()
	keys := make([]*chain.Signer, 5)
	for i := range keys {
		keys[i], _ = chain.GenerateSigner(c.ChainID())
	}
	e := &ticketEvent{c: c, nft: ticketnft.NewTicketNFT(),
		deployer: keys[0], boxOffice: keys[1], gate: keys[2], alice: keys[3], bob: keys[4]}

	start := uint64(time.Now().Add(time.Hour).Unix())
	addr, _, err := chain.Deploy(ctx, c, e.deployer, common.FromHex(ticketnft.TicketNFTMetaData.Bin),
		e.nft.PackConstructor("Go Test Fest", "GTF", 100, 11_000, start, "http://localhost:8080/tickets/",
			e.boxOffice.Address(), e.gate.Address()))
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	e.addr = addr
	return e
}

func (e *ticketEvent) send(ctx context.Context, t *testing.T, from *chain.Signer, data []byte) {
	t.Helper()
	if _, err := chain.SendAndWait(ctx, e.c, from, &e.addr, data); err != nil {
		t.Fatal(err)
	}
}

// pass is the holder's signed check-in pass for ticket id, valid for five minutes.
func (e *ticketEvent) pass(t *testing.T, holder *chain.Signer, id *big.Int) (*big.Int, []byte) {
	t.Helper()
	deadline := big.NewInt(time.Now().Add(5 * time.Minute).Unix())
	digest, err := ticketnft.CheckInDigest(e.c.ChainID(), e.addr, id, deadline)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := holder.SignHash(digest)
	if err != nil {
		t.Fatal(err)
	}
	return deadline, sig
}

func (e *ticketEvent) info(ctx context.Context, t *testing.T, id *big.Int) ticketnft.TicketInfoOutput {
	t.Helper()
	out, err := bind.Call(e.nft.Instance(e.c.Eth(), e.addr), &bind.CallOpts{Context: ctx},
		e.nft.PackTicketInfo(id), e.nft.UnpackTicketInfo)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func (e *ticketEvent) owner(ctx context.Context, id *big.Int) (common.Address, error) {
	return bind.Call(e.nft.Instance(e.c.Eth(), e.addr), &bind.CallOpts{Context: ctx},
		e.nft.PackOwnerOf(id), e.nft.UnpackOwnerOf)
}

// The thin thread on a real chain: mint to A, A lists at the cap, B buys,
// the gate checks B in.
func TestTicketNFTLifecycle(t *testing.T) {
	c := dialBesu(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	e := deployTicketEvent(ctx, t, c)
	id := big.NewInt(1)
	record := crypto.Keccak256Hash([]byte(`{"event":"go-test-fest","ticket":1}`))

	e.send(ctx, t, e.boxOffice, e.nft.PackMintTicket(e.alice.Address(), id, faceValue, record))
	if got, _ := e.owner(ctx, id); got != e.alice.Address() {
		t.Fatalf("after mint, owner = %s, want alice", got)
	}

	e.send(ctx, t, e.alice, e.nft.PackListForResale(id, resaleCap))
	if got := e.info(ctx, t, id).ListingPrice; got.Cmp(resaleCap) != 0 {
		t.Fatalf("listing price = %s, want %s", got, resaleCap)
	}

	rcpt, err := chain.SendAndWait(ctx, c, e.bob, &e.addr, e.nft.PackBuyResale(id, resaleCap))
	if err != nil {
		t.Fatal(err)
	}
	var resold *ticketnft.TicketNFTTicketResold
	for _, l := range rcpt.Logs {
		if ev, err := e.nft.UnpackTicketResoldEvent(l); err == nil {
			resold = ev
		}
	}
	if resold == nil || resold.From != e.alice.Address() || resold.To != e.bob.Address() || resold.Price.Cmp(resaleCap) != 0 {
		t.Fatalf("TicketResold event = %+v", resold)
	}
	if got, _ := e.owner(ctx, id); got != e.bob.Address() {
		t.Fatalf("after resale, owner = %s, want bob", got)
	}

	// The gate checks bob in with bob's pass, signed in Go: the contract
	// accepts the Go side's EIP-712 digest and signature.
	deadline, sig := e.pass(t, e.bob, id)
	e.send(ctx, t, e.gate, e.nft.PackCheckIn(id, deadline, sig))
	if _, err := e.owner(ctx, id); err == nil {
		t.Fatal("ticket still has an owner after check-in")
	}
	info := e.info(ctx, t, id)
	if info.FaceValue.Cmp(faceValue) != 0 || common.Hash(info.RecordHash) != record || info.ListingPrice.Sign() != 0 {
		t.Fatalf("after check-in, ticketInfo = %+v", info)
	}
}

// Rejected calls fail at gas estimation, before anything is signed, and the
// contract's custom error reaches Go intact. The API will use this to tell a
// buyer "the price changed" rather than "execution reverted".
func TestTicketNFTRevertsDecodeInGo(t *testing.T) {
	c := dialBesu(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	e := deployTicketEvent(ctx, t, c)
	id := big.NewInt(7)
	e.send(ctx, t, e.boxOffice, e.nft.PackMintTicket(e.alice.Address(), id, faceValue, common.Hash{1}))
	e.send(ctx, t, e.alice, e.nft.PackListForResale(id, faceValue))

	decode := func(from *chain.Signer, data []byte) any {
		t.Helper()
		_, err := c.EstimateGas(ctx, from.Address(), &e.addr, nil, data)
		if err == nil {
			t.Fatal("call succeeded, want a revert")
		}
		raw, ok := ethclient.RevertErrorData(err)
		if !ok {
			t.Fatalf("no revert data in %v", err)
		}
		decoded, err := e.nft.UnpackError(raw)
		if err != nil {
			t.Fatalf("decode revert: %v", err)
		}
		return decoded
	}

	overCap := new(big.Int).Add(resaleCap, big.NewInt(1))
	if got, ok := decode(e.alice, e.nft.PackListForResale(id, overCap)).(*ticketnft.TicketNFTPriceAboveCap); !ok ||
		got.Price.Cmp(overCap) != 0 || got.Cap.Cmp(resaleCap) != 0 {
		t.Errorf("list above cap: got %#v", got)
	}

	tooLow := big.NewInt(50_000_000)
	if got, ok := decode(e.bob, e.nft.PackBuyResale(id, tooLow)).(*ticketnft.TicketNFTPriceAboveMaxPrice); !ok ||
		got.Price.Cmp(faceValue) != 0 || got.MaxPrice.Cmp(tooLow) != 0 {
		t.Errorf("list above buyer’s max price: got %#v", got)
	}

	if _, ok := decode(e.alice, e.nft.PackTransferFrom(e.alice.Address(), e.bob.Address(), id)).(*ticketnft.TicketNFTTransfersDisabled); !ok {
		t.Error("transferFrom: want TransfersDisabled")
	}

	if got, ok := decode(e.bob, e.nft.PackMintTicket(e.bob.Address(), big.NewInt(8), faceValue, common.Hash{1})).(*ticketnft.TicketNFTAccessControlUnauthorizedAccount); !ok ||
		got.Account != e.bob.Address() {
		t.Errorf("mint by non-minter: got %#v", got)
	}

	// Check-in needs both parties: a pass signed by someone other than the
	// holder (alice) is refused, and so is the holder without the gate.
	deadline, bobsPass := e.pass(t, e.bob, id)
	if got, ok := decode(e.gate, e.nft.PackCheckIn(id, deadline, bobsPass)).(*ticketnft.TicketNFTNotSignedByHolder); !ok ||
		got.TokenId.Cmp(id) != 0 {
		t.Errorf("check-in with someone else's pass: got %#v", got)
	}
	deadline, alicesPass := e.pass(t, e.alice, id)
	if _, ok := decode(e.alice, e.nft.PackCheckIn(id, deadline, alicesPass)).(*ticketnft.TicketNFTAccessControlUnauthorizedAccount); !ok {
		t.Error("holder checks herself in: want AccessControlUnauthorizedAccount")
	}
	if _, err := c.EstimateGas(ctx, e.gate.Address(), &e.addr, nil, e.nft.PackCheckIn(id, deadline, alicesPass)); err != nil {
		t.Errorf("gate with the holder's pass: %v", err)
	}

	// The box office, not the deployer, is admin: the deployer can't grant itself minting
	minterRole := crypto.Keccak256Hash([]byte("MINTER_ROLE"))
	if got, ok := decode(e.deployer, e.nft.PackGrantRole(minterRole, e.deployer.Address())).(*ticketnft.TicketNFTAccessControlUnauthorizedAccount); !ok ||
		got.Account != e.deployer.Address() || got.NeededRole != [32]byte{} {
		t.Errorf("deployer grants itself minter: got %#v", got)
	}
	if _, err := c.EstimateGas(ctx, e.boxOffice.Address(), &e.addr, nil, e.nft.PackGrantRole(minterRole, e.bob.Address())); err != nil {
		t.Errorf("box office grants minter: %v", err)
	}
}
