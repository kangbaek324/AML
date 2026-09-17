-- name: CountValidAlertsByUser :many
SELECT userId, COUNT(*) AS count
FROM alerts
WHERE status = 'ABNORMAL'
  AND alerted_at >= sqlc.arg(period_start)
  AND alerted_at < sqlc.arg(period_end)
GROUP BY userId;
