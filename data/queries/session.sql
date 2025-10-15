-- name: SessionGetValid :one
SELECT * FROM session
WHERE
    id = $1 AND
    expires_at > $2 AND
    invalidated = FALSE AND
    account_id = $3
;

-- name: SessionCreate :exec
INSERT INTO session (
    id,
    account_id,
    expires_at,
    token_hash,
    created_at
)
VALUES
    ($1, $2, $3, $4, $5);

-- name: SessionInvalidate :exec
UPDATE session
SET
    invalidated = true,
    invalidated_at = $1
WHERE
    id = $2 AND
    account_id = $3;

-- name: SessionGetFromTokenHash :one
SELECT * FROM session
WHERE token_hash = $1;
