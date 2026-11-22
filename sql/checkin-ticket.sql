SELECT t.id as ticket_id, left(e.id::text, 13) as event_id, e.name as e_name, t.token_id as tid, u.subject as u_name,
    o.name as o_name, left(sk.encrypted_private_key::text, 10) as encrypt_p_k, sk.purpose as sk_pur
FROM tickets t
JOIN events e ON e.id = t.event_id
JOIN users u ON u.id = t.buyer_user_id
JOIN organizers o ON o.id = e.organizer_id
JOIN signing_keys sk ON sk.id = o.gate_key_id
ORDER BY t.created_at
