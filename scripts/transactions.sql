SELECT block_number as bnum, nonce, from_address, status, contract_address
FROM transactions
ORDER BY created_at;
