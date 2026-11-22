SELECT t.id as ticket_id, left(e.id::text, 18) as event_id, e.name as e_name, t.token_id as tid, u.subject as u_name,
    left(sk.encrypted_private_key::text, 10) as encrypt_p_k, sk.purpose as sk_pur
FROM tickets t
JOIN events e ON e.id = t.event_id
JOIN users u ON u.id = t.buyer_user_id
JOIN signing_keys sk ON sk.id = u.wallet_key_id
ORDER BY t.created_at
