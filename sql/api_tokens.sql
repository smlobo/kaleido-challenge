SELECT left(token_hash::text, 14) as token_hash, organizer_id, user_id
FROM api_tokens
ORDER BY created_at
