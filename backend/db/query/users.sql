-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?
LIMIT 1;

-- name: UpsertUserAssetTier :exec
-- updated_at을 명시적으로 지정해, 값이 안 바뀌어도 "언제 마지막으로 체크했는지"가 갱신되도록 한다.
-- (MySQL은 ON DUPLICATE KEY UPDATE에서 실제 값 변경이 없으면 ON UPDATE CURRENT_TIMESTAMP를 건너뛴다)
INSERT INTO users (id, average_asset, asset_tier)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
    average_asset = VALUES(average_asset),
    asset_tier = VALUES(asset_tier),
    updated_at = CURRENT_TIMESTAMP;

-- name: ListUserRiskLevels :many
SELECT id, risk_level FROM users
ORDER BY id;

-- name: UpdateUserRiskLevel :exec
UPDATE users
SET risk_level = ?
WHERE id = ?;
