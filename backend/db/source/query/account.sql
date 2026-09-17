-- name: ListAccountsByUser :many
SELECT id, account_number, balance, available_balance, status
FROM accounts
WHERE user_id = ?
ORDER BY id;

-- name: GetAccount :one
SELECT id, user_id, account_number, balance, available_balance, status
FROM accounts
WHERE id = ?
LIMIT 1;

-- name: ListStockHoldingsByAccount :many
SELECT s.id AS stock_id, s.name, us.quantity, us.average, s.price
FROM user_stocks us
JOIN stocks s ON s.id = us.stock_id
WHERE us.account_id = ?
ORDER BY s.id;
