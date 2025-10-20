CREATE TABLE account(
    id                          TEXT PRIMARY KEY,
    email                       TEXT UNIQUE NOT NULL,
    email_verified              BOOL NOT NULL DEFAULT FALSE,
    phone_number                TEXT,
    is_admin                    BOOL NOT NULL DEFAULT false,
    is_root                     BOOL NOT NULL DEFAULT false,
    created_at                  TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at                  TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    last_active                 TIMESTAMP,
    first_name                  TEXT,
    last_name                   TEXT,
    chef_status                 TEXT NOT NULL CHECK(chef_status IN ('none', 'pending', 'declined', 'verified', 'disabled')) DEFAULT 'none',
    password_hash               TEXT,
    provider                    TEXT CHECK(provider IN ('google', 'microsoft', 'okta', 'apple', 'facebook')),
    provider_token              TEXT,
    provider_refresh_token      TEXT,
    provider_last_refresh       TIMESTAMP,
    picture                     TEXT,
    disabled                    BOOL NOT NULL DEFAULT false,
    is_archived                 BOOL NOT NULL DEFAULT false,
    archived_at                 TIMESTAMP,
    archived_by                 TEXT references account(id)
);

CREATE TABLE session (
    id                  TEXT PRIMARY KEY,
    token_hash          TEXT NOT NULL,
    account_id          TEXT NOT NULL references account(id),
    created_at          TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at          TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    expires_at          TIMESTAMP NOT NULL,
    invalidated         BOOL NOT NULL DEFAULT false,
    invalidated_at      TIMESTAMP,
    invalidated_by      TEXT references account(id)
);

CREATE TABLE recipe (
    id              TEXT PRIMARY KEY,
    chef_id         TEXT NOT NULL references account(id),
    cover_img_id    TEXT NOT NULL,
    content         JSONB NOT NULL default '{}',
    tags            TEXT[],
    is_published    BOOL NOT NULL DEFAULT FALSE,
    is_hidden       BOOL NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at      TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

CREATE TABLE pocket (
    id              TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    account_id      TEXT NOT NULL references account(id),
    is_archived     BOOL NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at      TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

CREATE TABLE pocket_recipe ();

CREATE TABLE collection ();

CREATE TABLE featured_chef ();
