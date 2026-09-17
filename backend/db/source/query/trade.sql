-- name: ListTradesSince :many
SELECT
    t.id AS trade_id,
    t.quantity,
    t.price,
    t.matched_at,
    ma.user_id AS maker_user_id,
    ka.user_id AS taker_user_id
FROM trades t
JOIN orders mo ON mo.id = t.maker_order_id
JOIN accounts ma ON ma.id = mo.account_id
JOIN orders ko ON ko.id = t.taker_order_id
JOIN accounts ka ON ka.id = ko.account_id
WHERE t.matched_at > ?
ORDER BY t.matched_at;

-- name: ListSelfTradesByUserAndDate :many
SELECT
    t.id AS trade_id,
    t.quantity,
    t.price,
    t.matched_at
FROM trades t
JOIN orders mo ON mo.id = t.maker_order_id
JOIN accounts ma ON ma.id = mo.account_id
JOIN orders ko ON ko.id = t.taker_order_id
JOIN accounts ka ON ka.id = ko.account_id
WHERE ma.user_id = sqlc.arg(user_id)
  AND ka.user_id = sqlc.arg(user_id)
  AND DATE(t.matched_at) = sqlc.arg(date)
ORDER BY t.matched_at;

-- name: GetTradeDetail :one
SELECT
    t.id AS trade_id,
    t.stock_id,
    t.quantity,
    t.price,
    t.matched_at,
    ma.user_id AS maker_user_id,
    ka.user_id AS taker_user_id
FROM trades t
JOIN orders mo ON mo.id = t.maker_order_id
JOIN accounts ma ON ma.id = mo.account_id
JOIN orders ko ON ko.id = t.taker_order_id
JOIN accounts ka ON ka.id = ko.account_id
WHERE t.id = ?
LIMIT 1;
