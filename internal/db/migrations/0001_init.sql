-- Initial schema. Postgres is the only state store: everything the chain
-- cannot hold (keys, users, event details) plus a queryable copy of what it
-- does (owners, check-ins). Addresses and hashes are lowercase 0x-hex text.
-- uint256 values (token IDs, prices) are numeric(78,0).

CREATE DOMAIN eth_address AS text CHECK (VALUE ~ '^0x[0-9a-f]{40}$');
CREATE DOMAIN eth_hash    AS text CHECK (VALUE ~ '^0x[0-9a-f]{64}$');
CREATE DOMAIN uint256     AS numeric(78, 0) CHECK (VALUE >= 0);

-- Every key the service holds. The private key is AES-256-GCM encrypted under
-- KEY_ENCRYPTION_KEY (env), with the address as associated data, so a
-- ciphertext copied onto another row fails to decrypt.
CREATE TABLE signing_keys (
    id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    address               eth_address NOT NULL UNIQUE,
    purpose               text NOT NULL CHECK (purpose IN ('deployer', 'box_office', 'payout', 'gate', 'wallet')),
    encrypted_private_key bytea NOT NULL,
    next_nonce            bigint NOT NULL DEFAULT 0 CHECK (next_nonce >= 0),
    created_at            timestamptz NOT NULL DEFAULT now()
);

-- Exactly one deployer, even when several instances bootstrap at once.
CREATE UNIQUE INDEX signing_keys_one_deployer ON signing_keys (purpose) WHERE purpose = 'deployer';

CREATE TABLE organizers (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name              text NOT NULL,
    box_office_key_id uuid NOT NULL UNIQUE REFERENCES signing_keys (id),
    payout_key_id     uuid NOT NULL UNIQUE REFERENCES signing_keys (id),
    gate_key_id       uuid NOT NULL UNIQUE REFERENCES signing_keys (id),
    created_at        timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    subject       text NOT NULL UNIQUE, -- login identity (bearer token owner now, OIDC sub later)
    wallet_key_id uuid NOT NULL UNIQUE REFERENCES signing_keys (id),
    created_at    timestamptz NOT NULL DEFAULT now()
);

-- Outbox of signed transactions. Its id is the operation ID the API returns.
-- A row is committed before the transaction is sent, so any instance can
-- (re)send raw_tx after a crash.
CREATE TABLE transactions (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key text UNIQUE,
    from_address    eth_address NOT NULL,
    nonce           bigint NOT NULL CHECK (nonce >= 0),
    to_address      eth_address,          -- NULL for a contract creation
    calldata        bytea NOT NULL,
    raw_tx          bytea NOT NULL,
    tx_hash         eth_hash NOT NULL UNIQUE,
    status          text NOT NULL DEFAULT 'signed'
                    CHECK (status IN ('signed', 'submitted', 'mined', 'reverted', 'failed')),
    block_number    bigint,
    contract_address eth_address,         -- from the receipt of a creation
    revert_reason   text,
    attempts        integer NOT NULL DEFAULT 0,
    next_attempt_at timestamptz NOT NULL DEFAULT now(),
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    -- Safety net: a nonce handed out twice fails here instead of producing
    -- two conflicting signed transactions.
    UNIQUE (from_address, nonce)
);

-- Work queue for the submitter and receipt loops (FOR UPDATE SKIP LOCKED).
CREATE INDEX transactions_pending ON transactions (next_attempt_at)
    WHERE status IN ('signed', 'submitted');

CREATE TABLE events (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id     uuid NOT NULL REFERENCES organizers (id),
    name             text NOT NULL,
    symbol           text NOT NULL,
    starts_at        timestamptz NOT NULL,
    max_supply       integer NOT NULL CHECK (max_supply > 0),
    face_value       uint256 NOT NULL,
    max_resale_bps   integer NOT NULL CHECK (max_resale_bps >= 0),
    contract_address eth_address UNIQUE,  -- filled when the deploy is mined
    contract_version integer NOT NULL DEFAULT 1, -- which TicketNFT ABI
    deploy_tx_id     uuid REFERENCES transactions (id),
    status           text NOT NULL DEFAULT 'deploying'
                     CHECK (status IN ('deploying', 'on_sale', 'failed')),
    created_at       timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE tickets (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id      uuid NOT NULL REFERENCES events (id),
    token_id      uint256 NOT NULL,
    buyer_user_id uuid NOT NULL REFERENCES users (id),
    record        jsonb NOT NULL,         -- canonical off-chain record
    record_hash   eth_hash NOT NULL,      -- keccak256 of record, stored on-chain
    owner_address eth_address,            -- from indexed Transfer events only
    listing_price uint256,                -- open resale listing, from indexed events
    status        text NOT NULL DEFAULT 'pending_mint'
                  CHECK (status IN ('pending_mint', 'active', 'checked_in', 'failed')),
    mint_tx_id    uuid REFERENCES transactions (id),
    created_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (event_id, token_id)
);

-- Every contract event the indexer has seen. The primary key makes
-- re-indexing the same block a no-op.
CREATE TABLE chain_events (
    tx_hash      eth_hash NOT NULL,
    log_index    integer NOT NULL,
    contract     eth_address NOT NULL,
    name         text NOT NULL,
    token_id     uint256,
    from_address eth_address,
    to_address   eth_address,
    block_number bigint NOT NULL,
    args         jsonb NOT NULL DEFAULT '{}',
    PRIMARY KEY (tx_hash, log_index)
);

CREATE INDEX chain_events_token ON chain_events (contract, token_id);

-- Where the indexer resumes, and which chain it was indexing (genesis hash
-- detects a chain reset).
CREATE TABLE indexer_checkpoint (
    chain_id     bigint PRIMARY KEY,
    genesis_hash eth_hash NOT NULL,
    last_block   bigint NOT NULL DEFAULT -1
);
