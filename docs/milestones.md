# Build log and usage

How the service was built, milestone by milestone, with the commands to run each
step. For the design, see the [README](../README.md); for inspecting the chain
and database, [dev-tools.md](dev-tools.md); for the code layout, [layout.md](layout.md).

## Quick start

Start the chain (chain ID 1337, 2 s blocks, free gas) and Postgres (port 5432).
Jaeger starts with them, since Besu sends it traces (UI on port 16686):

```bash
docker compose -f docker/docker-compose.yaml up -d --wait besu postgres
```

## Milestone 1

Deploy a plain ERC-721 with a transaction the Go code signs itself and 
sends with `eth_sendRawTransaction`, mint a token and read its owner back:

```bash
go run ./cmd/m1deploy
```

## Milestone 2

The same deploy and mint, but the deployer key is stored encrypted in
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

## Milestone 3

The TicketNFT contract, one per event. The box office mints, holders
list at no more than the resale cap, anyone buys a listed ticket, and the gate burns
it at check-in, carrying a pass the holder signed (EIP-712), so both the organizer
and the holder authorize the check-in on-chain. Standard ERC-721 transfers and approvals are disabled, so resale is
the only way a ticket moves. Run the contract tests (needs Foundry):

```bash
cd contracts && forge test
```

## Milestone 4

The transaction outbox. Every transaction is signed and saved to
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

## Milestone 5

The service itself. On start it waits up to 30 s for Postgres and
Besu, checks the chain ID, applies migrations, creates the deployer key on first
start (or loads it), and decrypts it so a wrong `KEY_ENCRYPTION_KEY` fails at
once. Then it runs the outbox loops and serves `GET /healthz` and `GET /readyz`
on `HTTP_PORT` (8080). Ctrl-C or SIGTERM stops it cleanly. Configuration is
environment only: `KEY_ENCRYPTION_KEY` (required), `DATABASE_URL`,
`ETH_RPC_URL`, `CHAIN_ID`, `HTTP_PORT`, `PUBLIC_BASE_URL`, and from milestone 6
`ADMIN_TOKEN` (optional); the defaults match docker compose.

```bash
go run ./cmd/ticketing
```

```bash
curl -s localhost:8080/readyz
```

## Milestone 6

Events, tickets and check-in over HTTP. An admin creates an
organizer (with `ADMIN_TOKEN`), anyone creates a user, the organizer creates an
event (its TicketNFT is deployed through the outbox), a user buys tickets (the
organizer's box office mints them to the user's wallet), and the organizer
checks a ticket in with a pass the holder makes at the door (the gate key burns
it, and the contract checks the holder signed the pass). Writes answer `202` with an
operation ID and need an `Idempotency-Key` header; `GET /operations/{id}`
reports the result. The spec is
[internal/api/openapi.yaml](../internal/api/openapi.yaml), with Swagger UI at
`/docs`. Start the service with an admin token (keep the same
`KEY_ENCRYPTION_KEY` exported):

```bash
ADMIN_TOKEN=1234 go run ./cmd/ticketing
```

Exercise all endpoints:

```bash
./scripts/create_organizer.py "Rock Fest Ltd" --token 1234
./scripts/create_user.py alice
./scripts/create_event.py "Rock Fest 2027" --token <org-token>
./scripts/buy_ticket.py <event-id> --token <alice-token>
./scripts/create_pass.py <ticket-id> --token <alice-token> | ./scripts/check_in.py <ticket-id> --token <org-token> --pass -
./scripts/get_operation.py OPERATION_ID --token <respective-token>
curl localhost:8080/events/<event-id>
curl localhost:8080/tickets/<ticket-id>
curl localhost:8080/metadata/<event-id>/<token-id>
curl localhost:8080/events/<event-id>/listings
```

## Milestone 7

Resale. A holder lists a ticket at no more than the event's resale
cap (face value x `max_resale_bps` / 10000), can cancel the listing, and another
user buys it. Each step is a transaction from that user's own wallet, and the
contract checks it: only the holder lists, the price stays under the cap, and the
buyer pays no more than `max_price`. No money moves yet; the contract records the
price. A ticket has one transaction in flight at a time (`409 action_in_progress`
otherwise). Once sold, the seller's check-in passes no longer work; the buyer
makes a new one.

