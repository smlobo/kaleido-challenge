package chain

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// First Hardhat/Anvil dev account; public test vector, never a real key.
const hardhatKey0 = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func TestSignerAddress(t *testing.T) {
	s, err := NewSignerFromHex(hardhatKey0, ChainID1337)
	if err != nil {
		t.Fatal(err)
	}
	want := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	if s.Address() != want {
		t.Fatalf("address = %s, want %s", s.Address(), want)
	}
}

func TestSignProducesZeroFeeEIP1559ForChain(t *testing.T) {
	s, err := GenerateSigner(ChainID1337)
	if err != nil {
		t.Fatal(err)
	}
	to := common.HexToAddress("0x1")
	tx, err := s.Sign(TxRequest{Nonce: 7, To: &to, Data: []byte{0xde, 0xad}, Gas: 50_000})
	if err != nil {
		t.Fatal(err)
	}
	if tx.Type() != types.DynamicFeeTxType {
		t.Errorf("type = %d, want %d", tx.Type(), types.DynamicFeeTxType)
	}
	if tx.ChainId().Cmp(ChainID1337) != 0 {
		t.Errorf("chainId = %s", tx.ChainId())
	}
	if tx.GasFeeCap().Sign() != 0 || tx.GasTipCap().Sign() != 0 {
		t.Errorf("fee caps = %s/%s, want 0/0", tx.GasFeeCap(), tx.GasTipCap())
	}
	if tx.Nonce() != 7 || tx.Gas() != 50_000 {
		t.Errorf("nonce/gas = %d/%d", tx.Nonce(), tx.Gas())
	}
	from, err := types.Sender(types.LatestSignerForChainID(ChainID1337), tx)
	if err != nil {
		t.Fatal(err)
	}
	if from != s.Address() {
		t.Errorf("recovered sender %s, want %s", from, s.Address())
	}

	// Round-trips through the raw encoding sent to eth_sendRawTransaction.
	raw, err := tx.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var decoded types.Transaction
	if err := decoded.UnmarshalBinary(raw); err != nil {
		t.Fatal(err)
	}
	if decoded.Hash() != tx.Hash() {
		t.Errorf("hash changed after round trip")
	}
}

func TestNewSignerFromHexRejectsGarbage(t *testing.T) {
	if _, err := NewSignerFromHex("0xnothex", ChainID1337); err == nil {
		t.Fatal("expected error")
	}
}
