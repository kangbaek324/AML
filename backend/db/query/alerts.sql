-- name: ListAlerts :many
SELECT * FROM alerts
ORDER BY alerted_at DESC;

-- name: GetAlert :one
SELECT * FROM alerts
WHERE id = ?
LIMIT 1;

-- name: ListAlertTradeIDs :many
SELECT tradeId FROM alert_trades
WHERE alertId = ?;

-- name: ListAlertTransferIDs :many
SELECT transferId FROM alert_transfers
WHERE alertId = ?;

-- name: CreateAlert :execresult
INSERT INTO alerts (userId, type, reason, status)
VALUES (?, ?, ?, 'PENDING');

-- name: LinkAlertTrade :exec
INSERT INTO alert_trades (alertId, tradeId)
VALUES (?, ?);

-- name: CountValidAlertsByUser :many
SELECT userId, COUNT(*) AS count
FROM alerts
WHERE status = 'ABNORMAL'
  AND alerted_at >= sqlc.arg(period_start)
  AND alerted_at < sqlc.arg(period_end)
GROUP BY userId;
