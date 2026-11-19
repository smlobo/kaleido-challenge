-- Milestone 6: the first business endpoints.

-- API tokens. Only a SHA-256 of each token is stored; the token itself is
-- shown once, when its organizer or user is created. Each token belongs to
-- exactly one organizer or one user.
CREATE TABLE api_tokens (
    token_hash   bytea PRIMARY KEY CHECK (length(token_hash) = 32),
    organizer_id uuid REFERENCES organizers (id),
    user_id      uuid REFERENCES users (id),
    created_at   timestamptz NOT NULL DEFAULT now(),
    CHECK ((organizer_id IS NULL) <> (user_id IS NULL))
);

-- Ticket numbers are handed out by incrementing this counter, which also
-- enforces the supply: concurrent buyers cannot both get the last ticket.
ALTER TABLE events ADD COLUMN tickets_issued integer NOT NULL DEFAULT 0
    CHECK (tickets_issued >= 0 AND tickets_issued <= max_supply);

-- The check-in transaction, if one was sent.
ALTER TABLE tickets ADD COLUMN checkin_tx_id uuid REFERENCES transactions (id);

-- Who asked for the transaction ('organizer:<id>' or 'user:<id>'), so only
-- they can read its operation. NULL for transactions the service started.
ALTER TABLE transactions ADD COLUMN requested_by text;

-- The final-status hook finds the event or ticket a transaction belongs to.
CREATE INDEX events_deploy_tx ON events (deploy_tx_id);
CREATE INDEX tickets_mint_tx ON tickets (mint_tx_id);
CREATE INDEX tickets_checkin_tx ON tickets (checkin_tx_id);
