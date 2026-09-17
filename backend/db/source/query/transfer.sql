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
