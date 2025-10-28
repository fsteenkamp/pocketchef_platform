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
    account.id                              id,
    account.email                           email,
    account.phone_number                    phone_number,
    account.is_admin                        is_admin,
    account.is_root                         is_root,
    account.created_at                      created_at,
    account.updated_at                      updated_at,
    account.last_active                     last_active,
    account.first_name                      first_name,
    account.last_name                       last_name,
    account.provider                        provider,
    account.provider_last_refresh           provider_last_refresh,
    account.picture                         picture,
    account.disabled                        disabled,
    account.is_archived                     is_archived,
    account.archived_at                     archived_at,
    account.archived_by                     archived_by,
    json_agg(json_build_object(
        'id',                               chef.id,
        'account_id',                       chef.account_id,
        'display_name',                     chef.display_name,
        'description',                      chef.description,
        'picture',                          chef.picture,
        'phone_number',                     chef.phone_number,
        'chef_status',                      chef.chef_status,
        'created_at',                       chef.created_at,
        'updated_at',                       chef.updated_at,
        'archived_at',                      chef.archived_at,
        'archived_by',                      chef.archived_by,
        'social_link_instagram',            chef.social_link_instagram,
        'social_link_facebook',             chef.social_link_facebook,
        'social_link_website',              chef.social_link_website,
        'social_link_x',                    chef.social_link_x,
        'social_link_tiktok',               chef.social_link_tiktok,
        'social_link_youtube',              chef.social_link_youtube
    )) FILTER (WHERE chef.id IS NOT NULL)   chef_profile
FROM account
LEFT JOIN chef
ON
    chef.account_id = account.id
WHERE account.id = $1
GROUP BY account.id
;

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

-- name: AccountSetAdmin :exec
UPDATE account
SET
    is_admin = $2
WHERE
    id = $1;

-- name: AccountSetRoot :exec
UPDATE account
SET
    is_root  =  $2,
    is_admin =  $2
WHERE
    id = $1;
