// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ticketnft

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// TicketNFTMetaData contains all meta data concerning the TicketNFT contract.
var TicketNFTMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"maxSupply_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxResaleBps_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"eventStart_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"baseURI_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"boxOffice\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gate\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CHECKIN_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"GATE_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MINTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buyResale\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelListing\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkIn\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"holderSig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eventStart\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"listForResale\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxResaleBps\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mintTicket\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mintedCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resaleCap\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ticketInfo\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"faceValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"listingPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CheckedIn\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"holder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ListingCancelled\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TicketListed\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TicketMinted\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recordHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TicketResold\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BuyerIsSeller\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EmptyRecordHash\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EventStarted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotListed\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotSignedByHolder\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotTicketHolder\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"PassExpired\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"PriceAboveCap\",\"inputs\":[{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"PriceAboveMaxPrice\",\"inputs\":[{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SoldOut\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"TransfersDisabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroPrice\",\"inputs\":[]}]",
	ID:  "TicketNFT",
	Bin: "0x6101c0604052348015610010575f5ffd5b5060405161276438038061276483398101604081905261002f91610450565b60405180604001604052806009815260200168151a58dad95d13919560ba1b815250604051806040016040528060018152602001603160f81b8152508989815f908161007b91906105bf565b50600161008882826105bf565b506100989150839050600761024c565b610120526100a781600861024c565b61014052815160208084019190912060e052815190820120610100524660a05261013360e05161010051604080517f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f60208201529081019290925260608201524660808201523060a08201525f9060c00160405160208183030381529060405280519060200120905090565b60805250503060c05263ffffffff86161580610158575042846001600160401b031611155b15610176576040516306b7c75960e31b815260040160405180910390fd5b6001600160a01b038216158061019357506001600160a01b038116155b156101b1576040516306b7c75960e31b815260040160405180910390fd5b63ffffffff808716610160528516610180526001600160401b0384166101a052600b6101dd84826105bf565b506101e85f8361027e565b506102137f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a68361027e565b5061023e7f79aaa02ca9ee405c105f18ab3335a0dcabff839ebd075c93f39bce4cef295be48261027e565b5050505050505050506106d5565b5f6020835110156102675761026083610329565b9050610278565b8161027284826105bf565b5060ff90505b92915050565b5f8281526006602090815260408083206001600160a01b038516845290915281205460ff16610322575f8381526006602090815260408083206001600160a01b03861684529091529020805460ff191660011790556102da3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610278565b505f610278565b5f5f829050601f8151111561035c578260405163305a27a960e01b8152600401610353919061067d565b60405180910390fd5b8051610367826106b2565b179392505050565b634e487b7160e01b5f52604160045260245ffd5b5f82601f830112610392575f5ffd5b81516001600160401b038111156103ab576103ab61036f565b604051601f8201601f19908116603f011681016001600160401b03811182821017156103d9576103d961036f565b6040528181528382016020018510156103f0575f5ffd5b8160208501602083015e5f918101602001919091529392505050565b805163ffffffff8116811461041f575f5ffd5b919050565b80516001600160401b038116811461041f575f5ffd5b80516001600160a01b038116811461041f575f5ffd5b5f5f5f5f5f5f5f5f610100898b031215610468575f5ffd5b88516001600160401b0381111561047d575f5ffd5b6104898b828c01610383565b60208b015190995090506001600160401b038111156104a6575f5ffd5b6104b28b828c01610383565b9750506104c160408a0161040c565b95506104cf60608a0161040c565b94506104dd60808a01610424565b60a08a01519094506001600160401b038111156104f8575f5ffd5b6105048b828c01610383565b93505061051360c08a0161043a565b915061052160e08a0161043a565b90509295985092959890939650565b600181811c9082168061054457607f821691505b60208210810361056257634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156105ba57828211156105ba57805f5260205f20601f840160051c602085101561059357505f5b90810190601f840160051c035f5b818110156105b6575f838201556001016105a1565b5050505b505050565b81516001600160401b038111156105d8576105d861036f565b6105ec816105e68454610530565b84610568565b6020601f82116001811461061e575f83156106075750848201515b5f19600385901b1c1916600184901b178455610676565b5f84815260208120601f198516915b8281101561064d578785015182556020948501946001909201910161062d565b508482101561066a57868401515f19600387901b60f8161c191681555b505060018360011b0184555b5050505050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b80516020808301519190811015610562575f1960209190910360031b1b16919050565b60805160a05160c05160e05161010051610120516101405161016051610180516101a051611ffd6107675f395f81816102b7015281816106b0015281816108d40152610b1801525f81816103da015261068601525f81816104d8015261092601525f6111a901525f61117701525f6117cc01525f6117a401525f6116ff01525f61172901525f6117530152611ffd5ff3fe608060405234801561000f575f5ffd5b50600436106101fd575f3560e01c806391d1485411610114578063c87b56dd116100a9578063d5abeb0111610079578063d5abeb01146104d3578063df061368146104fa578063e985e9c514610521578063fa71308314610534578063fd649abf1461057b575f5ffd5b8063c87b56dd14610476578063cf721b1514610489578063d539139314610499578063d547741f146104c0575f5ffd5b8063a22cb465116100e4578063a22cb46514610420578063a89d2d381461042e578063b48c7f4914610441578063b88d4fde14610468575f5ffd5b806391d14854146103c2578063923aee8e146103d557806395d89b4114610411578063a217fddf14610419575f5ffd5b8063248a9ca31161019557806342842e0e1161016557806342842e0e1461035b57806359e1c2511461036e5780636352211e1461038157806370a082311461039457806384b0196e146103a7575f5ffd5b8063248a9ca3146103005780632f2ff15d14610322578063305a67a81461033557806336568abe14610348575f5ffd5b80630e0caabf116101d05780630e0caabf1461027e5780631f7d863c1461029f578063202887e7146102b257806323b872dd146102f2575f5ffd5b806301ffc9a71461020157806306fdde0314610229578063081812fc1461023e578063095ea7b314610269575b5f5ffd5b61021461020f366004611b14565b61058e565b60405190151581526020015b60405180910390f35b61023161059e565b6040516102209190611b69565b61025161024c366004611b7b565b61062d565b6040516001600160a01b039091168152602001610220565b61027c610277366004611bad565b610654565b005b61029161028c366004611b7b565b61066d565b604051908152602001610220565b61027c6102ad366004611bd5565b6106ae565b6102d97f000000000000000000000000000000000000000000000000000000000000000081565b60405167ffffffffffffffff9091168152602001610220565b61027c610277366004611bf5565b61029161030e366004611b7b565b5f9081526006602052604090206001015490565b61027c610330366004611c2f565b6107b3565b61027c610343366004611b7b565b6107dd565b61027c610356366004611c2f565b610856565b61027c610369366004611bf5565b61088e565b61027c61037c366004611c59565b6108a8565b61025161038f366004611b7b565b610a4c565b6102916103a2366004611c8f565b610a56565b6103af610a9b565b6040516102209796959493929190611ca8565b6102146103d0366004611c2f565b610add565b6103fc7f000000000000000000000000000000000000000000000000000000000000000081565b60405163ffffffff9091168152602001610220565b610231610b07565b6102915f81565b61027c610277366004611d3e565b61027c61043c366004611bd5565b610b16565b6102917fcc929a19526c6c9720325113f7f4aa9cec373f8599f0c95985df426ff92aeec081565b61027c610277366004611d8b565b610231610484366004611b7b565b610c65565b6009546103fc9063ffffffff1681565b6102917f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a681565b61027c6104ce366004611c2f565b610cca565b6103fc7f000000000000000000000000000000000000000000000000000000000000000081565b6102917f79aaa02ca9ee405c105f18ab3335a0dcabff839ebd075c93f39bce4cef295be481565b61021461052f366004611e68565b610cee565b610560610542366004611b7b565b5f908152600a60205260409020805460018201546002909201549092565b60408051938452602084019290925290820152606001610220565b61027c610589366004611e90565b610d1b565b5f61059882610ec6565b92915050565b60605f80546105ac90611f0c565b80601f01602080910402602001604051908101604052809291908181526020018280546105d890611f0c565b80156106235780601f106105fa57610100808354040283529160200191610623565b820191905f5260205f20905b81548152906001019060200180831161060657829003601f168201915b5050505050905090565b5f61063782610eea565b505f828152600460205260409020546001600160a01b0316610598565b604051638574adcf60e01b815260040160405180910390fd5b5f818152600a60205260408120546105989063ffffffff7f000000000000000000000000000000000000000000000000000000000000000016612710610f22565b7f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1642106106f85760405163292403c360e01b815260040160405180910390fd5b61070182610fd2565b805f0361072157604051634dfba02360e01b815260040160405180910390fd5b5f61072b8361066d565b90508082111561075d57604051634216244560e11b815260048101839052602481018290526044015b60405180910390fd5b5f838152600a60205260409081902060020183905551339084907f22968ec3be3687b180f3544163a034b88d57b48159d9a7e9a1b610fdb2d05f80906107a69086815260200190565b60405180910390a3505050565b5f828152600660205260409020600101546107cd81611009565b6107d78383611013565b50505050565b6107e681610fd2565b5f818152600a602052604081206002015490036108195760405163068bb67560e51b815260048101829052602401610754565b5f818152600a60205260408082206002018290555182917f411aee90354c51b1b04cd563fcab2617142a9d50da19232d888547c8a1b7fd8a91a250565b6001600160a01b038116331461087f5760405163334bd91960e11b815260040160405180910390fd5b61088982826110a4565b505050565b61088983838360405180602001604052805f815250610654565b7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a66108d281611009565b7f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff16421061091c5760405163292403c360e01b815260040160405180910390fd5b60095463ffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911610610967576040516352df9fe560e01b815260040160405180910390fd5b8161098557604051631ed0d70b60e11b815260040160405180910390fd5b6009805463ffffffff16905f61099a83611f44565b825463ffffffff9182166101009390930a9283029190920219909116179055506040805160608101825284815260208082018581525f838501818152898252600a90935293909320915182559151600182015590516002909101556109ff858561110f565b60408051848152602081018490526001600160a01b0387169186917f8811fc041e6a57d95124333f7756d53a5041b848d804206cba405dafd50661d4910160405180910390a35050505050565b5f61059882610eea565b5f6001600160a01b038216610a80576040516322718ad960e21b81525f6004820152602401610754565b506001600160a01b03165f9081526003602052604090205490565b5f6060805f5f5f6060610aac611170565b610ab46111a2565b604080515f80825260208201909252600f60f81b9b939a50919850469750309650945092509050565b5f9182526006602090815260408084206001600160a01b0393909316845291905290205460ff1690565b6060600180546105ac90611f0c565b7f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff164210610b605760405163292403c360e01b815260040160405180910390fd5b5f828152600a602052604081206002015490819003610b955760405163068bb67560e51b815260048101849052602401610754565b81811115610bc05760405163f319e5e360e01b81526004810182905260248101839052604401610754565b5f610bca84610a4c565b90506001600160a01b0381163303610bf557604051636183505960e01b815260040160405180910390fd5b5f848152600a6020526040812060020155610c118133866111cf565b336001600160a01b0316816001600160a01b0316857fd59a36d8310ffb140214be490baf3fe8961f859a17199ec1bf9015c8b734001185604051610c5791815260200190565b60405180910390a450505050565b6060610c7082610eea565b505f610c7a61127c565b90505f815111610c985760405180602001604052805f815250610cc3565b80610ca28461128b565b604051602001610cb3929190611f8b565b6040516020818303038152906040525b9392505050565b5f82815260066020526040902060010154610ce481611009565b6107d783836110a4565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b7f79aaa02ca9ee405c105f18ab3335a0dcabff839ebd075c93f39bce4cef295be4610d4581611009565b5f610d4f86610a4c565b905084421115610d7557604051630cb9443960e11b815260048101869052602401610754565b604080517fcc929a19526c6c9720325113f7f4aa9cec373f8599f0c95985df426ff92aeec06020820152908101879052606081018690525f90610dd0906080016040516020818303038152906040528051906020012061131b565b90505f5f610e138388888080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061134792505050565b5090925090505f816003811115610e2c57610e2c611f9f565b141580610e4b5750836001600160a01b0316826001600160a01b031614155b15610e6c57604051631cec443d60e01b8152600481018a9052602401610754565b5f898152600a6020526040812060020155610e8689611390565b6040516001600160a01b038516908a907f70790c0ab38e2aca9b5520985fe66db2a3700c3680b61c3853b5554f67a2c31c905f90a3505050505050505050565b5f6001600160e01b03198216637965db0b60e01b14806105985750610598826113cc565b5f818152600260205260408120546001600160a01b03168061059857604051637e27328960e01b815260048101849052602401610754565b5f5f5f610f2f868661141b565b91509150815f03610f5357838181610f4957610f49611fb3565b0492505050610cc3565b818411610f6a57610f6a6003851502601118611437565b5f848688095f868103871696879004966002600389028118808a02820302808a02820302808a02820302808a02820302808a02820302808a02909103029181900381900460010185841190960395909502919093039390930492909217029150509392505050565b33610fdc82610a4c565b6001600160a01b031614611006576040516316344afb60e31b815260048101829052602401610754565b50565b6110068133611448565b5f61101e8383610add565b61109d575f8381526006602090815260408083206001600160a01b03861684529091529020805460ff191660011790556110553390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610598565b505f610598565b5f6110af8383610add565b1561109d575f8381526006602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4506001610598565b6001600160a01b03821661113857604051633250574960e11b81525f6004820152602401610754565b5f61114483835f611481565b90506001600160a01b03811615610889576040516339e3563760e11b81525f6004820152602401610754565b606061119d7f00000000000000000000000000000000000000000000000000000000000000006007611573565b905090565b606061119d7f00000000000000000000000000000000000000000000000000000000000000006008611573565b6001600160a01b0382166111f857604051633250574960e11b81525f6004820152602401610754565b5f61120483835f611481565b90506001600160a01b03811661123057604051637e27328960e01b815260048101839052602401610754565b836001600160a01b0316816001600160a01b0316146107d7576040516364283d7b60e01b81526001600160a01b0380861660048301526024820184905282166044820152606401610754565b6060600b80546105ac90611f0c565b60605f6112978361161c565b60010190505f8167ffffffffffffffff8111156112b6576112b6611d77565b6040519080825280601f01601f1916602001820160405280156112e0576020820181803683370190505b5090508181016020015b5f19016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a85049450846112ea57509392505050565b5f6105986113276116f3565b8360405161190160f01b8152600281019290925260228201526042902090565b5f5f5f835160410361137e576020840151604085015160608601515f1a6113708882858561181c565b955095509550505050611389565b505081515f91506002905b9250925092565b5f61139c5f835f611481565b90506001600160a01b0381166113c857604051637e27328960e01b815260048101839052602401610754565b5050565b5f6001600160e01b031982166380ac58cd60e01b14806113fc57506001600160e01b03198216635b5e139f60e01b145b8061059857506301ffc9a760e01b6001600160e01b0319831614610598565b5f805f1983850993909202808410938190039390930393915050565b634e487b715f52806020526024601cfd5b6114528282610add565b6113c85760405163e2517d3f60e01b81526001600160a01b038216600482015260248101839052604401610754565b5f828152600260205260408120546001600160a01b03908116908316156114ad576114ad8184866118e4565b6001600160a01b038116156114e7576114c85f855f5f611948565b6001600160a01b0381165f90815260036020526040902080545f190190555b6001600160a01b03851615611515576001600160a01b0385165f908152600360205260409020805460010190555b5f8481526002602052604080822080546001600160a01b0319166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b606060ff831461158d5761158683611a4c565b9050610598565b81805461159990611f0c565b80601f01602080910402602001604051908101604052809291908181526020018280546115c590611f0c565b80156116105780601f106115e757610100808354040283529160200191611610565b820191905f5260205f20905b8154815290600101906020018083116115f357829003601f168201915b50505050509050610598565b5f8072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b831061165a5772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef81000000008310611686576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc1000083106116a457662386f26fc10000830492506010015b6305f5e10083106116bc576305f5e100830492506008015b61271083106116d057612710830492506004015b606483106116e2576064830492506002015b600a83106105985760010192915050565b5f306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614801561174b57507f000000000000000000000000000000000000000000000000000000000000000046145b1561177557507f000000000000000000000000000000000000000000000000000000000000000090565b61119d604080517f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f60208201527f0000000000000000000000000000000000000000000000000000000000000000918101919091527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a08201525f9060c00160405160208183030381529060405280519060200120905090565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a084111561185557505f915060039050826118da565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa1580156118a6573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b0381166118d157505f9250600191508290506118da565b92505f91508190505b9450945094915050565b6118ef838383611a89565b610889576001600160a01b03831661191d57604051637e27328960e01b815260048101829052602401610754565b60405163177e802f60e01b81526001600160a01b038316600482015260248101829052604401610754565b808061195c57506001600160a01b03821615155b15611a1d575f61196b84610eea565b90506001600160a01b038316158015906119975750826001600160a01b0316816001600160a01b031614155b80156119aa57506119a88184610cee565b155b156119d35760405163a9fbf51f60e01b81526001600160a01b0384166004820152602401610754565b8115611a1b5783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260046020526040902080546001600160a01b0319166001600160a01b0392909216919091179055565b60605f611a5883611aed565b6040805160208082528183019092529192505f91906020820181803683375050509182525060208101929092525090565b5f6001600160a01b03831615801590611ae55750826001600160a01b0316846001600160a01b03161480611ac25750611ac28484610cee565b80611ae557505f828152600460205260409020546001600160a01b038481169116145b949350505050565b5f60ff8216601f81111561059857604051632cd44ac360e21b815260040160405180910390fd5b5f60208284031215611b24575f5ffd5b81356001600160e01b031981168114610cc3575f5ffd5b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f610cc36020830184611b3b565b5f60208284031215611b8b575f5ffd5b5035919050565b80356001600160a01b0381168114611ba8575f5ffd5b919050565b5f5f60408385031215611bbe575f5ffd5b611bc783611b92565b946020939093013593505050565b5f5f60408385031215611be6575f5ffd5b50508035926020909101359150565b5f5f5f60608486031215611c07575f5ffd5b611c1084611b92565b9250611c1e60208501611b92565b929592945050506040919091013590565b5f5f60408385031215611c40575f5ffd5b82359150611c5060208401611b92565b90509250929050565b5f5f5f5f60808587031215611c6c575f5ffd5b611c7585611b92565b966020860135965060408601359560600135945092505050565b5f60208284031215611c9f575f5ffd5b610cc382611b92565b60ff60f81b8816815260e060208201525f611cc660e0830189611b3b565b8281036040840152611cd88189611b3b565b606084018890526001600160a01b038716608085015260a0840186905283810360c0850152845180825260208087019350909101905f5b81811015611d2d578351835260209384019390920191600101611d0f565b50909b9a5050505050505050505050565b5f5f60408385031215611d4f575f5ffd5b611d5883611b92565b915060208301358015158114611d6c575f5ffd5b809150509250929050565b634e487b7160e01b5f52604160045260245ffd5b5f5f5f5f60808587031215611d9e575f5ffd5b611da785611b92565b9350611db560208601611b92565b925060408501359150606085013567ffffffffffffffff811115611dd7575f5ffd5b8501601f81018713611de7575f5ffd5b803567ffffffffffffffff811115611e0157611e01611d77565b604051601f8201601f19908116603f0116810167ffffffffffffffff81118282101715611e3057611e30611d77565b604052818152828201602001891015611e47575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b5f5f60408385031215611e79575f5ffd5b611e8283611b92565b9150611c5060208401611b92565b5f5f5f5f60608587031215611ea3575f5ffd5b8435935060208501359250604085013567ffffffffffffffff811115611ec7575f5ffd5b8501601f81018713611ed7575f5ffd5b803567ffffffffffffffff811115611eed575f5ffd5b876020828401011115611efe575f5ffd5b949793965060200194505050565b600181811c90821680611f2057607f821691505b602082108103611f3e57634e487b7160e01b5f52602260045260245ffd5b50919050565b5f63ffffffff821663ffffffff8103611f6b57634e487b7160e01b5f52601160045260245ffd5b60010192915050565b5f81518060208401855e5f93019283525090919050565b5f611ae5611f998386611f74565b84611f74565b634e487b7160e01b5f52602160045260245ffd5b634e487b7160e01b5f52601260045260245ffdfea264697066735822122095f4a1a95db3156ed1b193766fa1b2bec6c684aa5a43b68bb854b338f3fbea4a64736f6c63430008260033",
}

