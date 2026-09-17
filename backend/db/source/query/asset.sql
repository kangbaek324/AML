-- name: ListUserIDs :many
SELECT id FROM users
ORDER BY id;

-- name: GetUserCashBalance :one
SELECT COALESCE(SUM(balance), 0) AS cash
FROM accounts
WHERE user_id = ?;

-- name: GetUserStockValue :one
SELECT COALESCE(SUM(us.quantity * s.price), 0) AS stock_value
FROM user_stocks us
JOIN accounts a ON a.id = us.account_id
JOIN stocks s ON s.id = us.stock_id
WHERE a.user_id = ?;
