package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// KEKSize is the key-encryption key length: AES-256.
const KEKSize = 32

// formatV1 prefixes every ciphertext, so the format (or the KEK) can be
// changed later without guessing what an old row holds.
const formatV1 byte = 1

// ErrDecrypt means a stored key could not be decrypted with this KEK.
var ErrDecrypt = errors.New("keystore: cannot decrypt private key (wrong KEY_ENCRYPTION_KEY, or the row was altered)")

// ParseKEK decodes KEY_ENCRYPTION_KEY: 32 bytes, base64 (standard or URL, padding optional).
// Generate one with: openssl rand -base64 32
func ParseKEK(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("KEY_ENCRYPTION_KEY is not set (generate one with: openssl rand -base64 32)")
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		if b, err := enc.DecodeString(s); err == nil {
			if len(b) != KEKSize {
				return nil, fmt.Errorf("KEY_ENCRYPTION_KEY must decode to %d bytes, got %d", KEKSize, len(b))
			}
			return b, nil
		}
	}
	return nil, errors.New("KEY_ENCRYPTION_KEY is not valid base64")
}

// sealer encrypts private keys with AES-256-GCM. The associated data binds
// each ciphertext to its address, so swapping ciphertexts between rows fails.
type sealer struct{ aead cipher.AEAD }

func newSealer(kek []byte) (*sealer, error) {
	if len(kek) != KEKSize {
		return nil, fmt.Errorf("key-encryption key must be %d bytes, got %d", KEKSize, len(kek))
	}
	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &sealer{aead: aead}, nil
}

// seal returns formatV1 || 12-byte random nonce || ciphertext+tag.
func (s *sealer) seal(plaintext []byte, address string) ([]byte, error) {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	out := make([]byte, 0, 1+len(nonce)+len(plaintext)+s.aead.Overhead())
	out = append(out, formatV1)
	out = append(out, nonce...)
	return s.aead.Seal(out, nonce, plaintext, aad(address)), nil
}

func (s *sealer) open(sealed []byte, address string) ([]byte, error) {
	ns := s.aead.NonceSize()
	if len(sealed) < 1+ns+s.aead.Overhead() || sealed[0] != formatV1 {
		return nil, ErrDecrypt
	}
	pt, err := s.aead.Open(nil, sealed[1:1+ns], sealed[1+ns:], aad(address))
	if err != nil {
		return nil, ErrDecrypt
	}
	return pt, nil
}

func aad(address string) []byte { return []byte("signing_keys:" + strings.ToLower(address)) }
