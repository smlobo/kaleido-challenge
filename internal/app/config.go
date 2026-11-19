// Package app is the ticketing service: it connects to Postgres and Besu,
// prepares the database and the deployer key, runs the outbox loops and
// serves HTTP. cmd/ticketing is a thin wrapper around it.
package app

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/keystore"
)

// Config is everything the service needs. It comes only from the
// environment (see FromEnv); keys never come from files.
type Config struct {
	DatabaseURL   string
	RPCURL        string
	ChainID       *big.Int
	KEK           []byte // KEY_ENCRYPTION_KEY, decoded; never logged
	HTTPAddr      string // host:port to listen on
	PublicBaseURL string // how clients reach the API; used in ticket URIs
	AdminToken    string // ADMIN_TOKEN, optional; allows POST /organizers; never logged

	// StartupTimeout is how long to keep retrying Postgres and Besu at
	// startup, so the service can start alongside them.
	StartupTimeout time.Duration
	// ShutdownTimeout is how long in-flight HTTP requests get to finish.
	ShutdownTimeout time.Duration
}

// Defaults match docker/docker-compose.yaml.
const (
	DefaultDatabaseURL = "postgres://ticketing:ticketing@localhost:5432/ticketing?sslmode=disable"
	DefaultRPCURL      = "http://localhost:8545"
	DefaultChainID     = 1337
	DefaultHTTPPort    = 8080
)

// FromEnv reads the configuration. Every problem is reported at once, so a
// misconfigured deployment is fixed in one go. getenv is os.Getenv in
// production.
func FromEnv(getenv func(string) string) (Config, error) {
	get := func(k, def string) string {
		if v := strings.TrimSpace(getenv(k)); v != "" {
			return v
		}
		return def
	}
	var errs []error
	cfg := Config{
		DatabaseURL:     get("DATABASE_URL", DefaultDatabaseURL),
		RPCURL:          get("ETH_RPC_URL", DefaultRPCURL),
		StartupTimeout:  30 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}

	if _, err := pgxpool.ParseConfig(cfg.DatabaseURL); err != nil {
		errs = append(errs, errors.New("DATABASE_URL is not a valid Postgres URL")) // the URL may hold a password
	}
	if u, err := url.Parse(cfg.RPCURL); err != nil || (u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "ws" && u.Scheme != "wss") || u.Host == "" {
		errs = append(errs, fmt.Errorf("ETH_RPC_URL must be an http(s) or ws(s) URL, got %q", cfg.RPCURL))
	}

	chainID, ok := new(big.Int).SetString(get("CHAIN_ID", strconv.Itoa(DefaultChainID)), 10)
	if !ok || chainID.Sign() <= 0 {
		errs = append(errs, fmt.Errorf("CHAIN_ID must be a positive integer, got %q", getenv("CHAIN_ID")))
	}
	cfg.ChainID = chainID

	kek, err := keystore.ParseKEK(getenv("KEY_ENCRYPTION_KEY"))
	if err != nil {
		errs = append(errs, err) // the messages never include the value
	}
	cfg.KEK = kek

	port, err := strconv.Atoi(get("HTTP_PORT", strconv.Itoa(DefaultHTTPPort)))
	if err != nil || port < 1 || port > 65535 {
		errs = append(errs, fmt.Errorf("HTTP_PORT must be 1-65535, got %q", getenv("HTTP_PORT")))
	}
	cfg.HTTPAddr = net.JoinHostPort("", strconv.Itoa(port))

	cfg.PublicBaseURL = strings.TrimRight(get("PUBLIC_BASE_URL", fmt.Sprintf("http://localhost:%d", port)), "/")
	if u, err := url.Parse(cfg.PublicBaseURL); err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		errs = append(errs, fmt.Errorf("PUBLIC_BASE_URL must be an http(s) URL, got %q", cfg.PublicBaseURL))
	}

	// Optional. Without it nobody can create organizers.
	cfg.AdminToken = strings.TrimSpace(getenv("ADMIN_TOKEN"))
	if cfg.AdminToken != "" && len(cfg.AdminToken) < 4 {
		errs = append(errs, errors.New("ADMIN_TOKEN must be at least 4 characters"))
	}

	if err := errors.Join(errs...); err != nil {
		return Config{}, fmt.Errorf("configuration:\n%w", err)
	}
	return cfg, nil
}
