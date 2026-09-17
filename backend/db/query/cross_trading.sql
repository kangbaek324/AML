-- name: UpsertCrossTradingCount :execresult
INSERT INTO cross_trading_count (userId, date, count)
VALUES (?, ?, LAST_INSERT_ID(1))
ON DUPLICATE KEY UPDATE count = LAST_INSERT_ID(count + 1);

-- name: ClaimCrossTradingAlert :execresult
UPDATE cross_trading_count
SET alerted = TRUE
WHERE userId = ? AND date = ? AND alerted = FALSE;
