# kaleido-challenge
Event ticketing service on a local Besu QBFT chain (ERC-721 tickets, ERC-20 payments).

## Quick start (milestone 1)

Start the chain (chain ID 1337, 2 s blocks, free gas):

```bash
docker compose -f docker/docker-compose.yaml up -d --wait
```

Deploy a plain ERC-721 with a transaction the Go code signs itself and sends with
`eth_sendRawTransaction`, mint a token and read its owner back:

```bash
go run ./cmd/m1deploy
```

Run the tests (the chain tests skip if Besu is not running):

```bash
go test ./...
```

## Layout

| Path | What |
|---|---|
| `contracts/src` | Solidity sources (OpenZeppelin 5.6.1, solc 0.8.38; settings in `contracts/foundry.toml`) |
| `contracts/build.sh` | Compiles with Foundry (`forge build`) and regenerates the Go bindings |
| `contracts/build` | Committed ABI and bytecode |
| `internal/contracts` | Committed abigen v2 bindings, so the service builds without solc |
| `internal/chain` | Signing, raw transaction send, receipt wait |
| `cmd/m1deploy` | Milestone 1 smoke test |
