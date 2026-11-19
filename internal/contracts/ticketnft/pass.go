package ticketnft

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// checkInTypes describes the message a holder signs, in EIP-712 terms. It
// must match TicketNFT.sol: the contract's domain is EIP712("TicketNFT", "1")
// (OpenZeppelin adds the chain ID and the contract's address), and its
// CHECKIN_TYPEHASH is keccak256("CheckIn(uint256 tokenId,uint256 deadline)").
var checkInTypes = apitypes.Types{
	"EIP712Domain": {
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	"CheckIn": {
		{Name: "tokenId", Type: "uint256"},
		{Name: "deadline", Type: "uint256"},
	},
}

// CheckInDigest is the hash a holder signs to present ticket tokenID at the
// gate of the event whose contract is at contract, valid until deadline
// (unix seconds). The domain ties it to that chain and contract, so a pass
// for ticket 7 at one event is useless at another.
//
// It uses go-ethereum's EIP-712 encoder, the one behind eth_signTypedData,
// rather than hashing by hand.
func CheckInDigest(chainID *big.Int, contract common.Address, tokenID, deadline *big.Int) (common.Hash, error) {
	digest, _, err := apitypes.TypedDataAndHash(apitypes.TypedData{
		Types:       checkInTypes,
		PrimaryType: "CheckIn",
		Domain: apitypes.TypedDataDomain{
			Name:              "TicketNFT",
			Version:           "1",
			ChainId:           (*math.HexOrDecimal256)(chainID),
			VerifyingContract: contract.Hex(),
		},
		Message: apitypes.TypedDataMessage{"tokenId": tokenID, "deadline": deadline},
	})
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(digest), nil
}
