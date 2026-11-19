package ticketing

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
)

var ticketABI = sync.OnceValue(func() *abi.ABI {
	a, err := ticketnft.TicketNFTMetaData.ParseABI()
	if err != nil {
		panic(err) // the committed ABI is fixed; a parse failure is a build bug
	}
	return a
})

// DecodeRevert turns revert data into something readable, such as
// "PriceAboveCap(120000000, 110000000)". TicketNFT's ABI includes its own
// errors and OpenZeppelin's; Error(string) and Panic(uint256) are decoded
// too. Anything else comes back as hex.
func DecodeRevert(data []byte) string {
	if len(data) == 0 {
		return "a revert with no reason"
	}
	if len(data) >= 4 {
		for _, e := range ticketABI().Errors {
			if !bytes.Equal(e.ID[:4], data[:4]) {
				continue
			}
			vals, err := e.Inputs.Unpack(data[4:])
			if err != nil {
				break
			}
			args := make([]string, len(vals))
			for i, v := range vals {
				if a, ok := v.(common.Address); ok {
					args[i] = strings.ToLower(a.Hex())
				} else {
					args[i] = fmt.Sprint(v)
				}
			}
			return e.Name + "(" + strings.Join(args, ", ") + ")"
		}
	}
	if reason, err := abi.UnpackRevert(data); err == nil {
		return fmt.Sprintf("Error(%q)", reason)
	}
	return hexutil.Encode(data)
}

// DecodeRevertHex is DecodeRevert for the hex string Besu puts in a receipt.
func DecodeRevertHex(s string) string {
	if s == "" {
		return ""
	}
	b, err := hexutil.Decode(s)
	if err != nil {
		return s
	}
	return DecodeRevert(b)
}
