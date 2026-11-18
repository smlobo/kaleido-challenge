SELECT nonce,
       from_address,
       status,
       attempts,
       get_byte(calldata, 67) AS token_id,
       left(tx_hash, 12)      AS tx_hash,
       created_at
FROM transactions
WHERE status IN ('signed', 'submitted')
ORDER BY from_address, nonce;