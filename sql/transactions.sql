SELECT left(id::text, 13) as id, block_number as bnum, nonce as non, left(from_address::text, 16) as from_address,
    left(to_address::text, 16) as to_address, status, contract_address
FROM transactions
ORDER BY created_at;
