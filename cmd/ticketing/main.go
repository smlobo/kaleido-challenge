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
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}
	log := slog.Default()
	cfg, err := app.FromEnv(os.Getenv)
	if err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a, err := app.Start(ctx, cfg, log)
	if err != nil {
		log.Error("startup failed", "err", err)
		os.Exit(1)
	}
	if err := a.Serve(ctx); err != nil {
		log.Error("stopped with an error", "err", err)
		os.Exit(1)
	}
	log.Info("stopped")
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
