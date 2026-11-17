// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package demonft

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

// DemoNFTMetaData contains all meta data concerning the DemoNFT contract.
var DemoNFTMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeMint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "DemoNFT",
	Bin: "0x608060405234801561000f575f5ffd5b506040516113e03803806113e083398101604081905261002e91610179565b8083835f61003c838261028d565b506001610049828261028d565b5050506001600160a01b03811661007957604051631e4fbdf760e01b81525f600482015260240160405180910390fd5b6100828161008b565b5050505061034b565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b634e487b7160e01b5f52604160045260245ffd5b5f82601f8301126100ff575f5ffd5b81516001600160401b03811115610118576101186100dc565b604051601f8201601f19908116603f011681016001600160401b0381118282101715610146576101466100dc565b60405281815283820160200185101561015d575f5ffd5b8160208501602083015e5f918101602001919091529392505050565b5f5f5f6060848603121561018b575f5ffd5b83516001600160401b038111156101a0575f5ffd5b6101ac868287016100f0565b602086015190945090506001600160401b038111156101c9575f5ffd5b6101d5868287016100f0565b604086015190935090506001600160a01b03811681146101f3575f5ffd5b809150509250925092565b600181811c9082168061021257607f821691505b60208210810361023057634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115610288578282111561028857805f5260205f20601f840160051c602085101561026157505f5b90810190601f840160051c035f5b81811015610284575f8382015560010161026f565b5050505b505050565b81516001600160401b038111156102a6576102a66100dc565b6102ba816102b484546101fe565b84610236565b6020601f8211600181146102ec575f83156102d55750848201515b5f19600385901b1c1916600184901b178455610344565b5f84815260208120601f198516915b8281101561031b57878501518255602094850194600190920191016102fb565b508482101561033857868401515f19600387901b60f8161c191681555b505060018360011b0184555b5050505050565b611088806103585f395ff3fe608060405234801561000f575f5ffd5b5060043610610106575f3560e01c8063715018a61161009e578063a22cb4651161006e578063a22cb46514610215578063b88d4fde14610228578063c87b56dd1461023b578063e985e9c51461024e578063f2fde38b14610261575f5ffd5b8063715018a6146101e15780638da5cb5b146101e957806395d89b41146101fa578063a144819414610202575f5ffd5b806323b872dd116100d957806323b872dd1461018757806342842e0e1461019a5780636352211e146101ad57806370a08231146101c0575f5ffd5b806301ffc9a71461010a57806306fdde0314610132578063081812fc14610147578063095ea7b314610172575b5f5ffd5b61011d610118366004610d35565b610274565b60405190151581526020015b60405180910390f35b61013a6102c5565b6040516101299190610d7e565b61015a610155366004610d90565b610354565b6040516001600160a01b039091168152602001610129565b610185610180366004610dc2565b61037b565b005b610185610195366004610dea565b61038a565b6101856101a8366004610dea565b610418565b61015a6101bb366004610d90565b610437565b6101d36101ce366004610e24565b610441565b604051908152602001610129565b610185610486565b6006546001600160a01b031661015a565b61013a610499565b610185610210366004610dc2565b6104a8565b610185610223366004610e3d565b6104ba565b610185610236366004610e8a565b6104c5565b61013a610249366004610d90565b6104dd565b61011d61025c366004610f67565b61054e565b61018561026f366004610e24565b61057b565b5f6001600160e01b031982166380ac58cd60e01b14806102a457506001600160e01b03198216635b5e139f60e01b145b806102bf57506301ffc9a760e01b6001600160e01b03198316145b92915050565b60605f80546102d390610f98565b80601f01602080910402602001604051908101604052809291908181526020018280546102ff90610f98565b801561034a5780601f106103215761010080835404028352916020019161034a565b820191905f5260205f20905b81548152906001019060200180831161032d57829003601f168201915b5050505050905090565b5f61035e826105b8565b505f828152600460205260409020546001600160a01b03166102bf565b6103868282336105f0565b5050565b6001600160a01b0382166103b857604051633250574960e11b81525f60048201526024015b60405180910390fd5b5f6103c48383336105fd565b9050836001600160a01b0316816001600160a01b031614610412576040516364283d7b60e01b81526001600160a01b03808616600483015260248201849052821660448201526064016103af565b50505050565b61043283838360405180602001604052805f8152506104c5565b505050565b5f6102bf826105b8565b5f6001600160a01b03821661046b576040516322718ad960e21b81525f60048201526024016103af565b506001600160a01b03165f9081526003602052604090205490565b61048e6106ef565b6104975f61071c565b565b6060600180546102d390610f98565b6104b06106ef565b610386828261076d565b610386338383610786565b6104d084848461038a565b610412338585858561084d565b60606104e8826105b8565b505f6104fe60408051602081019091525f815290565b90505f81511161051c5760405180602001604052805f815250610547565b8061052684610975565b604051602001610537929190610fe7565b6040516020818303038152906040525b9392505050565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b6105836106ef565b6001600160a01b0381166105ac57604051631e4fbdf760e01b81525f60048201526024016103af565b6105b58161071c565b50565b5f818152600260205260408120546001600160a01b0316806102bf57604051637e27328960e01b8152600481018490526024016103af565b6104328383836001610a05565b5f828152600260205260408120546001600160a01b039081169083161561062957610629818486610b09565b6001600160a01b03811615610663576106445f855f5f610a05565b6001600160a01b0381165f90815260036020526040902080545f190190555b6001600160a01b03851615610691576001600160a01b0385165f908152600360205260409020805460010190555b5f8481526002602052604080822080546001600160a01b0319166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b6006546001600160a01b031633146104975760405163118cdaa760e01b81523360048201526024016103af565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b610386828260405180602001604052805f815250610b6d565b6001600160a01b0383166107af5760405163a9fbf51f60e01b81525f60048201526024016103af565b6001600160a01b0382166107e157604051630b61174360e31b81526001600160a01b03831660048201526024016103af565b6001600160a01b038381165f81815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0383163b1561096e57604051630a85bd0160e11b81526001600160a01b0384169063150b7a029061088f908890889087908790600401610ffb565b6020604051808303815f875af19250505080156108c9575060408051601f3d908101601f191682019092526108c691810190611037565b60015b610930573d8080156108f6576040519150601f19603f3d011682016040523d82523d5f602084013e6108fb565b606091505b5080515f0361092857604051633250574960e11b81526001600160a01b03851660048201526024016103af565b805160208201fd5b6001600160e01b03198116630a85bd0160e11b1461096c57604051633250574960e11b81526001600160a01b03851660048201526024016103af565b505b5050505050565b60605f61098183610b84565b60010190505f8167ffffffffffffffff8111156109a0576109a0610e76565b6040519080825280601f01601f1916602001820160405280156109ca576020820181803683370190505b5090508181016020015b5f19016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a85049450846109d457509392505050565b8080610a1957506001600160a01b03821615155b15610ada575f610a28846105b8565b90506001600160a01b03831615801590610a545750826001600160a01b0316816001600160a01b031614155b8015610a675750610a65818461054e565b155b15610a905760405163a9fbf51f60e01b81526001600160a01b03841660048201526024016103af565b8115610ad85783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260046020526040902080546001600160a01b0319166001600160a01b0392909216919091179055565b610b14838383610c5b565b610432576001600160a01b038316610b4257604051637e27328960e01b8152600481018290526024016103af565b60405163177e802f60e01b81526001600160a01b0383166004820152602481018290526044016103af565b610b778383610cbf565b610432335f85858561084d565b5f8072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b8310610bc25772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef81000000008310610bee576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc100008310610c0c57662386f26fc10000830492506010015b6305f5e1008310610c24576305f5e100830492506008015b6127108310610c3857612710830492506004015b60648310610c4a576064830492506002015b600a83106102bf5760010192915050565b5f6001600160a01b03831615801590610cb75750826001600160a01b0316846001600160a01b03161480610c945750610c94848461054e565b80610cb757505f828152600460205260409020546001600160a01b038481169116145b949350505050565b6001600160a01b038216610ce857604051633250574960e11b81525f60048201526024016103af565b5f610cf483835f6105fd565b90506001600160a01b03811615610432576040516339e3563760e11b81525f60048201526024016103af565b6001600160e01b0319811681146105b5575f5ffd5b5f60208284031215610d45575f5ffd5b813561054781610d20565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f6105476020830184610d50565b5f60208284031215610da0575f5ffd5b5035919050565b80356001600160a01b0381168114610dbd575f5ffd5b919050565b5f5f60408385031215610dd3575f5ffd5b610ddc83610da7565b946020939093013593505050565b5f5f5f60608486031215610dfc575f5ffd5b610e0584610da7565b9250610e1360208501610da7565b929592945050506040919091013590565b5f60208284031215610e34575f5ffd5b61054782610da7565b5f5f60408385031215610e4e575f5ffd5b610e5783610da7565b915060208301358015158114610e6b575f5ffd5b809150509250929050565b634e487b7160e01b5f52604160045260245ffd5b5f5f5f5f60808587031215610e9d575f5ffd5b610ea685610da7565b9350610eb460208601610da7565b925060408501359150606085013567ffffffffffffffff811115610ed6575f5ffd5b8501601f81018713610ee6575f5ffd5b803567ffffffffffffffff811115610f0057610f00610e76565b604051601f8201601f19908116603f0116810167ffffffffffffffff81118282101715610f2f57610f2f610e76565b604052818152828201602001891015610f46575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b5f5f60408385031215610f78575f5ffd5b610f8183610da7565b9150610f8f60208401610da7565b90509250929050565b600181811c90821680610fac57607f821691505b602082108103610fca57634e487b7160e01b5f52602260045260245ffd5b50919050565b5f81518060208401855e5f93019283525090919050565b5f610cb7610ff58386610fd0565b84610fd0565b6001600160a01b03858116825284166020820152604081018390526080606082018190525f9061102d90830184610d50565b9695505050505050565b5f60208284031215611047575f5ffd5b815161054781610d2056fea264697066735822122035bd63a5b332b20ee619a68dca264779fbcca03e8a176c98d7bd9a32fb9fb10864736f6c63430008260033",
}

