package ticketing

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/smlobo/kaleido-challenge/internal/keystore"
)

// OrganizerInfo is an organizer and its keys' addresses.
type OrganizerInfo struct {
	ID, Name                string
	BoxOffice, Gate, Payout string // addresses
	Token                   string // set only when just created
}

// UserInfo is a user and their wallet address.
type UserInfo struct {
	ID, Subject, Wallet string
	Token               string // set only when just created
}

// Authenticate resolves a bearer token. An empty token is Anonymous.
func (s *Service) Authenticate(ctx context.Context, token string) (Principal, error) {
	if token == "" {
		return Principal{}, nil
	}
	h := sha256.Sum256([]byte(token))
	if s.adminHash != nil && subtle.ConstantTimeCompare(h[:], s.adminHash) == 1 {
		return Principal{Role: Admin}, nil
	}
	var org, user *string
	err := s.pool.QueryRow(ctx, `SELECT organizer_id, user_id FROM api_tokens WHERE token_hash = $1`, h[:]).Scan(&org, &user)
	switch {
	case isNoRows(err):
		return Principal{}, errBadToken
	case err != nil:
		return Principal{}, err
	case org != nil:
		return Principal{Role: Organizer, ID: *org}, nil
	default:
		return Principal{Role: User, ID: *user}, nil
	}
}

// newToken returns a random API token and the hash that is stored.
func newToken() (string, []byte) {
	b := make([]byte, 32)
	_, _ = rand.Read(b) // never fails (crypto/rand)
	t := "tk_" + base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(t))
	return t, h[:]
}

// CreateOrganizer creates an organizer with new box office, gate and payout
// keys, and its API token. Only the platform admin may.
func (s *Service) CreateOrganizer(ctx context.Context, p Principal, name string) (OrganizerInfo, error) {
	switch {
	case p.Role == Anonymous:
		return OrganizerInfo{}, errNoToken
	case s.adminHash == nil:
		return OrganizerInfo{}, errAdminUnset
	case p.Role != Admin:
		return OrganizerInfo{}, errNoAdmin
	}
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 200 {
		return OrganizerInfo{}, fail(Invalid, "invalid_name", "Invalid name", "The name must be 1 to 200 characters.")
	}
	out := OrganizerInfo{Name: name}
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var keys [3]keystore.Key
		for i, purpose := range []keystore.Purpose{keystore.BoxOffice, keystore.Gate, keystore.Payout} {
			k, err := s.keys.Create(ctx, tx, purpose)
			if err != nil {
				return err
			}
			keys[i] = k
		}
		if err := tx.QueryRow(ctx,
			`INSERT INTO organizers (name, box_office_key_id, gate_key_id, payout_key_id)
			 VALUES ($1, $2, $3, $4) RETURNING id`,
			name, keys[0].ID, keys[1].ID, keys[2].ID).Scan(&out.ID); err != nil {
			return err
		}
		token, hash := newToken()
		if _, err := tx.Exec(ctx, `INSERT INTO api_tokens (token_hash, organizer_id) VALUES ($1, $2)`, hash, out.ID); err != nil {
			return err
		}
		out.Token = token
		out.BoxOffice, out.Gate, out.Payout = hexAddr(keys[0].Address), hexAddr(keys[1].Address), hexAddr(keys[2].Address)
		return nil
	})
	return out, err
}

// CreateUser creates a user, their custodial wallet key and their API token.
// Anyone may.
func (s *Service) CreateUser(ctx context.Context, subject string) (UserInfo, error) {
	subject = strings.TrimSpace(subject)
	if subject == "" || len(subject) > 200 {
		return UserInfo{}, fail(Invalid, "invalid_subject", "Invalid subject", "The subject must be 1 to 200 characters.")
	}
	out := UserInfo{Subject: subject}
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		wallet, err := s.keys.Create(ctx, tx, keystore.Wallet)
		if err != nil {
			return err
		}
		if err := tx.QueryRow(ctx,
			`INSERT INTO users (subject, wallet_key_id) VALUES ($1, $2) RETURNING id`,
			subject, wallet.ID).Scan(&out.ID); err != nil {
			return err
		}
		token, hash := newToken()
		if _, err := tx.Exec(ctx, `INSERT INTO api_tokens (token_hash, user_id) VALUES ($1, $2)`, hash, out.ID); err != nil {
			return err
		}
		out.Token, out.Wallet = token, hexAddr(wallet.Address)
		return nil
	})
	if isUnique(err, "users_subject_key") {
		return UserInfo{}, fail(Conflict, "subject_taken", "Subject taken", "A user with that subject already exists.")
	}
	return out, err
}