```bash
./scripts/create_user.py bob
./scripts/list_ticket.py <ticket-id> --price 110000000 --token <alice-token>
./scripts/buy_resale.py <ticket-id> --token <bob-token>
./scripts/list_ticket.py <ticket-id> --price 90000000 --token <bob-token>
./scripts/cancel_listing.py <ticket-id> --token <bob-token>
```

## Milestone 9

The service as a container. The [Dockerfile](../Dockerfile) builds a
static binary in a Go stage and copies it into `distroless/static`, running as
non-root (about 25 MB). The binary needs no files: migrations, the OpenAPI spec
and the contract bytecode are compiled in, and configuration comes only from the
environment. Its health check runs `/ticketing healthcheck`, which asks the
service's own `/readyz`, since the image has no shell or curl. One compose file
runs Besu, Postgres and the service; it takes `KEY_ENCRYPTION_KEY` and
`ADMIN_TOKEN` from your shell:

```bash
docker compose -f docker/docker-compose.yaml up -d --build --wait
```

The scripts above then work unchanged against `localhost:8080`. Service logs:

```bash
docker compose -f docker/docker-compose.yaml logs -f ticketing
```

## Distributed tracing

The service and Besu send OpenTelemetry spans to Jaeger,
which compose also runs. Each API request is one trace: its SQL queries, its
JSON-RPC calls to Besu with Besu's own spans inside them, and, joined later
from the outbox, the send, each receipt check and the final-status hook. The
request's trace context is saved with its transaction (`transactions.traceparent`),
which is how the background loops add to the right trace. Besu reads the `b3`
header, not `traceparent`, so the service sends both. Health checks and the
loops' idle polling create no traces. Open the UI, pick service `ticketing`:

```bash
open http://localhost:16686
```

Tracing is on when `OTEL_EXPORTER_OTLP_ENDPOINT` is set (compose sets it). To
trace a `go run` instance, leave Jaeger running from compose:

```bash
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 ADMIN_TOKEN=1234 go run ./cmd/ticketing
```

## Demo frontend

A React + TypeScript page ([frontend/](../frontend)) with forms for
the scripts' first steps: create a user, an organizer (with `ADMIN_TOKEN`) and
an event, each with a table of what came back: the token (with a copy
button), the Postgres id, and the chain addresses (wallet, box office and gate,
event contract). An event's row shows `deploying` until its contract is mined.
A Tickets section buys a ticket for a user, and on each ticket lists it (at
most the resale cap, filled in), cancels the listing, or buys it as another
user. Each row shows the owner by name and its latest transaction as it goes
from `signed` to `mined`, with the block number. **Show pass** on a ticket
makes the holder's signed check-in pass and shows it as a QR code with its
five-minute countdown. The Gate section scans it with the camera (or from an
image, or takes it handed over or pasted), finds the event's organizer from
the pass's contract address, and checks the ticket in with that organizer's
token; the ticket's row turns `checked in` once the burn is mined.
**Show refused actions** (in the header) offers every action on every ticket,
as any user, plus pass tampering and the wrong gate, and lists the error paths
to try with the code each returns (for example `422 buyer_is_seller`,
`409 already_checked_in`, `422 pass_not_from_holder`).
The Vite dev server proxies `/api` to the service, so the browser sees one
origin and the API needs no CORS. With the service running on port 8080:

```bash
cd frontend && npm install && npm run dev
```

Then open http://localhost:5173. What the page creates is kept in that
browser's local storage until you press Clear; point it at another service with
`TICKETING_URL=http://host:port npm run dev`.

The compose file also runs the page, built and served by nginx, at
http://localhost:3000; `npm run dev` is for working on it.

## Development

After editing the spec, regenerate the server code:

```bash
go generate ./internal/api
```

Run the Go tests (chain tests skip if Besu is down, database tests if Postgres is
down; `-count=1` skips Go's test cache):

```bash
go test -count=1 ./...
```

More commands for inspecting the chain (`cast`, Besu logs, scripts): [dev-tools.md](dev-tools.md).
