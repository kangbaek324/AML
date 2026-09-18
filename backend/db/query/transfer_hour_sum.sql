-- name: UpsertTransferHourSum :exec
INSERT INTO transfer_hour_sum (userId, hour, amount)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE amount = amount + VALUES(amount);

-- name: GetTransferHourSum :one
SELECT amount FROM transfer_hour_sum
WHERE userId = ? AND hour = ?
LIMIT 1;

-- name: ClaimTransferAlert :execresult
UPDATE transfer_hour_sum
SET alerted = TRUE
WHERE userId = ? AND hour = ? AND alerted = FALSE;
