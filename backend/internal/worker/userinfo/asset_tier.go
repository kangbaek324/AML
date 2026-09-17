package userinfo

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
)

const (
	assetTierInterval = time.Hour

	// 단위: 원. docs 참고 - LOW ~2,000만원 / MEDIUM 2,000만원~2억원 / HIGH 2억원 이상
	assetTierLowMax    uint64 = 20_000_000
	assetTierMediumMax uint64 = 200_000_000
)

// AssetTierWorker는 유저의 현재 보유 현금 + 보유 주식 평가금액을 합산해
// average_asset을 갱신하고, 그 값에 따라 asset_tier를 재계산한다.
type AssetTierWorker struct {
	queries       *db.Queries
	sourceQueries *sourcedb.Queries
}

func NewAssetTierWorker(queries *db.Queries, sourceQueries *sourcedb.Queries) *AssetTierWorker {
	return &AssetTierWorker{queries: queries, sourceQueries: sourceQueries}
}

func (w *AssetTierWorker) Start(ctx context.Context) {
	w.run(ctx)

	ticker := time.NewTicker(assetTierInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.run(ctx)
		}
	}
}

func (w *AssetTierWorker) run(ctx context.Context) {
	userIDs, err := w.sourceQueries.ListUserIDs(ctx)
	if err != nil {
		log.Printf("asset tier worker: list users failed: %v", err)
		return
	}

	for _, userID := range userIDs {
		if err := w.updateUser(ctx, userID); err != nil {
			log.Printf("asset tier worker: update user %d failed: %v", userID, err)
		}
	}
}

func (w *AssetTierWorker) updateUser(ctx context.Context, userID int32) error {
	cash, err := w.sourceQueries.GetUserCashBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("get cash balance: %w", err)
	}

	stockValue, err := w.sourceQueries.GetUserStockValue(ctx, userID)
	if err != nil {
		return fmt.Errorf("get stock value: %w", err)
	}

	asset := toUint64(cash) + toUint64(stockValue)

	return w.queries.UpsertUserAssetTier(ctx, db.UpsertUserAssetTierParams{
		ID:           uint32(userID),
		AverageAsset: strconv.FormatUint(asset, 10),
		AssetTier:    classifyAssetTier(asset),
	})
}

func classifyAssetTier(asset uint64) db.UsersAssetTier {
	switch {
	case asset > assetTierMediumMax:
		return db.UsersAssetTierHIGH
	case asset > assetTierLowMax:
		return db.UsersAssetTierMEDIUM
	default:
		return db.UsersAssetTierLOW
	}
}

// COALESCE(SUM(...), 0)의 반환 타입을 sqlc가 정적으로 추론하지 못해 any로
// 생성되므로, 드라이버가 실제로 주는 형태([]byte 또는 int64)를 여기서 직접 처리한다.
func toUint64(v any) uint64 {
	switch t := v.(type) {
	case []byte:
		n, _ := strconv.ParseUint(string(t), 10, 64)
		return n
	case int64:
		return uint64(t)
	case uint64:
		return t
	default:
		return 0
	}
}
