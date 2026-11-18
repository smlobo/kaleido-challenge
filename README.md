# kaleido-challenge
Event ticketing service on a local Besu QBFT chain (ERC-721 tickets).

## Quick start

Start the chain (chain ID 1337, 2 s blocks, free gas) and Postgres (port 5432):

```bash
docker compose -f docker/docker-compose.yaml up -d --wait
```

Milestone 1: Deploy a plain ERC-721 with a transaction the Go code signs itself and 
sends with `eth_sendRawTransaction`, mint a token and read its owner back:

```bash
go run ./cmd/m1deploy
```

Milestone 2: the same deploy and mint, but the deployer key is stored encrypted in
Postgres and each nonce is reserved in the database transaction that saves the
signed transaction. Run it twice to see the nonces continue:

```bash
export KEY_ENCRYPTION_KEY=$(openssl rand -base64 32)
```

```bash
go run ./cmd/m2demo
```

`KEY_ENCRYPTION_KEY` only ever comes from the environment, never a file. Keep using
the same value for as long as the database lives: keys encrypted under one value
can't be decrypted with another.

Milestone 3: the TicketNFT contract, one per event. The box office mints, holders
list at no more than the resale cap, anyone buys a listed ticket, and the gate burns
it at check-in. Standard ERC-721 transfers and approvals are disabled, so resale is
the only way a ticket moves. Run the contract tests (needs Foundry):

```bash
cd contracts && forge test
```

Milestone 4: the transaction outbox. Every transaction is signed and saved to
Postgres before it is sent. A submitter loop sends anything not yet sent (after a
crash, say) and resends anything stuck, and a receipt loop records each result,
with the revert reason for a failed call. Any number of instances can run the
loops: `FOR UPDATE SKIP LOCKED` hands each row to one worker at a time. The demo
deploys a TicketNFT and commits six mints, one of them a duplicate that reverts,
without sending them; the loops do the rest. With `-crash` it exits before
sending, and the next run finishes the job:

```bash
go run ./cmd/m4demo -crash
```

```bash
go run ./cmd/m4demo
```

Milestone 5: the service itself. On start it waits up to 30 s for Postgres and
Besu, checks the chain ID, applies migrations, creates the deployer key on first
start (or loads it), and decrypts it so a wrong `KEY_ENCRYPTION_KEY` fails at
once. Then it runs the outbox loops and serves `GET /healthz` and `GET /readyz`
on `HTTP_PORT` (8080). Ctrl-C or SIGTERM stops it cleanly. Configuration is
environment only: `KEY_ENCRYPTION_KEY` (required), `DATABASE_URL`,
`ETH_RPC_URL`, `CHAIN_ID`, `HTTP_PORT`, `PUBLIC_BASE_URL`; the defaults match
docker compose.

```bash
go run ./cmd/ticketing
```

```bash
curl -s localhost:8080/readyz
```

Run the Go tests (chain tests skip if Besu is down, database tests if Postgres is
down; `-count=1` skips Go's test cache):

```bash
go test -count=1 ./...
```

More commands for inspecting the chain (`cast`, Besu logs, scripts): [docs/dev-tools.md](docs/dev-tools.md).

## Layout

| Path | What |
|---|---|
| `contracts/src` | Solidity sources (OpenZeppelin 5.6.1, solc 0.8.38; settings in `contracts/foundry.toml`) |
| `contracts/test` | Forge tests for TicketNFT |
| `contracts/build.sh` | Compiles with Foundry (`forge build`) and regenerates the Go bindings |
| `contracts/build` | Committed ABI and bytecode |
| `internal/contracts` | Committed abigen v2 bindings, so the service builds without solc |
| `internal/chain` | Signing, raw transaction send, receipts (with Besu's revert reason) |
| `internal/db` | Postgres pool and embedded schema migrations (`internal/db/migrations`) |
| `internal/keystore` | AES-256-GCM encrypted signing keys and nonce reservation |
| `internal/outbox` | Transaction outbox: enqueue, send, submitter and receipt loops |
| `internal/app` | The service: configuration, startup, health checks, shutdown |
| `internal/dbtest` | Test helper: a fresh migrated schema per test |
| `cmd/ticketing` | The service's entry point |
| `cmd/m1deploy` | Milestone 1 smoke test |
| `cmd/m2demo` | Milestone 2 smoke test |
| `cmd/m4demo` | Milestone 4 smoke test |
