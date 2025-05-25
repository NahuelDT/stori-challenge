-- name: CreateAccount :one
INSERT INTO accounts (id, email, created_at)
VALUES ($1, $2, NOW())
ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
RETURNING id;

-- name: GetAccountByEmail :one
SELECT id, email, created_at
FROM accounts
WHERE email = $1;

-- name: GetAccountByID :one
SELECT id, email, created_at
FROM accounts
WHERE id = $1;