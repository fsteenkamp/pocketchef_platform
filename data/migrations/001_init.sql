CREATE TABLE account(
    id              TEXT PRIMARY KEY,
    email           TEXT UNIQUE NOT NULL,
    is_admin        BOOL NOT NULL DEFAULT false,
    chef_status     TEXT NOT NULL CHECK(chef_status IN ('none', 'pending', 'verified', 'disabled'))
);

CREATE TABLE recipe ();

CREATE TABLE pocket ();

CREATE TABLE pocket_recipe ();

CREATE TABLE collection ();

CREATE TABLE featured_chef ();
