-- name: InsertTransaction :exec
INSERT INTO transactions (id, account_id, transaction_date, amount, transaction_type, processed_at)
VALUES ($1, $2, $3, $4, $5, NOW())
ON CONFLICT (id) DO NOTHING;

-- name: GetAccountBalance :one
SELECT COALESCE(SUM(amount), 0) as balance
FROM transactions
WHERE account_id = $1;

-- name: GetTransactionsByAccount :many
SELECT id, account_id, transaction_date, amount, transaction_type, processed_at
FROM transactions
WHERE account_id = $1
ORDER BY transaction_date DESC, processed_at DESC;

-- name: GetTransactionsByDateRange :many
SELECT id, account_id, transaction_date, amount, transaction_type, processed_at
FROM transactions
WHERE account_id = $1 
  AND transaction_date >= $2 
  AND transaction_date <= $3
ORDER BY transaction_date DESC, processed_at DESC;