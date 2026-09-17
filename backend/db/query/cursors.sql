-- name: GetCursor :one
SELECT timestamp FROM cursors
WHERE type = ?
LIMIT 1;

-- name: UpsertCursor :exec
INSERT INTO cursors (type, timestamp)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE timestamp = VALUES(timestamp);
