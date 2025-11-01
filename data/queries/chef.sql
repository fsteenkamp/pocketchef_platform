-- name: ChefProfileCreate :exec
INSERT INTO chef (
    id,
    account_id,
    display_name,
    description
)
VALUES ($1, $2, $3, $4);
