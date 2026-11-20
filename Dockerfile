# The ticketing service as a small, non-root container.
#
# Build from the repository root:
#   docker build -t kaleido-ticketing .
# or let compose build it: docker compose -f docker/docker-compose.yaml up -d --build --wait
#
# The binary needs no files at runtime: migrations, the OpenAPI spec and the
# contract bytecode (through the committed abigen bindings) are compiled in,
# and all configuration comes from the environment. So the runtime image has
# no shell, no solc and nothing else to patch.

# --- build ---
FROM golang:1.27.0-alpine3.23 AS build
WORKDIR /src

# Dependencies first, so editing the code does not re-download them.
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
# CGO off: a static binary that runs on distroless/static. go-ethereum then
# uses its pure-Go secp256k1, which signs the same way.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ticketing ./cmd/ticketing

# --- runtime ---
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/ticketing /ticketing
USER nonroot:nonroot
EXPOSE 8080
# No curl in the image: the binary asks its own /readyz.
HEALTHCHECK --interval=5s --timeout=5s --start-period=30s --retries=3 CMD ["/ticketing", "healthcheck"]
ENTRYPOINT ["/ticketing"]
