package rule

import (
	"context"
	"fmt"
	"log"
	"time"

	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
)

const largeTransactionInterval = 10 * time.Second

// 단일 거래 임계값 (k=0.5). docs 참고 - LOW 1,000만원 / MEDIUM 5,000만원 / HIGH 1억원
var largeTransactionThreshold = map[db.UsersAssetTier]uint64{
	db.UsersAssetTierLOW:    10_000_000,
	db.UsersAssetTierMEDIUM: 50_000_000,
	db.UsersAssetTierHIGH:   100_000_000,
}

// LargeTransactionRule은 단일 주식 거래(체결) 금액이 해당 유저 Asset Tier의
// 임계값 이상이면 Alert를 발생시킨다.
type LargeTransactionRule struct {
	queries       *db.Queries
	sourceQueries *sourcedb.Queries
}

func NewLargeTransactionRule(queries *db.Queries, sourceQueries *sourcedb.Queries) *LargeTransactionRule {
	return &LargeTransactionRule{queries: queries, sourceQueries: sourceQueries}
}

func (r *LargeTransactionRule) Name() string { return "large_transaction" }

func (r *LargeTransactionRule) Interval() time.Duration { return largeTransactionInterval }

func (r *LargeTransactionRule) Run(ctx context.Context, since time.Time) (time.Time, error) {
	trades, err := r.sourceQueries.ListTradesSince(ctx, since)
	if err != nil {
		return since, fmt.Errorf("list trades: %w", err)
	}
	if len(trades) == 0 {
		return since, nil
	}

	tiers, err := r.loadAssetTiers(ctx)
	if err != nil {
		return since, fmt.Errorf("load asset tiers: %w", err)
	}

	for _, t := range trades {
		amount := t.Quantity * t.Price

		r.checkUser(ctx, uint32(t.MakerUserID), amount, t.TradeID, tiers)
		if t.TakerUserID != t.MakerUserID {
			r.checkUser(ctx, uint32(t.TakerUserID), amount, t.TradeID, tiers)
		}
	}

	// trades는 matched_at 오름차순이므로 마지막 항목이 가장 최근 거래다.
	return trades[len(trades)-1].MatchedAt, nil
}

func (r *LargeTransactionRule) loadAssetTiers(ctx context.Context) (map[uint32]db.UsersAssetTier, error) {
	users, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	tiers := make(map[uint32]db.UsersAssetTier, len(users))
	for _, u := range users {
		tiers[u.ID] = u.AssetTier
	}
	return tiers, nil
}

func (r *LargeTransactionRule) checkUser(
	ctx context.Context,
	userID uint32,
	amount uint64,
	tradeID int64,
	tiers map[uint32]db.UsersAssetTier,
) {
	tier, ok := tiers[userID]
	if !ok {
		tier = db.UsersAssetTierLOW
	}

	threshold := largeTransactionThreshold[tier]
	if amount < threshold {
		return
	}

	reason := fmt.Sprintf("단일 거래 금액 %d원이 %s 등급 임계값(%d원) 이상입니다.", amount, tier, threshold)

	result, err := r.queries.CreateAlert(ctx, db.CreateAlertParams{
		Userid: userID,
		Type:   db.AlertsTypeLARGETRANSACTION,
		Reason: reason,
	})
	if err != nil {
		log.Printf("large transaction rule: create alert failed: %v", err)
		return
	}

	alertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("large transaction rule: get alert id failed: %v", err)
		return
	}

	if err := r.queries.LinkAlertTrade(ctx, db.LinkAlertTradeParams{
		Alertid: uint64(alertID),
		Tradeid: uint64(tradeID),
	}); err != nil {
		log.Printf("large transaction rule: link alert trade failed: %v", err)
	}
}
