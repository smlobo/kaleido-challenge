// Package dbtest gives each test its own migrated Postgres schema.
package dbtest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/db"
)

// DefaultURL is the Postgres from docker/docker-compose.yaml.
const DefaultURL = "postgres://ticketing:ticketing@localhost:5432/ticketing?sslmode=disable"

// URL is DATABASE_URL, or DefaultURL.
func URL() string {
	if u := os.Getenv("DATABASE_URL"); u != "" {
		return u
	}
	return DefaultURL
}

// Pool returns a pool whose search_path is a fresh schema with all migrations
// applied. The schema is dropped when the test ends. The test is skipped if
// Postgres is not reachable.
func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	admin, err := db.Open(ctx, URL())
	if err != nil {
		t.Skipf("Postgres not reachable (%v); run: docker compose -f docker/docker-compose.yaml up -d --wait", err)
	}
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	schema := "test_" + hex.EncodeToString(b)
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		admin.Close()
		t.Fatalf("create schema: %v", err)
	}

	cfg, err := pgxpool.ParseConfig(URL())
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	cfg.MaxConns = 64
	pool, err := db.OpenConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
		if _, err := admin.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Logf("drop schema %s: %v", schema, err)
		}
		admin.Close()
	})
	return pool
}
