-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?
LIMIT 1;

-- name: UpsertUserAssetTier :exec
INSERT INTO users (id, average_asset, asset_tier)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
    average_asset = VALUES(average_asset),
    asset_tier = VALUES(asset_tier);

-- name: ListUserRiskLevels :many
SELECT id, risk_level FROM users
ORDER BY id;

-- name: UpdateUserRiskLevel :exec
UPDATE users
SET risk_level = ?
WHERE id = ?;
