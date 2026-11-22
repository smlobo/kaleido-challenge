SELECT id, purpose, address, left(encrypted_private_key::text, 10) as encrypt_p_k, next_nonce as nxn
FROM signing_keys
ORDER BY created_at
