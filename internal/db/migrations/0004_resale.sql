-- Milestone 7: resale.

-- Every listing, cancellation and resale purchase is one transaction sent
-- from a user's wallet. This ties it to its ticket, so the final-status hook
-- knows what a mined transaction changed: a listing sets the ticket's
-- listing_price, a cancellation clears it, and a purchase moves the ticket to
-- the buyer's wallet. A reverted or failed one changes nothing.
CREATE TABLE resale_actions (
    tx_id      uuid PRIMARY KEY REFERENCES transactions (id),
    ticket_id  uuid NOT NULL REFERENCES tickets (id),
    kind       text NOT NULL CHECK (kind IN ('list', 'cancel', 'buy')),
    user_id    uuid NOT NULL REFERENCES users (id),
    -- list: the asking price; buy: the most the buyer agreed to pay.
    price      uint256,
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK ((kind = 'cancel') = (price IS NULL))
);

-- The service allows one chain action per ticket at a time, and finds the
-- ones still in flight through this.
CREATE INDEX resale_actions_ticket ON resale_actions (ticket_id);

-- GET /events/{id}/listings.
CREATE INDEX tickets_listed ON tickets (event_id) WHERE listing_price IS NOT NULL;
