// Command m2demo is the milestone 2 smoke test. It runs the same deploy and
// mint as m1deploy, but the deployer key lives encrypted in Postgres and each
// nonce comes from the database, reserved in the same transaction that saves
// the signed transaction (a first cut of the milestone 4 outbox).
//
//	DATABASE_URL       default postgres://ticketing:ticketing@localhost:5432/ticketing?sslmode=disable
//	KEY_ENCRYPTION_KEY required; 32 bytes base64 (openssl rand -base64 32)
//	RPC_URL            default http://localhost:8545
//
// Run it twice: the second run reuses the same deployer and continues its nonces.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/contracts/demonft"
	"github.com/smlobo/kaleido-challenge/internal/db"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
)

func main() {
	if err := run(); err != nil {
		slog.Error("m2demo failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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
	applied, err := db.Migrate(ctx, pool)
	if err != nil {
		return err
	}
	slog.Info("database ready", "migrationsApplied", applied)

	client, err := chain.Dial(ctx, envOr("RPC_URL", "http://localhost:8545"), chain.ChainID1337)
	if err != nil {
		return err
	}
	defer client.Close()

	ks, err := keystore.New(kek, client.ChainID())
	if err != nil {
		return err
	}
	deployer, err := ks.EnsureDeployer(ctx, pool)
	if err != nil {
		return err
	}
	slog.Info("deployer key", "id", deployer.ID, "address", deployer.Address, "nextNonce", deployer.NextNonce)

	if err := checkNonceInSync(ctx, client, deployer); err != nil {
		return err
	}

	nft := demonft.NewDemoNFT()
	deploy, err := sendViaOutbox(ctx, pool, ks, client, deployer, nil,
		append(common.FromHex(demonft.DemoNFTMetaData.Bin), nft.PackConstructor("Demo Ticket", "DEMO", deployer.Address)...))
	if err != nil {
		return fmt.Errorf("deploy: %w", err)
	}
	addr := deploy.ContractAddress
	slog.Info("deployed DemoNFT", "address", addr, "block", deploy.BlockNumber)

	holder := common.HexToAddress("0x000000000000000000000000000000000000c001")
	tokenID := big.NewInt(1)
	mintDemo, err := sendViaOutbox(ctx, pool, ks, client, deployer, &addr, nft.PackSafeMint(holder, tokenID))
	if err != nil {
		return fmt.Errorf("safeMint: %w", err)
	}
	slog.Info("minted DemoNFT", "contract", addr, "holder", holder, "block", mintDemo.BlockNumber)

	owner, err := bind.Call(nft.Instance(client.Eth(), addr), &bind.CallOpts{Context: ctx},
		nft.PackOwnerOf(tokenID), nft.UnpackOwnerOf)
	if err != nil {
		return err
	}
	if owner != holder {
		return fmt.Errorf("ownerOf(1) = %s, want %s", owner, holder)
	}
	slog.Info("verified on-chain", "ownerOf(1)", owner)
	return nil
}

// sendViaOutbox is the write path in miniature:
//  1. one DB transaction: reserve the nonce, sign, insert the signed bytes; commit [db status: signed]
//  2. send the saved bytes with eth_sendRawTransaction [db status: submitted]
//  3. wait for the receipt and record the result [db status: mined or reverted]
//
// Nothing reaches the chain before step 1 commits, so a crash can always be
// recovered by resending what is in the table.
func sendViaOutbox(ctx context.Context, pool *pgxpool.Pool, ks *keystore.Store, c *chain.Client,
	key keystore.Key, to *common.Address, data []byte) (*types.Receipt, error) {

	gas, err := c.EstimateGas(ctx, key.Address, to, nil, data)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	var opID string
	err = pgx.BeginFunc(ctx, pool, func(dbtx pgx.Tx) error {
		nonce, err := keystore.ReserveNonce(ctx, dbtx, key.ID)
		if err != nil {
			return err
		}
		signer, err := ks.Signer(ctx, dbtx, key.ID)
		if err != nil {
			return err
		}
		tx, err = signer.Sign(chain.TxRequest{Nonce: nonce, To: to, Data: data, Gas: gas})
		if err != nil {
			return err
		}
		raw, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		return dbtx.QueryRow(ctx,
			`INSERT INTO transactions (from_address, nonce, to_address, calldata, raw_tx, tx_hash)
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			lower(key.Address), int64(nonce), lowerPtr(to), data, raw, tx.Hash().Hex()).Scan(&opID)
	})
	if err != nil {
		return nil, err
	}
	slog.Info("signed and saved", "op", opID, "nonce", tx.Nonce(), "tx", tx.Hash())

	if _, err := c.SendRaw(ctx, tx); err != nil {
		return nil, err // the row stays 'signed'; TODO: needs retry
	}
	if _, err := pool.Exec(ctx, `UPDATE transactions SET status = 'submitted', attempts = attempts + 1,
		updated_at = now() WHERE id = $1`, opID); err != nil {
		return nil, err
	}

	rcpt, waitErr := c.WaitMined(ctx, tx.Hash(), time.Second)
	if rcpt != nil {
		status := "mined"
		if errors.Is(waitErr, chain.ErrReverted) {
			status = "reverted"
		}
		var contract *string
		if rcpt.ContractAddress != (common.Address{}) {
			s := lower(rcpt.ContractAddress)
			contract = &s
		}
		if _, err := pool.Exec(ctx, `UPDATE transactions SET status = $2, block_number = $3,
			contract_address = $4, updated_at = now() WHERE id = $1`,
			opID, status, rcpt.BlockNumber.Int64(), contract); err != nil {
			return nil, err
		}
	}
	return rcpt, waitErr
}

// checkNonceInSync compares the database counter with the chain. They differ
// if the chain was reset without wiping Postgres (database ahead), or if the
// key was used outside this service (chain ahead).
func checkNonceInSync(ctx context.Context, c *chain.Client, k keystore.Key) error {
	onChain, err := c.PendingNonce(ctx, k.Address)
	if err != nil {
		return err
	}
	switch {
	case k.NextNonce > onChain:
		return fmt.Errorf("database next_nonce %d is ahead of the chain (%d): the chain was probably reset; "+
			"reset both with: docker compose -f docker/docker-compose.yaml down && docker compose -f docker/docker-compose.yaml up -d --wait",
			k.NextNonce, onChain)
	case k.NextNonce < onChain:
		return fmt.Errorf("chain nonce %d is ahead of the database (%d): the key was used outside this service", onChain, k.NextNonce)
	}
	return nil
}

func lower(a common.Address) string { return strings.ToLower(a.Hex()) }

func lowerPtr(a *common.Address) *string {
	if a == nil {
		return nil
	}
	s := lower(*a)
	return &s
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
