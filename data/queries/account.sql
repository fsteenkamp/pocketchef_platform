-- name: AccountGetByID :one
SELECT * FROM account WHERE id = $1;

-- name: AccountGetByEmail :one
SELECT * FROM account WHERE email = $1;

-- name: AccountCreate :exec
INSERT INTO account (
    id,
    email,
    is_admin,
    is_root
)
VALUES
    ($1, $2, $3, $4);
