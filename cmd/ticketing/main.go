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
//
// On first start it applies the migrations and creates the deployer key.
// It stops cleanly on Ctrl-C or SIGTERM.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/smlobo/kaleido-challenge/internal/app"
)

func main() {
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
