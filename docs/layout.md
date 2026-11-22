# Code layout

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
| `internal/api` | HTTP layer: `openapi.yaml`, generated server code, handlers, bearer tokens, Swagger UI |
| `internal/ticketing` | Business rules: organizers, users, events, tickets, check-in, resale, and the outbox's final-status hook |
| `docs/schema.sql` | The schema all migrations produce, in one file, for reading; regenerate with `scripts/schema-dump.sh` |
| `internal/dbtest` | Test helper: a fresh migrated schema per test |
| `internal/telemetry` | OpenTelemetry setup, SQL spans, JSON-RPC spans with trace headers, trace context for the outbox |
| `cmd/ticketing` | The service's entry point |
| `frontend` | Demo page: React + TypeScript on Vite, proxying `/api` to the service |
| `Dockerfile` | The service image: Go build stage, distroless non-root runtime |
| `docker/docker-compose.yaml` | Besu, Postgres and the service |
| `cmd/m1deploy` | Milestone 1 smoke test |
| `cmd/m2demo` | Milestone 2 smoke test |
| `cmd/m4demo` | Milestone 4 smoke test |
| `deploy/k8s` | Kubernetes manifests (microk8s) for Postgres, Besu, Jaeger, two service pods and the frontend; see its [README](../deploy/k8s/README.md) |
