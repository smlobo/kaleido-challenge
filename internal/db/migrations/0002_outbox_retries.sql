-- Milestone 4: the submitter and receipt loops.
--
-- A transaction the node rejects outright (bad gas limit, wrong chain) is
-- marked failed and its nonce is reused by a zero-value transfer to self, so
-- later transactions from the same key are not stuck behind the gap. The
-- one-row-per-nonce safety net therefore only covers rows that can still be
-- mined.
ALTER TABLE transactions DROP CONSTRAINT transactions_from_address_nonce_key;
CREATE UNIQUE INDEX transactions_live_nonce ON transactions (from_address, nonce)
    WHERE status <> 'failed';

-- The most recent send error. Kept while the submitter retries, cleared once a
-- send succeeds; for a failed row it says why.
ALTER TABLE transactions ADD COLUMN last_error text;