// TicketNFT is an auto generated Go binding around an Ethereum contract.
type TicketNFT struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *TicketNFT) GetABI() abi.ABI {
	return c.abi
}

// NewTicketNFT creates a new instance of TicketNFT.
func NewTicketNFT() *TicketNFT {
	parsed, err := TicketNFTMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &TicketNFT{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *TicketNFT) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(string name_, string symbol_, uint32 maxSupply_, uint32 maxResaleBps_, uint64 eventStart_, string baseURI_, address boxOffice, address gate) returns()
func (ticketNFT *TicketNFT) PackConstructor(name_ string, symbol_ string, maxSupply_ uint32, maxResaleBps_ uint32, eventStart_ uint64, baseURI_ string, boxOffice common.Address, gate common.Address) []byte {
	enc, err := ticketNFT.abi.Pack("", name_, symbol_, maxSupply_, maxResaleBps_, eventStart_, baseURI_, boxOffice, gate)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCHECKINTYPEHASH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb48c7f49.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function CHECKIN_TYPEHASH() view returns(bytes32)
func (ticketNFT *TicketNFT) PackCHECKINTYPEHASH() []byte {
	enc, err := ticketNFT.abi.Pack("CHECKIN_TYPEHASH")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCHECKINTYPEHASH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb48c7f49.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function CHECKIN_TYPEHASH() view returns(bytes32)
func (ticketNFT *TicketNFT) TryPackCHECKINTYPEHASH() ([]byte, error) {
	return ticketNFT.abi.Pack("CHECKIN_TYPEHASH")
}

// UnpackCHECKINTYPEHASH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb48c7f49.
//
// Solidity: function CHECKIN_TYPEHASH() view returns(bytes32)
func (ticketNFT *TicketNFT) UnpackCHECKINTYPEHASH(data []byte) ([32]byte, error) {
	out, err := ticketNFT.abi.Unpack("CHECKIN_TYPEHASH", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) PackDEFAULTADMINROLE() []byte {
	enc, err := ticketNFT.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return ticketNFT.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := ticketNFT.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGATEROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf061368.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function GATE_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) PackGATEROLE() []byte {
	enc, err := ticketNFT.abi.Pack("GATE_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGATEROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf061368.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function GATE_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) TryPackGATEROLE() ([]byte, error) {
	return ticketNFT.abi.Pack("GATE_ROLE")
}

// UnpackGATEROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdf061368.
//
// Solidity: function GATE_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) UnpackGATEROLE(data []byte) ([32]byte, error) {
	out, err := ticketNFT.abi.Unpack("GATE_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackMINTERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5391393.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) PackMINTERROLE() []byte {
	enc, err := ticketNFT.abi.Pack("MINTER_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINTERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5391393.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) TryPackMINTERROLE() ([]byte, error) {
	return ticketNFT.abi.Pack("MINTER_ROLE")
}

// UnpackMINTERROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (ticketNFT *TicketNFT) UnpackMINTERROLE(data []byte) ([32]byte, error) {
	out, err := ticketNFT.abi.Unpack("MINTER_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address , uint256 ) pure returns()
func (ticketNFT *TicketNFT) PackApprove(arg0 common.Address, arg1 *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("approve", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address , uint256 ) pure returns()
func (ticketNFT *TicketNFT) TryPackApprove(arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("approve", arg0, arg1)
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (ticketNFT *TicketNFT) PackBalanceOf(owner common.Address) []byte {
	enc, err := ticketNFT.abi.Pack("balanceOf", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (ticketNFT *TicketNFT) TryPackBalanceOf(owner common.Address) ([]byte, error) {
	return ticketNFT.abi.Pack("balanceOf", owner)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (ticketNFT *TicketNFT) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := ticketNFT.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBuyResale is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa89d2d38.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function buyResale(uint256 tokenId, uint256 maxPrice) returns()
func (ticketNFT *TicketNFT) PackBuyResale(tokenId *big.Int, maxPrice *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("buyResale", tokenId, maxPrice)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBuyResale is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa89d2d38.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function buyResale(uint256 tokenId, uint256 maxPrice) returns()
func (ticketNFT *TicketNFT) TryPackBuyResale(tokenId *big.Int, maxPrice *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("buyResale", tokenId, maxPrice)
}

// PackCancelListing is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x305a67a8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelListing(uint256 tokenId) returns()
func (ticketNFT *TicketNFT) PackCancelListing(tokenId *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("cancelListing", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelListing is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x305a67a8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelListing(uint256 tokenId) returns()
func (ticketNFT *TicketNFT) TryPackCancelListing(tokenId *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("cancelListing", tokenId)
}

// PackCheckIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfd649abf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkIn(uint256 tokenId, uint256 deadline, bytes holderSig) returns()
func (ticketNFT *TicketNFT) PackCheckIn(tokenId *big.Int, deadline *big.Int, holderSig []byte) []byte {
	enc, err := ticketNFT.abi.Pack("checkIn", tokenId, deadline, holderSig)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfd649abf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkIn(uint256 tokenId, uint256 deadline, bytes holderSig) returns()
func (ticketNFT *TicketNFT) TryPackCheckIn(tokenId *big.Int, deadline *big.Int, holderSig []byte) ([]byte, error) {
	return ticketNFT.abi.Pack("checkIn", tokenId, deadline, holderSig)
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (ticketNFT *TicketNFT) PackEip712Domain() []byte {
	enc, err := ticketNFT.abi.Pack("eip712Domain")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (ticketNFT *TicketNFT) TryPackEip712Domain() ([]byte, error) {
	return ticketNFT.abi.Pack("eip712Domain")
}

// Eip712DomainOutput serves as a container for the return parameters of contract
// method Eip712Domain.
type Eip712DomainOutput struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

// UnpackEip712Domain is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (ticketNFT *TicketNFT) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := ticketNFT.abi.Unpack("eip712Domain", data)
	outstruct := new(Eip712DomainOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, nil
}

// PackEventStart is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x202887e7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eventStart() view returns(uint64)
func (ticketNFT *TicketNFT) PackEventStart() []byte {
	enc, err := ticketNFT.abi.Pack("eventStart")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEventStart is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x202887e7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eventStart() view returns(uint64)
func (ticketNFT *TicketNFT) TryPackEventStart() ([]byte, error) {
	return ticketNFT.abi.Pack("eventStart")
}

// UnpackEventStart is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x202887e7.
//
// Solidity: function eventStart() view returns(uint64)
func (ticketNFT *TicketNFT) UnpackEventStart(data []byte) (uint64, error) {
	out, err := ticketNFT.abi.Unpack("eventStart", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackGetApproved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081812fc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (ticketNFT *TicketNFT) PackGetApproved(tokenId *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("getApproved", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetApproved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081812fc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (ticketNFT *TicketNFT) TryPackGetApproved(tokenId *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("getApproved", tokenId)
}

// UnpackGetApproved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (ticketNFT *TicketNFT) UnpackGetApproved(data []byte) (common.Address, error) {
	out, err := ticketNFT.abi.Unpack("getApproved", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (ticketNFT *TicketNFT) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := ticketNFT.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (ticketNFT *TicketNFT) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return ticketNFT.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (ticketNFT *TicketNFT) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := ticketNFT.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (ticketNFT *TicketNFT) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := ticketNFT.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (ticketNFT *TicketNFT) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return ticketNFT.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (ticketNFT *TicketNFT) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := ticketNFT.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (ticketNFT *TicketNFT) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return ticketNFT.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (ticketNFT *TicketNFT) UnpackHasRole(data []byte) (bool, error) {
	out, err := ticketNFT.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (ticketNFT *TicketNFT) PackIsApprovedForAll(owner common.Address, operator common.Address) []byte {
	enc, err := ticketNFT.abi.Pack("isApprovedForAll", owner, operator)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (ticketNFT *TicketNFT) TryPackIsApprovedForAll(owner common.Address, operator common.Address) ([]byte, error) {
	return ticketNFT.abi.Pack("isApprovedForAll", owner, operator)
}

// UnpackIsApprovedForAll is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (ticketNFT *TicketNFT) UnpackIsApprovedForAll(data []byte) (bool, error) {
	out, err := ticketNFT.abi.Unpack("isApprovedForAll", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackListForResale is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f7d863c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function listForResale(uint256 tokenId, uint256 price) returns()
func (ticketNFT *TicketNFT) PackListForResale(tokenId *big.Int, price *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("listForResale", tokenId, price)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackListForResale is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f7d863c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function listForResale(uint256 tokenId, uint256 price) returns()
func (ticketNFT *TicketNFT) TryPackListForResale(tokenId *big.Int, price *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("listForResale", tokenId, price)
}

// PackMaxResaleBps is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x923aee8e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxResaleBps() view returns(uint32)
func (ticketNFT *TicketNFT) PackMaxResaleBps() []byte {
	enc, err := ticketNFT.abi.Pack("maxResaleBps")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxResaleBps is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x923aee8e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxResaleBps() view returns(uint32)
func (ticketNFT *TicketNFT) TryPackMaxResaleBps() ([]byte, error) {
	return ticketNFT.abi.Pack("maxResaleBps")
}

// UnpackMaxResaleBps is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x923aee8e.
//
// Solidity: function maxResaleBps() view returns(uint32)
func (ticketNFT *TicketNFT) UnpackMaxResaleBps(data []byte) (uint32, error) {
	out, err := ticketNFT.abi.Unpack("maxResaleBps", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackMaxSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5abeb01.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxSupply() view returns(uint32)
func (ticketNFT *TicketNFT) PackMaxSupply() []byte {
	enc, err := ticketNFT.abi.Pack("maxSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5abeb01.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxSupply() view returns(uint32)
func (ticketNFT *TicketNFT) TryPackMaxSupply() ([]byte, error) {
	return ticketNFT.abi.Pack("maxSupply")
}

// UnpackMaxSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd5abeb01.
//
// Solidity: function maxSupply() view returns(uint32)
func (ticketNFT *TicketNFT) UnpackMaxSupply(data []byte) (uint32, error) {
	out, err := ticketNFT.abi.Unpack("maxSupply", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackMintTicket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x59e1c251.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mintTicket(address to, uint256 tokenId, uint256 faceValue, bytes32 recordHash) returns()
func (ticketNFT *TicketNFT) PackMintTicket(to common.Address, tokenId *big.Int, faceValue *big.Int, recordHash [32]byte) []byte {
	enc, err := ticketNFT.abi.Pack("mintTicket", to, tokenId, faceValue, recordHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMintTicket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x59e1c251.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mintTicket(address to, uint256 tokenId, uint256 faceValue, bytes32 recordHash) returns()
func (ticketNFT *TicketNFT) TryPackMintTicket(to common.Address, tokenId *big.Int, faceValue *big.Int, recordHash [32]byte) ([]byte, error) {
	return ticketNFT.abi.Pack("mintTicket", to, tokenId, faceValue, recordHash)
}

// PackMintedCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf721b15.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mintedCount() view returns(uint32)
func (ticketNFT *TicketNFT) PackMintedCount() []byte {
	enc, err := ticketNFT.abi.Pack("mintedCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMintedCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf721b15.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mintedCount() view returns(uint32)
func (ticketNFT *TicketNFT) TryPackMintedCount() ([]byte, error) {
	return ticketNFT.abi.Pack("mintedCount")
}

// UnpackMintedCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcf721b15.
//
// Solidity: function mintedCount() view returns(uint32)
func (ticketNFT *TicketNFT) UnpackMintedCount(data []byte) (uint32, error) {
	out, err := ticketNFT.abi.Unpack("mintedCount", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (ticketNFT *TicketNFT) PackName() []byte {
	enc, err := ticketNFT.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function name() view returns(string)
func (ticketNFT *TicketNFT) TryPackName() ([]byte, error) {
	return ticketNFT.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (ticketNFT *TicketNFT) UnpackName(data []byte) (string, error) {
	out, err := ticketNFT.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackOwnerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6352211e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (ticketNFT *TicketNFT) PackOwnerOf(tokenId *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("ownerOf", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwnerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6352211e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (ticketNFT *TicketNFT) TryPackOwnerOf(tokenId *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("ownerOf", tokenId)
}

// UnpackOwnerOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (ticketNFT *TicketNFT) UnpackOwnerOf(data []byte) (common.Address, error) {
	out, err := ticketNFT.abi.Unpack("ownerOf", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (ticketNFT *TicketNFT) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := ticketNFT.abi.Pack("renounceRole", role, callerConfirmation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (ticketNFT *TicketNFT) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return ticketNFT.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackResaleCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0e0caabf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function resaleCap(uint256 tokenId) view returns(uint256)
func (ticketNFT *TicketNFT) PackResaleCap(tokenId *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("resaleCap", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackResaleCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0e0caabf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function resaleCap(uint256 tokenId) view returns(uint256)
func (ticketNFT *TicketNFT) TryPackResaleCap(tokenId *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("resaleCap", tokenId)
}

// UnpackResaleCap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0e0caabf.
//
// Solidity: function resaleCap(uint256 tokenId) view returns(uint256)
func (ticketNFT *TicketNFT) UnpackResaleCap(data []byte) (*big.Int, error) {
	out, err := ticketNFT.abi.Unpack("resaleCap", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (ticketNFT *TicketNFT) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := ticketNFT.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (ticketNFT *TicketNFT) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return ticketNFT.abi.Pack("revokeRole", role, account)
}

// PackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42842e0e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (ticketNFT *TicketNFT) PackSafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("safeTransferFrom", from, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42842e0e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (ticketNFT *TicketNFT) TryPackSafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("safeTransferFrom", from, to, tokenId)
}

// PackSafeTransferFrom0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb88d4fde.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeTransferFrom(address , address , uint256 , bytes ) pure returns()
func (ticketNFT *TicketNFT) PackSafeTransferFrom0(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) []byte {
	enc, err := ticketNFT.abi.Pack("safeTransferFrom0", arg0, arg1, arg2, arg3)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSafeTransferFrom0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb88d4fde.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function safeTransferFrom(address , address , uint256 , bytes ) pure returns()
func (ticketNFT *TicketNFT) TryPackSafeTransferFrom0(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([]byte, error) {
	return ticketNFT.abi.Pack("safeTransferFrom0", arg0, arg1, arg2, arg3)
}

// PackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setApprovalForAll(address , bool ) pure returns()
func (ticketNFT *TicketNFT) PackSetApprovalForAll(arg0 common.Address, arg1 bool) []byte {
	enc, err := ticketNFT.abi.Pack("setApprovalForAll", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setApprovalForAll(address , bool ) pure returns()
func (ticketNFT *TicketNFT) TryPackSetApprovalForAll(arg0 common.Address, arg1 bool) ([]byte, error) {
	return ticketNFT.abi.Pack("setApprovalForAll", arg0, arg1)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (ticketNFT *TicketNFT) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := ticketNFT.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (ticketNFT *TicketNFT) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return ticketNFT.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (ticketNFT *TicketNFT) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := ticketNFT.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (ticketNFT *TicketNFT) PackSymbol() []byte {
	enc, err := ticketNFT.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(string)
func (ticketNFT *TicketNFT) TryPackSymbol() ([]byte, error) {
	return ticketNFT.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (ticketNFT *TicketNFT) UnpackSymbol(data []byte) (string, error) {
	out, err := ticketNFT.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTicketInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa713083.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ticketInfo(uint256 tokenId) view returns(uint256 faceValue, bytes32 recordHash, uint256 listingPrice)
func (ticketNFT *TicketNFT) PackTicketInfo(tokenId *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("ticketInfo", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTicketInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa713083.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ticketInfo(uint256 tokenId) view returns(uint256 faceValue, bytes32 recordHash, uint256 listingPrice)
func (ticketNFT *TicketNFT) TryPackTicketInfo(tokenId *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("ticketInfo", tokenId)
}

// TicketInfoOutput serves as a container for the return parameters of contract
// method TicketInfo.
type TicketInfoOutput struct {
	FaceValue    *big.Int
	RecordHash   [32]byte
	ListingPrice *big.Int
}

// UnpackTicketInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfa713083.
//
// Solidity: function ticketInfo(uint256 tokenId) view returns(uint256 faceValue, bytes32 recordHash, uint256 listingPrice)
func (ticketNFT *TicketNFT) UnpackTicketInfo(data []byte) (TicketInfoOutput, error) {
	out, err := ticketNFT.abi.Unpack("ticketInfo", data)
	outstruct := new(TicketInfoOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.FaceValue = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.RecordHash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ListingPrice = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackTokenURI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87b56dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (ticketNFT *TicketNFT) PackTokenURI(tokenId *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("tokenURI", tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTokenURI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87b56dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (ticketNFT *TicketNFT) TryPackTokenURI(tokenId *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("tokenURI", tokenId)
}

// UnpackTokenURI is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (ticketNFT *TicketNFT) UnpackTokenURI(data []byte) (string, error) {
	out, err := ticketNFT.abi.Unpack("tokenURI", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address , address , uint256 ) pure returns()
func (ticketNFT *TicketNFT) PackTransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) []byte {
	enc, err := ticketNFT.abi.Pack("transferFrom", arg0, arg1, arg2)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address , address , uint256 ) pure returns()
func (ticketNFT *TicketNFT) TryPackTransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) ([]byte, error) {
	return ticketNFT.abi.Pack("transferFrom", arg0, arg1, arg2)
}

// TicketNFTApproval represents a Approval event raised by the TicketNFT contract.
type TicketNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const TicketNFTApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (TicketNFTApproval) ContractEventName() string {
	return TicketNFTApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (ticketNFT *TicketNFT) UnpackApprovalEvent(log *types.Log) (*TicketNFTApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTApproval)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTApprovalForAll represents a ApprovalForAll event raised by the TicketNFT contract.
type TicketNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const TicketNFTApprovalForAllEventName = "ApprovalForAll"

// ContractEventName returns the user-defined event name.
func (TicketNFTApprovalForAll) ContractEventName() string {
	return TicketNFTApprovalForAllEventName
}

// UnpackApprovalForAllEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (ticketNFT *TicketNFT) UnpackApprovalForAllEvent(log *types.Log) (*TicketNFTApprovalForAll, error) {
	event := "ApprovalForAll"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTApprovalForAll)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTCheckedIn represents a CheckedIn event raised by the TicketNFT contract.
type TicketNFTCheckedIn struct {
	TokenId *big.Int
	Holder  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const TicketNFTCheckedInEventName = "CheckedIn"

// ContractEventName returns the user-defined event name.
func (TicketNFTCheckedIn) ContractEventName() string {
	return TicketNFTCheckedInEventName
}

// UnpackCheckedInEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CheckedIn(uint256 indexed tokenId, address indexed holder)
func (ticketNFT *TicketNFT) UnpackCheckedInEvent(log *types.Log) (*TicketNFTCheckedIn, error) {
	event := "CheckedIn"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTCheckedIn)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTEIP712DomainChanged represents a EIP712DomainChanged event raised by the TicketNFT contract.
type TicketNFTEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const TicketNFTEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (TicketNFTEIP712DomainChanged) ContractEventName() string {
	return TicketNFTEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (ticketNFT *TicketNFT) UnpackEIP712DomainChangedEvent(log *types.Log) (*TicketNFTEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTEIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTListingCancelled represents a ListingCancelled event raised by the TicketNFT contract.
type TicketNFTListingCancelled struct {
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const TicketNFTListingCancelledEventName = "ListingCancelled"

// ContractEventName returns the user-defined event name.
func (TicketNFTListingCancelled) ContractEventName() string {
	return TicketNFTListingCancelledEventName
}

// UnpackListingCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ListingCancelled(uint256 indexed tokenId)
func (ticketNFT *TicketNFT) UnpackListingCancelledEvent(log *types.Log) (*TicketNFTListingCancelled, error) {
	event := "ListingCancelled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTListingCancelled)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTRoleAdminChanged represents a RoleAdminChanged event raised by the TicketNFT contract.
type TicketNFTRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const TicketNFTRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (TicketNFTRoleAdminChanged) ContractEventName() string {
	return TicketNFTRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (ticketNFT *TicketNFT) UnpackRoleAdminChangedEvent(log *types.Log) (*TicketNFTRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTRoleGranted represents a RoleGranted event raised by the TicketNFT contract.
type TicketNFTRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const TicketNFTRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (TicketNFTRoleGranted) ContractEventName() string {
	return TicketNFTRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (ticketNFT *TicketNFT) UnpackRoleGrantedEvent(log *types.Log) (*TicketNFTRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTRoleGranted)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTRoleRevoked represents a RoleRevoked event raised by the TicketNFT contract.
type TicketNFTRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const TicketNFTRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (TicketNFTRoleRevoked) ContractEventName() string {
	return TicketNFTRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (ticketNFT *TicketNFT) UnpackRoleRevokedEvent(log *types.Log) (*TicketNFTRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTRoleRevoked)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTTicketListed represents a TicketListed event raised by the TicketNFT contract.
type TicketNFTTicketListed struct {
	TokenId *big.Int
	Seller  common.Address
	Price   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const TicketNFTTicketListedEventName = "TicketListed"

// ContractEventName returns the user-defined event name.
func (TicketNFTTicketListed) ContractEventName() string {
	return TicketNFTTicketListedEventName
}

// UnpackTicketListedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TicketListed(uint256 indexed tokenId, address indexed seller, uint256 price)
func (ticketNFT *TicketNFT) UnpackTicketListedEvent(log *types.Log) (*TicketNFTTicketListed, error) {
	event := "TicketListed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTTicketListed)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTTicketMinted represents a TicketMinted event raised by the TicketNFT contract.
type TicketNFTTicketMinted struct {
	TokenId    *big.Int
	To         common.Address
	FaceValue  *big.Int
	RecordHash [32]byte
	Raw        *types.Log // Blockchain specific contextual infos
}

const TicketNFTTicketMintedEventName = "TicketMinted"

// ContractEventName returns the user-defined event name.
func (TicketNFTTicketMinted) ContractEventName() string {
	return TicketNFTTicketMintedEventName
}

// UnpackTicketMintedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TicketMinted(uint256 indexed tokenId, address indexed to, uint256 faceValue, bytes32 recordHash)
func (ticketNFT *TicketNFT) UnpackTicketMintedEvent(log *types.Log) (*TicketNFTTicketMinted, error) {
	event := "TicketMinted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTTicketMinted)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTTicketResold represents a TicketResold event raised by the TicketNFT contract.
type TicketNFTTicketResold struct {
	TokenId *big.Int
	From    common.Address
	To      common.Address
	Price   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const TicketNFTTicketResoldEventName = "TicketResold"

// ContractEventName returns the user-defined event name.
func (TicketNFTTicketResold) ContractEventName() string {
	return TicketNFTTicketResoldEventName
}

// UnpackTicketResoldEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TicketResold(uint256 indexed tokenId, address indexed from, address indexed to, uint256 price)
func (ticketNFT *TicketNFT) UnpackTicketResoldEvent(log *types.Log) (*TicketNFTTicketResold, error) {
	event := "TicketResold"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTTicketResold)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// TicketNFTTransfer represents a Transfer event raised by the TicketNFT contract.
type TicketNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const TicketNFTTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (TicketNFTTransfer) ContractEventName() string {
	return TicketNFTTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (ticketNFT *TicketNFT) UnpackTransferEvent(log *types.Log) (*TicketNFTTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != ticketNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(TicketNFTTransfer)
	if len(log.Data) > 0 {
		if err := ticketNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range ticketNFT.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (ticketNFT *TicketNFT) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["BuyerIsSeller"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackBuyerIsSellerError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ERC721IncorrectOwner"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackERC721IncorrectOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ERC721InsufficientApproval"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackERC721InsufficientApprovalError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ERC721InvalidApprover"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackERC721InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ERC721InvalidOperator"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackERC721InvalidOperatorError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ERC721InvalidOwner"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackERC721InvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ERC721InvalidReceiver"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackERC721InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ERC721InvalidSender"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackERC721InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ERC721NonexistentToken"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackERC721NonexistentTokenError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["EmptyRecordHash"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackEmptyRecordHashError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["EventStarted"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackEventStartedError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["InvalidConfig"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackInvalidConfigError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["InvalidShortString"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackInvalidShortStringError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["NotListed"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackNotListedError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["NotSignedByHolder"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackNotSignedByHolderError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["NotTicketHolder"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackNotTicketHolderError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["PassExpired"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackPassExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["PriceAboveCap"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackPriceAboveCapError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["PriceAboveMaxPrice"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackPriceAboveMaxPriceError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["SoldOut"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackSoldOutError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["StringTooLong"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackStringTooLongError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["TransfersDisabled"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackTransfersDisabledError(raw[4:])
	}
	if bytes.Equal(raw[:4], ticketNFT.abi.Errors["ZeroPrice"].ID.Bytes()[:4]) {
		return ticketNFT.UnpackZeroPriceError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TicketNFTAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the TicketNFT contract.
type TicketNFTAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func TicketNFTAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (ticketNFT *TicketNFT) UnpackAccessControlBadConfirmationError(raw []byte) (*TicketNFTAccessControlBadConfirmation, error) {
	out := new(TicketNFTAccessControlBadConfirmation)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the TicketNFT contract.
type TicketNFTAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func TicketNFTAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (ticketNFT *TicketNFT) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*TicketNFTAccessControlUnauthorizedAccount, error) {
	out := new(TicketNFTAccessControlUnauthorizedAccount)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTBuyerIsSeller represents a BuyerIsSeller error raised by the TicketNFT contract.
type TicketNFTBuyerIsSeller struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BuyerIsSeller()
func TicketNFTBuyerIsSellerErrorID() common.Hash {
	return common.HexToHash("0x618350596b18803f7b9674dba6324378e0091c153139de70f0c90a97e84e19a5")
}

// UnpackBuyerIsSellerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BuyerIsSeller()
func (ticketNFT *TicketNFT) UnpackBuyerIsSellerError(raw []byte) (*TicketNFTBuyerIsSeller, error) {
	out := new(TicketNFTBuyerIsSeller)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "BuyerIsSeller", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTERC721IncorrectOwner represents a ERC721IncorrectOwner error raised by the TicketNFT contract.
type TicketNFTERC721IncorrectOwner struct {
	Sender  common.Address
	TokenId *big.Int
	Owner   common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)
func TicketNFTERC721IncorrectOwnerErrorID() common.Hash {
	return common.HexToHash("0x64283d7b313c8117c125f736876fa2b4e90ea3831a4716dfdb87d2f540e26289")
}

// UnpackERC721IncorrectOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)
func (ticketNFT *TicketNFT) UnpackERC721IncorrectOwnerError(raw []byte) (*TicketNFTERC721IncorrectOwner, error) {
	out := new(TicketNFTERC721IncorrectOwner)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ERC721IncorrectOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTERC721InsufficientApproval represents a ERC721InsufficientApproval error raised by the TicketNFT contract.
type TicketNFTERC721InsufficientApproval struct {
	Operator common.Address
	TokenId  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InsufficientApproval(address operator, uint256 tokenId)
func TicketNFTERC721InsufficientApprovalErrorID() common.Hash {
	return common.HexToHash("0x177e802f6f313bc89797ecace66d6d29ab4719cbaaacbb87367264048b1eb861")
}

// UnpackERC721InsufficientApprovalError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InsufficientApproval(address operator, uint256 tokenId)
func (ticketNFT *TicketNFT) UnpackERC721InsufficientApprovalError(raw []byte) (*TicketNFTERC721InsufficientApproval, error) {
	out := new(TicketNFTERC721InsufficientApproval)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ERC721InsufficientApproval", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTERC721InvalidApprover represents a ERC721InvalidApprover error raised by the TicketNFT contract.
type TicketNFTERC721InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidApprover(address approver)
func TicketNFTERC721InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xa9fbf51f86b8e03595d59dc726bb10c329bb24f62589be276d8dd193ca0b69ea")
}

// UnpackERC721InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidApprover(address approver)
func (ticketNFT *TicketNFT) UnpackERC721InvalidApproverError(raw []byte) (*TicketNFTERC721InvalidApprover, error) {
	out := new(TicketNFTERC721InvalidApprover)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ERC721InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTERC721InvalidOperator represents a ERC721InvalidOperator error raised by the TicketNFT contract.
type TicketNFTERC721InvalidOperator struct {
	Operator common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidOperator(address operator)
func TicketNFTERC721InvalidOperatorErrorID() common.Hash {
	return common.HexToHash("0x5b08ba185e8f577075361f3a3555a6580a227ce22734dcc979c1aeadf894658b")
}

// UnpackERC721InvalidOperatorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidOperator(address operator)
func (ticketNFT *TicketNFT) UnpackERC721InvalidOperatorError(raw []byte) (*TicketNFTERC721InvalidOperator, error) {
	out := new(TicketNFTERC721InvalidOperator)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ERC721InvalidOperator", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTERC721InvalidOwner represents a ERC721InvalidOwner error raised by the TicketNFT contract.
type TicketNFTERC721InvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidOwner(address owner)
func TicketNFTERC721InvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x89c62b6479af2e623826dcc39c5133061d35b66d72de92833401dd2fd6567480")
}

// UnpackERC721InvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidOwner(address owner)
func (ticketNFT *TicketNFT) UnpackERC721InvalidOwnerError(raw []byte) (*TicketNFTERC721InvalidOwner, error) {
	out := new(TicketNFTERC721InvalidOwner)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ERC721InvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTERC721InvalidReceiver represents a ERC721InvalidReceiver error raised by the TicketNFT contract.
type TicketNFTERC721InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidReceiver(address receiver)
func TicketNFTERC721InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0x64a0ae9278f805eaf991dcd18ca78756d280b7508b764ef1b255c55845c11df9")
}

// UnpackERC721InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidReceiver(address receiver)
func (ticketNFT *TicketNFT) UnpackERC721InvalidReceiverError(raw []byte) (*TicketNFTERC721InvalidReceiver, error) {
	out := new(TicketNFTERC721InvalidReceiver)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ERC721InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTERC721InvalidSender represents a ERC721InvalidSender error raised by the TicketNFT contract.
type TicketNFTERC721InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidSender(address sender)
func TicketNFTERC721InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x73c6ac6e10798e95d99e1f130d923eb40193ecb8d094ec3dce93292564eb3b17")
}

// UnpackERC721InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidSender(address sender)
func (ticketNFT *TicketNFT) UnpackERC721InvalidSenderError(raw []byte) (*TicketNFTERC721InvalidSender, error) {
	out := new(TicketNFTERC721InvalidSender)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ERC721InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTERC721NonexistentToken represents a ERC721NonexistentToken error raised by the TicketNFT contract.
type TicketNFTERC721NonexistentToken struct {
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721NonexistentToken(uint256 tokenId)
func TicketNFTERC721NonexistentTokenErrorID() common.Hash {
	return common.HexToHash("0x7e273289a3a9ef6670f06df7dca227856fc925e956db96980692764a8bc734d7")
}

// UnpackERC721NonexistentTokenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721NonexistentToken(uint256 tokenId)
func (ticketNFT *TicketNFT) UnpackERC721NonexistentTokenError(raw []byte) (*TicketNFTERC721NonexistentToken, error) {
	out := new(TicketNFTERC721NonexistentToken)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ERC721NonexistentToken", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTEmptyRecordHash represents a EmptyRecordHash error raised by the TicketNFT contract.
type TicketNFTEmptyRecordHash struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyRecordHash()
func TicketNFTEmptyRecordHashErrorID() common.Hash {
	return common.HexToHash("0x3da1ae16ccc54aff6dcf02e753311248334fe80db5a081fbec5989ccbb2fd9b5")
}

// UnpackEmptyRecordHashError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyRecordHash()
func (ticketNFT *TicketNFT) UnpackEmptyRecordHashError(raw []byte) (*TicketNFTEmptyRecordHash, error) {
	out := new(TicketNFTEmptyRecordHash)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "EmptyRecordHash", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTEventStarted represents a EventStarted error raised by the TicketNFT contract.
type TicketNFTEventStarted struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EventStarted()
func TicketNFTEventStartedErrorID() common.Hash {
	return common.HexToHash("0x292403c31f9543e41c6458bb9eb5eac3516353be8e91cdd4c3238f478142f059")
}

// UnpackEventStartedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EventStarted()
func (ticketNFT *TicketNFT) UnpackEventStartedError(raw []byte) (*TicketNFTEventStarted, error) {
	out := new(TicketNFTEventStarted)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "EventStarted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTInvalidConfig represents a InvalidConfig error raised by the TicketNFT contract.
type TicketNFTInvalidConfig struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidConfig()
func TicketNFTInvalidConfigErrorID() common.Hash {
	return common.HexToHash("0x35be3ac88eaaf385e6b33d44d849a44bdc3ae215676d5d7e0194dfe20d80c66e")
}

// UnpackInvalidConfigError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidConfig()
func (ticketNFT *TicketNFT) UnpackInvalidConfigError(raw []byte) (*TicketNFTInvalidConfig, error) {
	out := new(TicketNFTInvalidConfig)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "InvalidConfig", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTInvalidShortString represents a InvalidShortString error raised by the TicketNFT contract.
type TicketNFTInvalidShortString struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidShortString()
func TicketNFTInvalidShortStringErrorID() common.Hash {
	return common.HexToHash("0xb3512b0c6163e5f0bafab72bb631b9d58cd7a731b082f910338aa21c83d5c274")
}

// UnpackInvalidShortStringError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidShortString()
func (ticketNFT *TicketNFT) UnpackInvalidShortStringError(raw []byte) (*TicketNFTInvalidShortString, error) {
	out := new(TicketNFTInvalidShortString)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "InvalidShortString", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTNotListed represents a NotListed error raised by the TicketNFT contract.
type TicketNFTNotListed struct {
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotListed(uint256 tokenId)
func TicketNFTNotListedErrorID() common.Hash {
	return common.HexToHash("0xd176cea0d27fd63dcbcf21630f5c3452bee00c238c2ad1142db4ae64e0f26fd0")
}

// UnpackNotListedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotListed(uint256 tokenId)
func (ticketNFT *TicketNFT) UnpackNotListedError(raw []byte) (*TicketNFTNotListed, error) {
	out := new(TicketNFTNotListed)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "NotListed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTNotSignedByHolder represents a NotSignedByHolder error raised by the TicketNFT contract.
type TicketNFTNotSignedByHolder struct {
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotSignedByHolder(uint256 tokenId)
func TicketNFTNotSignedByHolderErrorID() common.Hash {
	return common.HexToHash("0x1cec443d96f4eeee4d54548d8d2797d5cff8f04fa928854eff000a5c34fcd30c")
}

// UnpackNotSignedByHolderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotSignedByHolder(uint256 tokenId)
func (ticketNFT *TicketNFT) UnpackNotSignedByHolderError(raw []byte) (*TicketNFTNotSignedByHolder, error) {
	out := new(TicketNFTNotSignedByHolder)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "NotSignedByHolder", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTNotTicketHolder represents a NotTicketHolder error raised by the TicketNFT contract.
type TicketNFTNotTicketHolder struct {
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotTicketHolder(uint256 tokenId)
func TicketNFTNotTicketHolderErrorID() common.Hash {
	return common.HexToHash("0xb1a257d892b9a020fd0b795fd263c5427eccc65cffedf7dffb475c6a7b2f5d77")
}

// UnpackNotTicketHolderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotTicketHolder(uint256 tokenId)
func (ticketNFT *TicketNFT) UnpackNotTicketHolderError(raw []byte) (*TicketNFTNotTicketHolder, error) {
	out := new(TicketNFTNotTicketHolder)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "NotTicketHolder", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTPassExpired represents a PassExpired error raised by the TicketNFT contract.
type TicketNFTPassExpired struct {
	Deadline *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PassExpired(uint256 deadline)
func TicketNFTPassExpiredErrorID() common.Hash {
	return common.HexToHash("0x19728872b7a2221ddd22f5ba6b83d274f8af01630f9871c7ba61f106cd4f1eee")
}

// UnpackPassExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PassExpired(uint256 deadline)
func (ticketNFT *TicketNFT) UnpackPassExpiredError(raw []byte) (*TicketNFTPassExpired, error) {
	out := new(TicketNFTPassExpired)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "PassExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTPriceAboveCap represents a PriceAboveCap error raised by the TicketNFT contract.
type TicketNFTPriceAboveCap struct {
	Price *big.Int
	Cap   *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PriceAboveCap(uint256 price, uint256 cap)
func TicketNFTPriceAboveCapErrorID() common.Hash {
	return common.HexToHash("0x842c488ab9609c345f75f58201dac86cf45ff680afc0f4f3009f49a8236da33c")
}

// UnpackPriceAboveCapError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PriceAboveCap(uint256 price, uint256 cap)
func (ticketNFT *TicketNFT) UnpackPriceAboveCapError(raw []byte) (*TicketNFTPriceAboveCap, error) {
	out := new(TicketNFTPriceAboveCap)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "PriceAboveCap", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTPriceAboveMaxPrice represents a PriceAboveMaxPrice error raised by the TicketNFT contract.
type TicketNFTPriceAboveMaxPrice struct {
	Price    *big.Int
	MaxPrice *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PriceAboveMaxPrice(uint256 price, uint256 maxPrice)
func TicketNFTPriceAboveMaxPriceErrorID() common.Hash {
	return common.HexToHash("0xf319e5e3575dfacb8b91a7f606229c29969bfd23cb9f9ce803c562b7f56c6301")
}

// UnpackPriceAboveMaxPriceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PriceAboveMaxPrice(uint256 price, uint256 maxPrice)
func (ticketNFT *TicketNFT) UnpackPriceAboveMaxPriceError(raw []byte) (*TicketNFTPriceAboveMaxPrice, error) {
	out := new(TicketNFTPriceAboveMaxPrice)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "PriceAboveMaxPrice", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTSoldOut represents a SoldOut error raised by the TicketNFT contract.
type TicketNFTSoldOut struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SoldOut()
func TicketNFTSoldOutErrorID() common.Hash {
	return common.HexToHash("0x52df9fe5b9c9a7b0b4fdc2c9f89387959e35e4209c2a8d133a2b8165edad2a04")
}

// UnpackSoldOutError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SoldOut()
func (ticketNFT *TicketNFT) UnpackSoldOutError(raw []byte) (*TicketNFTSoldOut, error) {
	out := new(TicketNFTSoldOut)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "SoldOut", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTStringTooLong represents a StringTooLong error raised by the TicketNFT contract.
type TicketNFTStringTooLong struct {
	Str string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringTooLong(string str)
func TicketNFTStringTooLongErrorID() common.Hash {
	return common.HexToHash("0x305a27a93f8e33b7392df0a0f91d6fc63847395853c45991eec52dbf24d72381")
}

// UnpackStringTooLongError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringTooLong(string str)
func (ticketNFT *TicketNFT) UnpackStringTooLongError(raw []byte) (*TicketNFTStringTooLong, error) {
	out := new(TicketNFTStringTooLong)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "StringTooLong", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTTransfersDisabled represents a TransfersDisabled error raised by the TicketNFT contract.
type TicketNFTTransfersDisabled struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TransfersDisabled()
func TicketNFTTransfersDisabledErrorID() common.Hash {
	return common.HexToHash("0x8574adcfb0eb8c46c67db25ca9cd9d7b2e5c6606b1da1f502824783bba70aeef")
}

// UnpackTransfersDisabledError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TransfersDisabled()
func (ticketNFT *TicketNFT) UnpackTransfersDisabledError(raw []byte) (*TicketNFTTransfersDisabled, error) {
	out := new(TicketNFTTransfersDisabled)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "TransfersDisabled", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TicketNFTZeroPrice represents a ZeroPrice error raised by the TicketNFT contract.
type TicketNFTZeroPrice struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroPrice()
func TicketNFTZeroPriceErrorID() common.Hash {
	return common.HexToHash("0x4dfba023c19aaad0ad1005b8a52f2ae0c7a1acb3f45ecd31562b980620d65f93")
}

// UnpackZeroPriceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroPrice()
func (ticketNFT *TicketNFT) UnpackZeroPriceError(raw []byte) (*TicketNFTZeroPrice, error) {
	out := new(TicketNFTZeroPrice)
	if err := ticketNFT.abi.UnpackIntoInterface(out, "ZeroPrice", raw); err != nil {
		return nil, err
	}
	return out, nil
}
