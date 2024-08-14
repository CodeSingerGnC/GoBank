-- name: CreateSession :execresult
INSERT INTO sessions (
  id,
  user_account,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at,
  created_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, DEFAULT
);

-- name: GetSession :one
SELECT *
FROM sessions
WHERE id = ?
LIMIT 1;