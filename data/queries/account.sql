-- name: AccountGetByID :one
SELECT
    id,
    email,
    phone_number,
    is_admin,
    is_root,
    created_at,
    updated_at,
    last_active,
    first_name,
    last_name,
    chef_status,
    provider,
    provider_last_refresh,
    picture,
    disabled,
    is_archived,
    archived_at,
    archived_by
FROM account WHERE id = $1;

-- name: AccountGetByEmail :one
SELECT * FROM account WHERE email = $1;

-- name: AccountCreate :exec
INSERT INTO account (
    id,
    email,
    is_admin,
    is_root,
    verified,
    password_hash,
    verify_code_hash
)
VALUES
    ($1, $2, $3, $4, $5, $6, $7);

-- name: AccountSetLastActive :exec
UPDATE account
SET
    last_active = $2
WHERE id = $1;

-- name: AccountList :many
SELECT
    id,
    email,
    phone_number,
    is_admin,
    is_root,
    created_at,
    updated_at,
    last_active,
    first_name,
    last_name,
    chef_status,
    provider,
    provider_last_refresh,
    picture,
    disabled,
    is_archived,
    archived_at,
    archived_by
FROM
    account;

-- name: AccountInit :one
SELECT
    id,
    email,
    phone_number,
    is_admin,
    is_root,
    created_at,
    updated_at,
    last_active,
    first_name,
    last_name,
    chef_status,
    provider,
    provider_last_refresh,
    picture,
    disabled,
    is_archived,
    archived_at,
    archived_by
FROM account WHERE id = $1;

-- name: AccountSetVerified :exec
UPDATE account
SET verified = true
WHERE verify_code_hash = $1;
