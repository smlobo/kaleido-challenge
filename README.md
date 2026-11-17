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
| `internal/chain` | Signing, raw transaction send, receipt wait |
| `internal/db` | Postgres pool and embedded schema migrations (`internal/db/migrations`) |
| `internal/keystore` | AES-256-GCM encrypted signing keys and nonce reservation |
| `internal/dbtest` | Test helper: a fresh migrated schema per test |
| `cmd/m1deploy` | Milestone 1 smoke test |
| `cmd/m2demo` | Milestone 2 smoke test |
