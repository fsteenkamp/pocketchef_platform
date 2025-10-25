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
    verify_code_hash,
    provider,
    provider_token,
    provider_refresh_token,
    provider_last_refresh,
    picture
)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

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

-- name: AccountRefreshProviderDetails :exec
UPDATE account
SET
    first_name = $2,
    last_name = $3,
    phone_number = $4,
    picture = $5
WHERE
    id = $1;
