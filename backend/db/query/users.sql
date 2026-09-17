-- name: UpsertUserAssetTier :exec
INSERT INTO users (id, average_asset, asset_tier)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
    average_asset = VALUES(average_asset),
    asset_tier = VALUES(asset_tier);
