package db_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/smlobo/kaleido-challenge/internal/db"
	"github.com/smlobo/kaleido-challenge/internal/dbtest"
)

func TestMigrateIsIdempotent(t *testing.T) {
	pool := dbtest.Pool(t) // already migrated once
	applied, err := db.Migrate(context.Background(), pool)
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 0 {
		t.Fatalf("second Migrate applied %v, want nothing", applied)
	}
}

// fakeHash repeats one hex digit into a value the eth_hash domain accepts
// (0x + 64 lowercase hex). Different digits give different hashes, since
// tx_hash is unique.
func fakeHash(digit string) string { return "0x" + strings.Repeat(digit, 64) }

// Postgres SQLSTATE codes for the two kinds of constraint tested here.
const (
	checkViolation  = "23514"
	uniqueViolation = "23505"
)

// Each insert breaks exactly one rule. Checking the SQLSTATE code and the
// constraint name proves that rule rejected it, not a typo or another rule.
func TestSchemaRejectsBadValues(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	for name, tc := range map[string]struct {
		sql        string
		code       string
		constraint string
	}{
		"mixed-case address": {
			sql: `INSERT INTO signing_keys (address, purpose, encrypted_private_key)
				VALUES ('0xAbCdEf0000000000000000000000000000000000', 'wallet', '\x00')`,
			code: checkViolation, constraint: "eth_address_check",
		},
		"unknown purpose": {
			sql: `INSERT INTO signing_keys (address, purpose, encrypted_private_key)
				VALUES ('0x00000000000000000000000000000000000000aa', 'admin', '\x00')`,
			code: checkViolation, constraint: "signing_keys_purpose_check",
		},
		"second deployer": {
			sql: `INSERT INTO signing_keys (address, purpose, encrypted_private_key) VALUES
				('0x00000000000000000000000000000000000000a1', 'deployer', '\x00'),
				('0x00000000000000000000000000000000000000a2', 'deployer', '\x00')`,
			code: uniqueViolation, constraint: "signing_keys_one_deployer",
		},
		"negative nonce": {
			sql: `INSERT INTO signing_keys (address, purpose, encrypted_private_key, next_nonce)
				VALUES ('0x00000000000000000000000000000000000000bb', 'wallet', '\x00', -1)`,
			code: checkViolation, constraint: "signing_keys_next_nonce_check",
		},
		"second live transaction at a nonce": {
			sql: fmt.Sprintf(`INSERT INTO transactions (from_address, nonce, calldata, raw_tx, tx_hash) VALUES
				('0x00000000000000000000000000000000000000cc', 7, '\x', '\x', '%s'),
				('0x00000000000000000000000000000000000000cc', 7, '\x', '\x', '%s')`, fakeHash("a"), fakeHash("b")),
			code: uniqueViolation, constraint: "transactions_live_nonce",
		},
	} {
		_, err := pool.Exec(ctx, tc.sql)
		var pgErr *pgconn.PgError
		switch {
		case err == nil:
			t.Errorf("%s: accepted", name)
		case !errors.As(err, &pgErr):
			t.Errorf("%s: %v, want a Postgres error", name, err)
		case pgErr.Code != tc.code || pgErr.ConstraintName != tc.constraint:
			t.Errorf("%s: got %s on %q (%s), want %s on %q",
				name, pgErr.Code, pgErr.ConstraintName, pgErr.Message, tc.code, tc.constraint)
		}
	}
}

// A failed transaction gives up its nonce: the outbox's filler takes it.
func TestFailedTransactionFreesItsNonce(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	insert := func(hash, status string) error {
		_, err := pool.Exec(ctx, `INSERT INTO transactions (from_address, nonce, calldata, raw_tx, tx_hash, status)
			VALUES ('0x00000000000000000000000000000000000000dd', 3, '\x', '\x', $1, $2)`, fakeHash(hash), status)
		return err
	}
	if err := insert("a", "failed"); err != nil {
		t.Fatal(err)
	}
	if err := insert("b", "signed"); err != nil {
		t.Fatalf("filler at a failed row's nonce: %v", err)
	}
	if err := insert("c", "signed"); err == nil {
		t.Fatal("second live transaction at the same nonce was accepted")
	}
}
