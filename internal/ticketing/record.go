package ticketing

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// A ticket's record is the off-chain data its on-chain recordHash commits to.
// Its canonical form is compact JSON with keys in sorted order (what
// encoding/json does for a map), amounts as strings and times in UTC
// seconds. Postgres stores it as jsonb, which reorders keys and adds spaces,
// so the hash is always taken over the canonical form, never jsonb's text.

func newRecord(eventID string, ticketNumber int64, faceValue string, issuedAt time.Time) []byte {
	b, err := json.Marshal(map[string]any{
		"event_id":      eventID,
		"ticket_number": ticketNumber,
		"face_value":    faceValue,
		"issued_at":     issuedAt.UTC().Truncate(time.Second).Format(time.RFC3339),
	})
	if err != nil {
		panic(err) // a map of strings and an int always marshals
	}
	return b
}

// Canonical re-encodes a record in canonical form, whatever its spacing or
// key order. Numbers are kept exactly as written.
func Canonical(record []byte) ([]byte, error) {
	d := json.NewDecoder(bytes.NewReader(record))
	d.UseNumber()
	var m map[string]any
	if err := d.Decode(&m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// RecordHash is the keccak256 of a record's canonical form: the value
// mintTicket stores on-chain.
func RecordHash(record []byte) (common.Hash, error) {
	c, err := Canonical(record)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(c), nil
}
