// Command m4demo is the milestone 4 smoke test: everything goes through the
// outbox, and its submitter and receipt loops finish the work.
//
// It deploys a TicketNFT signed by the deployer, then commits five mints and a
// duplicate mint (which reverts) from a new box office key in one database
// transaction, without sending them. The loops send and confirm them.
//
//	-crash   exit right after committing the mints, as if the instance died.
//	         The next run's loops find the unsent rows and finish them.
//
//	DATABASE_URL       default postgres://ticketing:ticketing@localhost:5432/ticketing?sslmode=disable
//	KEY_ENCRYPTION_KEY required; 32 bytes base64 (openssl rand -base64 32)
//	RPC_URL            default http://localhost:8545
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
	"github.com/smlobo/kaleido-challenge/internal/db"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
	"github.com/smlobo/kaleido-challenge/internal/outbox"
)

func main() {
	crash := flag.Bool("crash", false, "exit right after committing the mints, without sending them")
	flag.Parse()
	if err := run(*crash); err != nil {
		slog.Error("m4demo failed", "err", err)
		os.Exit(1)
	}
}

func run(crash bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	kek, err := keystore.ParseKEK(os.Getenv("KEY_ENCRYPTION_KEY"))
	if err != nil {
		return err
	}
	pool, err := db.Open(ctx, envOr("DATABASE_URL", "postgres://ticketing:ticketing@localhost:5432/ticketing?sslmode=disable"))
	if err != nil {
		return err
	}
	defer pool.Close()
	if _, err := db.Migrate(ctx, pool); err != nil {
		return err
	}
	client, err := chain.Dial(ctx, envOr("RPC_URL", "http://localhost:8545"), chain.ChainID1337)
	if err != nil {
		return err
	}
	defer client.Close()
	ks, err := keystore.New(kek, client.ChainID())
	if err != nil {
		return err
	}

	// Rows a previous run left unfinished, e.g. after -crash.
	resumed, err := unfinished(ctx, pool)
	if err != nil {
		return err
	}
	if len(resumed) > 0 {
		slog.Info("resuming unfinished transactions from an earlier run", "count", len(resumed))
	}

	ob := outbox.New(pool, ks, client, outbox.Config{})
	loopCtx, stopLoops := context.WithCancel(ctx)
	loopsDone := make(chan struct{})
	go func() { ob.Run(loopCtx); close(loopsDone) }()
	defer func() { stopLoops(); <-loopsDone }()

	deployer, err := ks.EnsureDeployer(ctx, pool)
	if err != nil {
		return err
	}
	boxOffice, err := ks.Create(ctx, pool, keystore.BoxOffice)
	if err != nil {
		return err
	}
	gate, err := ks.Create(ctx, pool, keystore.Gate)
	if err != nil {
		return err
	}

	// Deploy, using the fast path: Send right after the commit.
	nft := ticketnft.NewTicketNFT()
	code := append(common.FromHex(ticketnft.TicketNFTMetaData.Bin), nft.PackConstructor("Outbox Demo Fest", "OBX", 100, 11_000,
		uint64(time.Now().Add(30*24*time.Hour).Unix()), "http://localhost:8080/tickets/", boxOffice.Address, gate.Address)...)
	gas, err := client.EstimateGas(ctx, deployer.Address, nil, nil, code)
	if err != nil {
		return err
	}
	deploy, err := enqueue(ctx, pool, ob, outbox.Request{KeyID: deployer.ID, Data: code, Gas: gas})
	if err != nil {
		return err
	}
	if err := ob.Send(ctx, deploy[0].ID); err != nil {
		return err
	}
	deployed, err := ob.Wait(ctx, deploy[0].ID)
	if err != nil {
		return err
	}
	if deployed.Status != outbox.Mined {
		return fmt.Errorf("deploy %s: %s", deployed.Status, deployed.LastError)
	}
	addr := *deployed.ContractAddress
	slog.Info("deployed TicketNFT", "address", addr, "block", *deployed.BlockNumber)

	// With -crash, stop this run's loops before committing the mints, so
	// nothing here can send them: only the next run's loops can.
	if crash {
		stopLoops()
		<-loopsDone
	}

	// Five mints and a duplicate of token 3, committed together and not sent.
	// The duplicate uses a fixed gas limit, since estimating would catch the revert.
	holder := common.HexToAddress("0x000000000000000000000000000000000000dead")
	var reqs []outbox.Request
	for _, id := range []int64{1, 2, 3, 4, 5, 3} {
		reqs = append(reqs, outbox.Request{KeyID: boxOffice.ID, To: &addr, Gas: 200_000,
			Data: nft.PackMintTicket(holder, big.NewInt(id), big.NewInt(100_000_000), common.Hash{byte(id)})})
	}
	mints, err := enqueue(ctx, pool, ob, reqs...)
	if err != nil {
		return err
	}
	slog.Info("committed 6 mints; the submitter will send them", "boxOffice", boxOffice.Address)
	if crash {
		slog.Warn("-crash: exiting before anything is sent; run again to see them finished")
		os.Exit(3)
	}

	ids := append([]string{deploy[0].ID}, resumed...)
	for _, m := range mints {
		ids = append(ids, m.ID)
	}
	fmt.Printf("\n%-6s %-44s %-11s %-9s %-6s %s\n", "nonce", "from", "call", "status", "block", "revert")
	for _, id := range ids {
		op, err := ob.Wait(ctx, id)
		if err != nil {
			return err
		}
		fmt.Printf("%-6d %-44s %-11s %-9s %-6s %s\n", op.Nonce, op.From.Hex(), callName(nft, op),
			op.Status, block(op), revert(nft, op))
	}
	return nil
}

// enqueue signs and commits the requests in one database transaction.
func enqueue(ctx context.Context, pool *pgxpool.Pool, ob *outbox.Outbox, reqs ...outbox.Request) ([]outbox.Op, error) {
	var ops []outbox.Op
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		for _, r := range reqs {
			op, err := ob.Enqueue(ctx, tx, r)
			if err != nil {
				return err
			}
			ops = append(ops, op)
		}
		return nil
	})
	return ops, err
}

func unfinished(ctx context.Context, pool *pgxpool.Pool) ([]string, error) {
	rows, err := pool.Query(ctx, `SELECT id FROM transactions WHERE status IN ('signed', 'submitted') ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func callName(nft *ticketnft.TicketNFT, op outbox.Op) string {
	switch {
	case op.To == nil:
		return "deploy"
	case len(op.Calldata) == 0:
		return "filler"
	}
	if m, err := ticketnft.TicketNFTMetaData.ParseABI(); err == nil {
		if method, err := m.MethodById(op.Calldata); err == nil {
			return method.Name
		}
	}
	return "call"
}

func block(op outbox.Op) string {
	if op.BlockNumber == nil {
		return "-"
	}
	return fmt.Sprint(*op.BlockNumber)
}

func revert(nft *ticketnft.TicketNFT, op outbox.Op) string {
	if op.LastError != "" {
		return op.LastError
	}
	if op.RevertReason == "" {
		return ""
	}
	e, err := nft.UnpackError(common.FromHex(op.RevertReason))
	if err != nil {
		return op.RevertReason
	}
	name := strings.TrimPrefix(fmt.Sprintf("%T", e), "*ticketnft.TicketNFT")
	return name + strings.TrimPrefix(fmt.Sprintf("%+v", e), "&")
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
