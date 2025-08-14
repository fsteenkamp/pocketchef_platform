CREATE TABLE account(
    id          TEXT PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    is_admin    BOOL NOT NULL DEFAULT false
);
