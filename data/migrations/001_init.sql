CREATE TABLE account(
    id                          text primary key,
    email                       text unique not null,
    verified                    bool not null default false,
    verify_code_hash            text,
    profile_configured          bool not null default false,
    phone_number                text,
    is_admin                    bool not null default false,
    is_root                     bool not null default false,
    created_at                  timestamp not null default (now() at time zone 'utc'),
    updated_at                  timestamp not null default (now() at time zone 'utc'),
    last_active                 timestamp,
    first_name                  text,
    last_name                   text,
    password_hash               text,
    provider                    text check(provider in ('google', 'microsoft', 'okta', 'apple', 'facebook')),
    provider_token              text,
    provider_refresh_token      text,
    provider_last_refresh       timestamp,
    picture                     text,
    disabled                    bool not null default false,
    is_archived                 bool not null default false,
    archived_at                 timestamp,
    archived_by                 text references account(id)
);

CREATE TABLE session (
    id                  text primary key,
    token_hash          text not null,
    account_id          text not null references account(id),
    created_at          timestamp not null default (now() at time zone 'utc'),
    updated_at          timestamp not null default (now() at time zone 'utc'),
    expires_at          timestamp not null,
    invalidated         bool not null default false,
    invalidated_at      timestamp,
    invalidated_by      text references account(id)
);

CREATE TABLE chef (
    id                          text primary key,
    account_id                  text not null references account(id),
    display_name                text,
    description                 text not null,
    picture                     text not null,
    phone_number                text not null,
    chef_status                 text not null check(chef_status in ('pending', 'rejected', 'verified', 'disabled')) default 'none',
    created_at                  timestamp not null default (now() at time zone 'utc'),
    updated_at                  timestamp not null default (now() at time zone 'utc'),
    archived_at                 timestamp,
    archived_by                 text references account(id),
    social_link_instagram       text,
    social_link_facebook        text,
    social_link_website         text,
    social_link_x               text,
    social_link_tiktok          text,
    social_link_youtube         text
);

CREATE TABLE chef_profile_review (
    id                      text primary key,
    chef_id                 text not null references chef(id),
    created_by              text not null references account(id),
    created_at              timestamp not null default (now() at time zone 'utc'),
    reviewed_at             timestamp not null default (now() at time zone 'utc'),
    reviewer                text references account(id),
    review_private_note     text,
    review_public_note      text,
    review_outcome          text not null check(review_outcome in ('rejected', 'verified', 'disabled')) default 'none'
);

CREATE TABLE recipe (
    id              text primary key,
    chef_id         text not null references chef(id),
    cover_img_id    text not null,
    content         jsonb not null default '{}',
    tags            text[],
    is_published    bool not null default false,
    is_hidden       bool not null default false,
    created_at      timestamp not null default (now() at time zone 'utc'),
    updated_at      timestamp not null default (now() at time zone 'utc')
);

CREATE TABLE pocket (
    id              text primary key,
    name            text not null,
    account_id      text not null references account(id),
    is_archived     bool not null default false,
    created_at      timestamp not null default (now() at time zone 'utc'),
    updated_at      timestamp not null default (now() at time zone 'utc')
);

CREATE TABLE pocket_recipe ();

CREATE TABLE collection ();

CREATE TABLE featured_chef ();

CREATE TABLE subscription ();
CREATE TABLE subscription_plan ();
CREATE TABLE invoice ();
