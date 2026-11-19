SELECT left(id::text, 13) as id, left(organizer_id::text, 13) as organizer_id, name,
    left(contract_address::text, 16) as contract_address, status, max_supply as supply, tickets_issued as issue
FROM events
ORDER BY created_at
