package keystore

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"
)

func testKEK(t *testing.T) []byte {
	t.Helper()
	k := make([]byte, KEKSize)
	if _, err := rand.Read(k); err != nil {
		t.Fatal(err)
	}
	return k
}

func TestParseKEK(t *testing.T) {
	good := testKEK(t)
	for _, s := range []string{
		base64.StdEncoding.EncodeToString(good),
		base64.RawStdEncoding.EncodeToString(good),
		base64.URLEncoding.EncodeToString(good),
		"  " + base64.StdEncoding.EncodeToString(good) + "\n",
	} {
		got, err := ParseKEK(s)
		if err != nil || !bytes.Equal(got, good) {
			t.Errorf("ParseKEK(%q) = %x, %v", s, got, err)
		}
	}
	for _, s := range []string{"", "not base64!!", base64.StdEncoding.EncodeToString(good[:16])} {
		if _, err := ParseKEK(s); err == nil {
			t.Errorf("ParseKEK(%q) succeeded, want error", s)
		}
	}
}

func TestSealRoundTrip(t *testing.T) {
	s, err := newSealer(testKEK(t))
	if err != nil {
		t.Fatal(err)
	}
	secret := []byte("0123456789abcdef0123456789abcdef")
	addr := "0x00000000000000000000000000000000000000aa"

	a, _ := s.seal(secret, addr)
	b, _ := s.seal(secret, addr)
	if bytes.Equal(a, b) {
		t.Error("two seals of the same key are identical; nonce not random")
	}
	if bytes.Contains(a, secret) {
		t.Error("ciphertext contains the plaintext")
	}
	got, err := s.open(a, addr)
	if err != nil || !bytes.Equal(got, secret) {
		t.Fatalf("open = %x, %v", got, err)
	}
	// Address comparison is case-insensitive (checksummed vs lowercase).
	if _, err := s.open(a, "0x"+strings.ToUpper(addr[2:])); err != nil {
		t.Errorf("open with checksummed address: %v", err)
	}
}

func TestOpenRejects(t *testing.T) {
	kek := testKEK(t)
	s, _ := newSealer(kek)
	addr := "0x00000000000000000000000000000000000000aa"
	sealed, _ := s.seal([]byte("secret-secret-secret-secret-1234"), addr)

	other, _ := newSealer(testKEK(t))
	if _, err := other.open(sealed, addr); err == nil {
		t.Error("opened with the wrong KEK")
	}
	if _, err := s.open(sealed, "0x00000000000000000000000000000000000000bb"); err == nil {
		t.Error("opened under a different address (ciphertext moved between rows)")
	}
	tampered := bytes.Clone(sealed)
	tampered[len(tampered)-1] ^= 1
	if _, err := s.open(tampered, addr); err == nil {
		t.Error("opened a tampered ciphertext")
	}
	if _, err := s.open(sealed[:10], addr); err == nil {
		t.Error("opened a truncated ciphertext")
	}
	wrongVersion := bytes.Clone(sealed)
	wrongVersion[0] = 9
	if _, err := s.open(wrongVersion, addr); err == nil {
		t.Error("opened an unknown format version")
	}
}
