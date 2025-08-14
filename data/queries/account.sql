-- name: AccountGetByID :one
SELECT * FROM account WHERE id = $1;
