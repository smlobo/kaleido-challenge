// Command ticketing is the event ticketing service.
//
// Configuration comes only from the environment:
//
//	KEY_ENCRYPTION_KEY required; 32 bytes base64 (openssl rand -base64 32)
//	DATABASE_URL       default postgres://ticketing:ticketing@localhost:5432/ticketing?sslmode=disable
//	ETH_RPC_URL        default http://localhost:8545
//	CHAIN_ID           default 1337
//	HTTP_PORT          default 8080
//	PUBLIC_BASE_URL    default http://localhost:$HTTP_PORT
//	ADMIN_TOKEN        optional; allows POST /organizers (at least 4 characters)
//
// Tracing is on when OTEL_EXPORTER_OTLP_ENDPOINT is set, e.g.
// http://localhost:4318 (OTLP over HTTP); OTEL_SERVICE_NAME names the service
// (default ticketing). See internal/telemetry.
//
// On first start it applies the migrations and creates the deployer key.
// It stops cleanly on Ctrl-C or SIGTERM.
//
// "ticketing healthcheck" asks the running service's /readyz and exits 0 if
// it is ready, 1 if not. The container image has no shell or curl, so its
// HEALTHCHECK runs the binary itself.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smlobo/kaleido-challenge/internal/app"
	"github.com/smlobo/kaleido-challenge/internal/telemetry"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}
	os.Exit(serve())
}

// serve runs the service and returns the exit code. It is separate from main
// so its deferred calls (flushing the last spans) run before the exit.
func serve() int {
	log := slog.Default()
	cfg, err := app.FromEnv(os.Getenv)
	if err != nil {
		log.Error(err.Error())
		return 2
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownTracing, err := telemetry.Setup(ctx)
	if err != nil {
		log.Error("tracing setup failed", "err", err)
		return 2
	}
	defer func() {
		flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracing(flushCtx); err != nil {
			log.Warn("flushing traces", "err", err)
		}
	}()
	if telemetry.Enabled() {
		log.Info("tracing on", "endpoint", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	}

	a, err := app.Start(ctx, cfg, log)
	if err != nil {
		log.Error("startup failed", "err", err)
		return 1
	}
	if err := a.Serve(ctx); err != nil {
		log.Error("stopped with an error", "err", err)
		return 1
	}
	log.Info("stopped")
	return 0
}

// healthcheck returns 0 if GET /readyz on this host's HTTP_PORT answers 200.
func healthcheck() int {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = fmt.Sprint(app.DefaultHTTPPort)
	}
	c := http.Client{Timeout: 3 * time.Second}
	resp, err := c.Get("http://127.0.0.1:" + port + "/readyz")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "readyz:", resp.Status)
		return 1
	}
	return 0
}
