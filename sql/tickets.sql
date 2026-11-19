SELECT left(id::text, 13) as id, left(event_id::text, 13) as event_id, token_id as tID,
    left(mint_tx_id::text, 13) as mint_tx_id, left(buyer_user_id::text, 13) as buyer_user_id,
    left(owner_address::text, 16) as owner_address, status
FROM tickets
ORDER BY created_at
