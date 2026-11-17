package db_test

import (
	"context"
	"errors"
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
