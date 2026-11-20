SELECT left(ticket_id::text, 13) as ticket_id, kind, left(user_id::text, 13) as user_id, price
FROM resale_actions
ORDER BY created_at
