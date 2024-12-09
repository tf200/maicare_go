-- name: CreateSession :one
INSERT INTO sessions (
    id,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expires_at,
    created_at,
    user_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;