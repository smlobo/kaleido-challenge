package chain

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Signer holds one private key and signs EIP-1559 transactions for a fixed chain.
// The service builds and signs every transaction itself; nothing is delegated to
// the node (no eth_sendTransaction, no unlocked accounts).
type Signer struct {
	key     *ecdsa.PrivateKey
	address common.Address
	chainID *big.Int
	signer  types.Signer
}

// NewSigner wraps a private key for the given chain.
func NewSigner(key *ecdsa.PrivateKey, chainID *big.Int) *Signer {
	return &Signer{
		key:     key,
		address: crypto.PubkeyToAddress(key.PublicKey),
		chainID: new(big.Int).Set(chainID),
		signer:  types.LatestSignerForChainID(chainID),
	}
}

// NewSignerFromHex parses a hex private key (with or without 0x).
func NewSignerFromHex(hexKey string, chainID *big.Int) (*Signer, error) {
	key, err := crypto.HexToECDSA(strings.TrimPrefix(strings.TrimSpace(hexKey), "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	return NewSigner(key, chainID), nil
}

// GenerateSigner creates a signer with a fresh random key.
func GenerateSigner(chainID *big.Int) (*Signer, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return NewSigner(key, chainID), nil
}

// Address is the account this signer sends from.
func (s *Signer) Address() common.Address { return s.address }

// SignHash signs a 32-byte digest, such as an EIP-712 message hash, and
// returns the 65-byte r ‖ s ‖ v signature that Solidity's ecrecover and
// OpenZeppelin's ECDSA expect (v is 27 or 28). It sends nothing.
func (s *Signer) SignHash(digest common.Hash) ([]byte, error) {
	sig, err := crypto.Sign(digest[:], s.key)
	if err != nil {
		return nil, err
	}
	sig[64] += 27 // crypto.Sign gives v as 0 or 1
	return sig, nil
}

// TxRequest describes an unsigned call or deployment. A nil To deploys a contract.
type TxRequest struct {
	Nonce uint64
	To    *common.Address
	Value *big.Int
	Data  []byte
	Gas   uint64
}

// Sign builds an EIP-1559 transaction with zero fee caps (the chain charges no
// gas) and signs it. The nonce comes from the caller, so a nonce manager or
// outbox controls ordering, not this type.
func (s *Signer) Sign(req TxRequest) (*types.Transaction, error) {
	value := req.Value
	if value == nil {
		value = new(big.Int)
	}
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   s.chainID,
		Nonce:     req.Nonce,
		GasTipCap: new(big.Int),
		GasFeeCap: new(big.Int),
		Gas:       req.Gas,
		To:        req.To,
		Value:     value,
		Data:      req.Data,
	})
	return types.SignTx(tx, s.signer, s.key)
}
