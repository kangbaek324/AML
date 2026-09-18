-- name: ListTransfersSince :many
SELECT
    t.id AS transfer_id,
    t.amount,
    t.completed_at,
    ra.user_id AS recipient_user_id
FROM transfers t
JOIN accounts ra ON ra.id = t.recipient_account_id
WHERE t.status = 'COMPLETED'
  AND t.completed_at > sqlc.arg(completed_at)
ORDER BY t.completed_at;

-- name: ListReceivedTransfersByUserAndHour :many
SELECT
    t.id AS transfer_id,
    t.amount,
    t.completed_at
FROM transfers t
JOIN accounts ra ON ra.id = t.recipient_account_id
WHERE ra.user_id = sqlc.arg(user_id)
  AND t.status = 'COMPLETED'
  AND t.completed_at >= sqlc.arg(hour_start)
  AND t.completed_at < sqlc.arg(hour_end)
ORDER BY t.completed_at;

-- name: GetTransferDetail :one
SELECT
    t.id AS transfer_id,
    t.amount,
    t.status,
    t.created_at,
    t.completed_at,
    sa.user_id AS sender_user_id,
    ra.user_id AS recipient_user_id
FROM transfers t
JOIN accounts sa ON sa.id = t.sender_account_id
JOIN accounts ra ON ra.id = t.recipient_account_id
WHERE t.id = ?
LIMIT 1;
