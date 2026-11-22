SELECT e.id as event_id, e.name as e_name, o.name as o_name,
    left(sk.encrypted_private_key::text, 10) as encrypt_p_k, sk.purpose as sk_pur
FROM events e
JOIN organizers o ON o.id = e.organizer_id
JOIN signing_keys sk ON sk.id = o.box_office_key_id
ORDER BY e.created_at