// DemoNFT is an auto generated Go binding around an Ethereum contract.
type DemoNFT struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *DemoNFT) GetABI() abi.ABI {
	return c.abi
}

// NewDemoNFT creates a new instance of DemoNFT.
func NewDemoNFT() *DemoNFT {
	parsed, err := DemoNFTMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &DemoNFT{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *DemoNFT) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(string name_, string symbol_, address owner_) returns()
func (demoNFT *DemoNFT) PackConstructor(name_ string, symbol_ string, owner_ common.Address) []byte {
	enc, err := demoNFT.abi.Pack("", name_, symbol_, owner_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (demoNFT *DemoNFT) PackApprove(to common.Address, tokenId *big.Int) []byte {
	enc, err := demoNFT.abi.Pack("approve", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (demoNFT *DemoNFT) TryPackApprove(to common.Address, tokenId *big.Int) ([]byte, error) {
	return demoNFT.abi.Pack("approve", to, tokenId)
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (demoNFT *DemoNFT) PackBalanceOf(owner common.Address) []byte {
	enc, err := demoNFT.abi.Pack("balanceOf", owner)
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
func (demoNFT *DemoNFT) TryPackBalanceOf(owner common.Address) ([]byte, error) {
	return demoNFT.abi.Pack("balanceOf", owner)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (demoNFT *DemoNFT) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := demoNFT.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetApproved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x081812fc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (demoNFT *DemoNFT) PackGetApproved(tokenId *big.Int) []byte {
	enc, err := demoNFT.abi.Pack("getApproved", tokenId)
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
func (demoNFT *DemoNFT) TryPackGetApproved(tokenId *big.Int) ([]byte, error) {
	return demoNFT.abi.Pack("getApproved", tokenId)
}

// UnpackGetApproved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (demoNFT *DemoNFT) UnpackGetApproved(data []byte) (common.Address, error) {
	out, err := demoNFT.abi.Unpack("getApproved", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackIsApprovedForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe985e9c5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (demoNFT *DemoNFT) PackIsApprovedForAll(owner common.Address, operator common.Address) []byte {
	enc, err := demoNFT.abi.Pack("isApprovedForAll", owner, operator)
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
func (demoNFT *DemoNFT) TryPackIsApprovedForAll(owner common.Address, operator common.Address) ([]byte, error) {
	return demoNFT.abi.Pack("isApprovedForAll", owner, operator)
}

// UnpackIsApprovedForAll is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (demoNFT *DemoNFT) UnpackIsApprovedForAll(data []byte) (bool, error) {
	out, err := demoNFT.abi.Unpack("isApprovedForAll", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (demoNFT *DemoNFT) PackName() []byte {
	enc, err := demoNFT.abi.Pack("name")
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
func (demoNFT *DemoNFT) TryPackName() ([]byte, error) {
	return demoNFT.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (demoNFT *DemoNFT) UnpackName(data []byte) (string, error) {
	out, err := demoNFT.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (demoNFT *DemoNFT) PackOwner() []byte {
	enc, err := demoNFT.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (demoNFT *DemoNFT) TryPackOwner() ([]byte, error) {
	return demoNFT.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (demoNFT *DemoNFT) UnpackOwner(data []byte) (common.Address, error) {
	out, err := demoNFT.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOwnerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6352211e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (demoNFT *DemoNFT) PackOwnerOf(tokenId *big.Int) []byte {
	enc, err := demoNFT.abi.Pack("ownerOf", tokenId)
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
func (demoNFT *DemoNFT) TryPackOwnerOf(tokenId *big.Int) ([]byte, error) {
	return demoNFT.abi.Pack("ownerOf", tokenId)
}

// UnpackOwnerOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (demoNFT *DemoNFT) UnpackOwnerOf(data []byte) (common.Address, error) {
	out, err := demoNFT.abi.Unpack("ownerOf", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (demoNFT *DemoNFT) PackRenounceOwnership() []byte {
	enc, err := demoNFT.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() returns()
func (demoNFT *DemoNFT) TryPackRenounceOwnership() ([]byte, error) {
	return demoNFT.abi.Pack("renounceOwnership")
}

// PackSafeMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1448194.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeMint(address to, uint256 tokenId) returns()
func (demoNFT *DemoNFT) PackSafeMint(to common.Address, tokenId *big.Int) []byte {
	enc, err := demoNFT.abi.Pack("safeMint", to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSafeMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1448194.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function safeMint(address to, uint256 tokenId) returns()
func (demoNFT *DemoNFT) TryPackSafeMint(to common.Address, tokenId *big.Int) ([]byte, error) {
	return demoNFT.abi.Pack("safeMint", to, tokenId)
}

// PackSafeTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42842e0e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (demoNFT *DemoNFT) PackSafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := demoNFT.abi.Pack("safeTransferFrom", from, to, tokenId)
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
func (demoNFT *DemoNFT) TryPackSafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) ([]byte, error) {
	return demoNFT.abi.Pack("safeTransferFrom", from, to, tokenId)
}

// PackSafeTransferFrom0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb88d4fde.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (demoNFT *DemoNFT) PackSafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) []byte {
	enc, err := demoNFT.abi.Pack("safeTransferFrom0", from, to, tokenId, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSafeTransferFrom0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb88d4fde.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (demoNFT *DemoNFT) TryPackSafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) ([]byte, error) {
	return demoNFT.abi.Pack("safeTransferFrom0", from, to, tokenId, data)
}

// PackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (demoNFT *DemoNFT) PackSetApprovalForAll(operator common.Address, approved bool) []byte {
	enc, err := demoNFT.abi.Pack("setApprovalForAll", operator, approved)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetApprovalForAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa22cb465.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (demoNFT *DemoNFT) TryPackSetApprovalForAll(operator common.Address, approved bool) ([]byte, error) {
	return demoNFT.abi.Pack("setApprovalForAll", operator, approved)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (demoNFT *DemoNFT) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := demoNFT.abi.Pack("supportsInterface", interfaceId)
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
func (demoNFT *DemoNFT) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return demoNFT.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (demoNFT *DemoNFT) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := demoNFT.abi.Unpack("supportsInterface", data)
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
func (demoNFT *DemoNFT) PackSymbol() []byte {
	enc, err := demoNFT.abi.Pack("symbol")
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
func (demoNFT *DemoNFT) TryPackSymbol() ([]byte, error) {
	return demoNFT.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (demoNFT *DemoNFT) UnpackSymbol(data []byte) (string, error) {
	out, err := demoNFT.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTokenURI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87b56dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (demoNFT *DemoNFT) PackTokenURI(tokenId *big.Int) []byte {
	enc, err := demoNFT.abi.Pack("tokenURI", tokenId)
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
func (demoNFT *DemoNFT) TryPackTokenURI(tokenId *big.Int) ([]byte, error) {
	return demoNFT.abi.Pack("tokenURI", tokenId)
}

// UnpackTokenURI is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (demoNFT *DemoNFT) UnpackTokenURI(data []byte) (string, error) {
	out, err := demoNFT.abi.Unpack("tokenURI", data)
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
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (demoNFT *DemoNFT) PackTransferFrom(from common.Address, to common.Address, tokenId *big.Int) []byte {
	enc, err := demoNFT.abi.Pack("transferFrom", from, to, tokenId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (demoNFT *DemoNFT) TryPackTransferFrom(from common.Address, to common.Address, tokenId *big.Int) ([]byte, error) {
	return demoNFT.abi.Pack("transferFrom", from, to, tokenId)
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (demoNFT *DemoNFT) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := demoNFT.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (demoNFT *DemoNFT) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return demoNFT.abi.Pack("transferOwnership", newOwner)
}

// DemoNFTApproval represents a Approval event raised by the DemoNFT contract.
type DemoNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const DemoNFTApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (DemoNFTApproval) ContractEventName() string {
	return DemoNFTApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (demoNFT *DemoNFT) UnpackApprovalEvent(log *types.Log) (*DemoNFTApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != demoNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DemoNFTApproval)
	if len(log.Data) > 0 {
		if err := demoNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range demoNFT.abi.Events[event].Inputs {
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

// DemoNFTApprovalForAll represents a ApprovalForAll event raised by the DemoNFT contract.
type DemoNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const DemoNFTApprovalForAllEventName = "ApprovalForAll"

// ContractEventName returns the user-defined event name.
func (DemoNFTApprovalForAll) ContractEventName() string {
	return DemoNFTApprovalForAllEventName
}

// UnpackApprovalForAllEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (demoNFT *DemoNFT) UnpackApprovalForAllEvent(log *types.Log) (*DemoNFTApprovalForAll, error) {
	event := "ApprovalForAll"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != demoNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DemoNFTApprovalForAll)
	if len(log.Data) > 0 {
		if err := demoNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range demoNFT.abi.Events[event].Inputs {
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

// DemoNFTOwnershipTransferred represents a OwnershipTransferred event raised by the DemoNFT contract.
type DemoNFTOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const DemoNFTOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (DemoNFTOwnershipTransferred) ContractEventName() string {
	return DemoNFTOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (demoNFT *DemoNFT) UnpackOwnershipTransferredEvent(log *types.Log) (*DemoNFTOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != demoNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DemoNFTOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := demoNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range demoNFT.abi.Events[event].Inputs {
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

// DemoNFTTransfer represents a Transfer event raised by the DemoNFT contract.
type DemoNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const DemoNFTTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (DemoNFTTransfer) ContractEventName() string {
	return DemoNFTTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (demoNFT *DemoNFT) UnpackTransferEvent(log *types.Log) (*DemoNFTTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != demoNFT.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DemoNFTTransfer)
	if len(log.Data) > 0 {
		if err := demoNFT.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range demoNFT.abi.Events[event].Inputs {
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
func (demoNFT *DemoNFT) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["ERC721IncorrectOwner"].ID.Bytes()[:4]) {
		return demoNFT.UnpackERC721IncorrectOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["ERC721InsufficientApproval"].ID.Bytes()[:4]) {
		return demoNFT.UnpackERC721InsufficientApprovalError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["ERC721InvalidApprover"].ID.Bytes()[:4]) {
		return demoNFT.UnpackERC721InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["ERC721InvalidOperator"].ID.Bytes()[:4]) {
		return demoNFT.UnpackERC721InvalidOperatorError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["ERC721InvalidOwner"].ID.Bytes()[:4]) {
		return demoNFT.UnpackERC721InvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["ERC721InvalidReceiver"].ID.Bytes()[:4]) {
		return demoNFT.UnpackERC721InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["ERC721InvalidSender"].ID.Bytes()[:4]) {
		return demoNFT.UnpackERC721InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["ERC721NonexistentToken"].ID.Bytes()[:4]) {
		return demoNFT.UnpackERC721NonexistentTokenError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return demoNFT.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], demoNFT.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return demoNFT.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// DemoNFTERC721IncorrectOwner represents a ERC721IncorrectOwner error raised by the DemoNFT contract.
type DemoNFTERC721IncorrectOwner struct {
	Sender  common.Address
	TokenId *big.Int
	Owner   common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)
func DemoNFTERC721IncorrectOwnerErrorID() common.Hash {
	return common.HexToHash("0x64283d7b313c8117c125f736876fa2b4e90ea3831a4716dfdb87d2f540e26289")
}

// UnpackERC721IncorrectOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)
func (demoNFT *DemoNFT) UnpackERC721IncorrectOwnerError(raw []byte) (*DemoNFTERC721IncorrectOwner, error) {
	out := new(DemoNFTERC721IncorrectOwner)
	if err := demoNFT.abi.UnpackIntoInterface(out, "ERC721IncorrectOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTERC721InsufficientApproval represents a ERC721InsufficientApproval error raised by the DemoNFT contract.
type DemoNFTERC721InsufficientApproval struct {
	Operator common.Address
	TokenId  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InsufficientApproval(address operator, uint256 tokenId)
func DemoNFTERC721InsufficientApprovalErrorID() common.Hash {
	return common.HexToHash("0x177e802f6f313bc89797ecace66d6d29ab4719cbaaacbb87367264048b1eb861")
}

// UnpackERC721InsufficientApprovalError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InsufficientApproval(address operator, uint256 tokenId)
func (demoNFT *DemoNFT) UnpackERC721InsufficientApprovalError(raw []byte) (*DemoNFTERC721InsufficientApproval, error) {
	out := new(DemoNFTERC721InsufficientApproval)
	if err := demoNFT.abi.UnpackIntoInterface(out, "ERC721InsufficientApproval", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTERC721InvalidApprover represents a ERC721InvalidApprover error raised by the DemoNFT contract.
type DemoNFTERC721InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidApprover(address approver)
func DemoNFTERC721InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xa9fbf51f86b8e03595d59dc726bb10c329bb24f62589be276d8dd193ca0b69ea")
}

// UnpackERC721InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidApprover(address approver)
func (demoNFT *DemoNFT) UnpackERC721InvalidApproverError(raw []byte) (*DemoNFTERC721InvalidApprover, error) {
	out := new(DemoNFTERC721InvalidApprover)
	if err := demoNFT.abi.UnpackIntoInterface(out, "ERC721InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTERC721InvalidOperator represents a ERC721InvalidOperator error raised by the DemoNFT contract.
type DemoNFTERC721InvalidOperator struct {
	Operator common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidOperator(address operator)
func DemoNFTERC721InvalidOperatorErrorID() common.Hash {
	return common.HexToHash("0x5b08ba185e8f577075361f3a3555a6580a227ce22734dcc979c1aeadf894658b")
}

// UnpackERC721InvalidOperatorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidOperator(address operator)
func (demoNFT *DemoNFT) UnpackERC721InvalidOperatorError(raw []byte) (*DemoNFTERC721InvalidOperator, error) {
	out := new(DemoNFTERC721InvalidOperator)
	if err := demoNFT.abi.UnpackIntoInterface(out, "ERC721InvalidOperator", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTERC721InvalidOwner represents a ERC721InvalidOwner error raised by the DemoNFT contract.
type DemoNFTERC721InvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidOwner(address owner)
func DemoNFTERC721InvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x89c62b6479af2e623826dcc39c5133061d35b66d72de92833401dd2fd6567480")
}

// UnpackERC721InvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidOwner(address owner)
func (demoNFT *DemoNFT) UnpackERC721InvalidOwnerError(raw []byte) (*DemoNFTERC721InvalidOwner, error) {
	out := new(DemoNFTERC721InvalidOwner)
	if err := demoNFT.abi.UnpackIntoInterface(out, "ERC721InvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTERC721InvalidReceiver represents a ERC721InvalidReceiver error raised by the DemoNFT contract.
type DemoNFTERC721InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidReceiver(address receiver)
func DemoNFTERC721InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0x64a0ae9278f805eaf991dcd18ca78756d280b7508b764ef1b255c55845c11df9")
}

// UnpackERC721InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidReceiver(address receiver)
func (demoNFT *DemoNFT) UnpackERC721InvalidReceiverError(raw []byte) (*DemoNFTERC721InvalidReceiver, error) {
	out := new(DemoNFTERC721InvalidReceiver)
	if err := demoNFT.abi.UnpackIntoInterface(out, "ERC721InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTERC721InvalidSender represents a ERC721InvalidSender error raised by the DemoNFT contract.
type DemoNFTERC721InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721InvalidSender(address sender)
func DemoNFTERC721InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x73c6ac6e10798e95d99e1f130d923eb40193ecb8d094ec3dce93292564eb3b17")
}

// UnpackERC721InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721InvalidSender(address sender)
func (demoNFT *DemoNFT) UnpackERC721InvalidSenderError(raw []byte) (*DemoNFTERC721InvalidSender, error) {
	out := new(DemoNFTERC721InvalidSender)
	if err := demoNFT.abi.UnpackIntoInterface(out, "ERC721InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTERC721NonexistentToken represents a ERC721NonexistentToken error raised by the DemoNFT contract.
type DemoNFTERC721NonexistentToken struct {
	TokenId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC721NonexistentToken(uint256 tokenId)
func DemoNFTERC721NonexistentTokenErrorID() common.Hash {
	return common.HexToHash("0x7e273289a3a9ef6670f06df7dca227856fc925e956db96980692764a8bc734d7")
}

// UnpackERC721NonexistentTokenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC721NonexistentToken(uint256 tokenId)
func (demoNFT *DemoNFT) UnpackERC721NonexistentTokenError(raw []byte) (*DemoNFTERC721NonexistentToken, error) {
	out := new(DemoNFTERC721NonexistentToken)
	if err := demoNFT.abi.UnpackIntoInterface(out, "ERC721NonexistentToken", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the DemoNFT contract.
type DemoNFTOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func DemoNFTOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (demoNFT *DemoNFT) UnpackOwnableInvalidOwnerError(raw []byte) (*DemoNFTOwnableInvalidOwner, error) {
	out := new(DemoNFTOwnableInvalidOwner)
	if err := demoNFT.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoNFTOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the DemoNFT contract.
type DemoNFTOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func DemoNFTOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (demoNFT *DemoNFT) UnpackOwnableUnauthorizedAccountError(raw []byte) (*DemoNFTOwnableUnauthorizedAccount, error) {
	out := new(DemoNFTOwnableUnauthorizedAccount)
	if err := demoNFT.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}
