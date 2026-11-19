package ticketnft

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// A fixed example, so an accidental change to the types or the domain shows
// up here rather than as every check-in reverting. The expected digest does
// not come from this code: Foundry's cast computes it by following the
// EIP-712 steps directly.
//
//	domainType=$(cast keccak "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)")
//	domain=$(cast keccak $(cast abi-encode "f(bytes32,bytes32,bytes32,uint256,address)" \
//	    $domainType $(cast keccak "TicketNFT") $(cast keccak "1") 1337 0x2afded4d3f878e2f506df1606482f73136cfa843))
//	struct=$(cast keccak $(cast abi-encode "f(bytes32,uint256,uint256)" \
//	    $(cast keccak "CheckIn(uint256 tokenId,uint256 deadline)") 7 1763571234))
//	cast keccak 0x1901${domain#0x}${struct#0x}
//
// The Foundry tests and the Besu tests check that the contract itself
// accepts passes signed over this digest.
func TestCheckInDigest(t *testing.T) {
	type inputs struct {
		chain, token, deadline *big.Int
		contract               common.Address
	}
	base := inputs{big.NewInt(1337), big.NewInt(7), big.NewInt(1763571234),
		common.HexToAddress("0x2afded4d3f878e2f506df1606482f73136cfa843")}
	got, err := CheckInDigest(base.chain, base.contract, base.token, base.deadline)
	if err != nil {
		t.Fatal(err)
	}
	if want := "0xc205da2633a4fcfa55df7e8ad6f26d82a626b374bc6f1b459e29c2883d704c6b"; got.Hex() != want {
		t.Fatalf("CheckInDigest = %s, want %s", got, want)
	}

	// Every input changes the digest: the base pass, with one thing changed.
	for name, in := range map[string]inputs{
		"chain":    {big.NewInt(1338), base.token, base.deadline, base.contract},
		"token":    {base.chain, big.NewInt(8), base.deadline, base.contract},
		"deadline": {base.chain, base.token, big.NewInt(1763571235), base.contract},
		"contract": {base.chain, base.token, base.deadline, common.HexToAddress("0x01")},
	} {
		d, err := CheckInDigest(in.chain, in.contract, in.token, in.deadline)
		if err != nil || d == got {
			t.Errorf("changing the %s left the digest the same (%v)", name, err)
		}
	}
}
