package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/api"
	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/db"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
	"github.com/smlobo/kaleido-challenge/internal/outbox"
	"github.com/smlobo/kaleido-challenge/internal/ticketing"
)

// App is a started service: connected, migrated, with a deployer key and a
// listening HTTP socket. Serve runs it; Close releases it without serving.
type App struct {
	cfg      Config
	log      *slog.Logger
	pool     *pgxpool.Pool
	chain    *chain.Client
	keys     *keystore.Store
	outbox   *outbox.Outbox
	deployer keystore.Key
	svc      *ticketing.Service
	ln       net.Listener
	srv      *http.Server
}

// Start brings the service up, in order:
//
//  1. connect to Postgres and Besu, retrying for StartupTimeout;
//  2. check the node's chain ID;
//  3. apply migrations (one instance at a time, under an advisory lock);
//  4. create the deployer key on first start, or load it;
//  5. decrypt it, so a wrong KEY_ENCRYPTION_KEY fails now, not at the first deploy;
//  6. listen on HTTPAddr.
//
// Nothing runs until Serve.
func Start(ctx context.Context, cfg Config, log *slog.Logger) (_ *App, err error) {
	if log == nil {
		log = slog.Default()
	}
	a := &App{cfg: cfg, log: log}
	defer func() {
		if err != nil {
			a.Close()
		}
	}()

	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_URL: %w", err)
	}
	deadline := time.Now().Add(cfg.StartupTimeout)
	a.pool, err = retry(ctx, log, deadline, "Postgres", func(ctx context.Context) (*pgxpool.Pool, error) {
		return db.OpenConfig(ctx, poolCfg.Copy())
	})
	if err != nil {
		return nil, err
	}
	a.chain, err = retry(ctx, log, deadline, "Besu", func(ctx context.Context) (*chain.Client, error) {
		return chain.Dial(ctx, cfg.RPCURL, nil)
	})
	if err != nil {
		return nil, err
	}
	if got := a.chain.ChainID(); got.Cmp(cfg.ChainID) != 0 {
		return nil, fmt.Errorf("the node at ETH_RPC_URL serves chain %s, but CHAIN_ID is %s", got, cfg.ChainID)
	}

	applied, err := db.Migrate(ctx, a.pool)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if len(applied) > 0 {
		log.Info("applied migrations", "migrations", applied)
	}

	a.keys, err = keystore.New(cfg.KEK, cfg.ChainID)
	if err != nil {
		return nil, err
	}
	a.deployer, err = a.keys.EnsureDeployer(ctx, a.pool)
	if err != nil {
		return nil, fmt.Errorf("deployer key: %w", err)
	}
	if _, err := a.keys.Signer(ctx, a.pool, a.deployer.ID); err != nil {
		if errors.Is(err, keystore.ErrDecrypt) {
			return nil, fmt.Errorf("KEY_ENCRYPTION_KEY does not match the keys in this database "+
				"(the deployer key %s was encrypted with a different one): %w", a.deployer.Address, err)
		}
		return nil, fmt.Errorf("deployer key: %w", err)
	}

	a.outbox = outbox.New(a.pool, a.keys, a.chain, outbox.Config{Logger: log, OnFinal: ticketing.OnFinal})
	a.svc = ticketing.New(ticketing.Config{
		Pool: a.pool, Keys: a.keys, Outbox: a.outbox, Node: a.chain, Deployer: a.deployer,
		PublicBaseURL: cfg.PublicBaseURL, AdminToken: cfg.AdminToken, Logger: log,
	})

	a.ln, err = net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", cfg.HTTPAddr, err)
	}
	a.srv = &http.Server{Handler: a.routes(), ReadHeaderTimeout: 10 * time.Second}

	if cfg.AdminToken == "" {
		log.Warn("ADMIN_TOKEN is not set: POST /organizers is disabled, so no new organizers can be created")
	}
	log.Info("ready", "chainID", cfg.ChainID, "deployer", a.deployer.Address, "http", a.Addr(),
		"admin", cfg.AdminToken != "")
	return a, nil
}

// Serve runs the outbox loops and the HTTP server until ctx ends, then shuts
// down: HTTP first (in-flight requests get ShutdownTimeout), then the loops,
// then the connections.
func (a *App) Serve(ctx context.Context) error {
	defer a.Close()

	// The loops get their own context so they stop after HTTP, not with it:
	// a request still finishing can rely on them being there.
	loopCtx, stopLoops := context.WithCancel(context.WithoutCancel(ctx))
	loopsDone := make(chan struct{})
	go func() { a.outbox.Run(loopCtx); close(loopsDone) }()
	defer func() { stopLoops(); <-loopsDone }()

	httpErr := make(chan error, 1)
	go func() { httpErr <- a.srv.Serve(a.ln) }()

	select {
	case err := <-httpErr:
		return fmt.Errorf("http: %w", err)
	case <-ctx.Done():
	}
	a.log.Info("shutting down")
	shutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), a.cfg.ShutdownTimeout)
	defer cancel()
	if err := a.srv.Shutdown(shutCtx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}
	return nil
}

// Close releases the connections. Serve calls it; call it yourself only if
// Start succeeded and Serve is never called.
func (a *App) Close() {
	if a.ln != nil {
		_ = a.ln.Close() // already closed after Shutdown; harmless
	}
	if a.chain != nil {
		a.chain.Close()
		a.chain = nil
	}
	if a.pool != nil {
		a.pool.Close()
		a.pool = nil
	}
}

// Addr is the address the HTTP server listens on (useful with port 0).
func (a *App) Addr() string { return a.ln.Addr().String() }

// Deployer is the platform's deployer key.
func (a *App) Deployer() keystore.Key { return a.deployer }

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("GET /readyz", readiness(map[string]func(context.Context) error{
		"postgres": func(ctx context.Context) error { return a.pool.Ping(ctx) },
		"besu": func(ctx context.Context) error {
			_, err := a.chain.Eth().BlockNumber(ctx)
			return err
		},
	}))
	api.Register(mux, a.svc, a.log)
	return mux
}

// readiness answers 200 when every check passes within 2 s, else 503 with
// the failing checks' errors.
func readiness(checks map[string]func(context.Context) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		status, results := http.StatusOK, map[string]string{}
		for name, check := range checks {
			if err := check(ctx); err != nil {
				status, results[name] = http.StatusServiceUnavailable, err.Error()
				continue
			}
			results[name] = "ok"
		}
		writeJSON(w, status, results)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// retry calls connect until it succeeds, ctx ends or the deadline passes,
// once a second.
func retry[T any](ctx context.Context, log *slog.Logger, deadline time.Time, what string,
	connect func(context.Context) (T, error)) (T, error) {
	for attempt := 1; ; attempt++ {
		tryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		v, err := connect(tryCtx)
		cancel()
		if err == nil {
			return v, nil
		}
		if ctx.Err() != nil || time.Now().Add(time.Second).After(deadline) {
			var zero T
			return zero, fmt.Errorf("connect to %s: %w", what, err)
		}
		if attempt == 1 {
			log.Info("waiting for "+what, "err", err)
		}
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
		}
	}
}
