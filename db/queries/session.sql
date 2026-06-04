-- name: CreateSession :one
INSERT INTO sessions (
    user_id, refresh_token_hash, user_agent, expires_at
) VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: FindSessionByHash :one
SELECT *
    FROM sessions
    WHERE refresh_token_hash = $1
    AND expires_at > now();

-- name: RotateSession :one
UPDATE sessions
    SET refresh_token_hash = $1, expires_at = $2, last_used_at = now()
    WHERE refresh_token_hash = $3
    AND expires_at > now()
RETURNING *;

-- name: DeleteSessionByHash :exec
DELETE FROM sessions WHERE refresh_token_hash = $1;

-- name: DeleteSessionsByUserId :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < now();
