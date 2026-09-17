package rule

import (
	"context"
	"fmt"
	"log"
	"time"

	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
)

const crossTradingInterval = 10 * time.Second

// 하루 동안 자전거래 횟수가 이 값을 초과하면 Alert를 발생시킨다.
const crossTradingAlertThreshold = 100

// CrossTradingRule은 동일 유저가 매수/매도 양쪽 주문의 당사자인 거래(자전거래)를
// 날짜별로 집계하다가, 하루 누적 횟수가 임계값을 초과하면 Alert를 발생시킨다.
type CrossTradingRule struct {
	queries       *db.Queries
	sourceQueries *sourcedb.Queries
}

func NewCrossTradingRule(queries *db.Queries, sourceQueries *sourcedb.Queries) *CrossTradingRule {
	return &CrossTradingRule{queries: queries, sourceQueries: sourceQueries}
}

func (r *CrossTradingRule) Name() string { return "cross_trading" }

func (r *CrossTradingRule) Interval() time.Duration { return crossTradingInterval }

func (r *CrossTradingRule) Run(ctx context.Context, since time.Time) (time.Time, error) {
	trades, err := r.sourceQueries.ListTradesSince(ctx, since)
	if err != nil {
		return since, fmt.Errorf("list trades: %w", err)
	}
	if len(trades) == 0 {
		return since, nil
	}

	for _, t := range trades {
		if t.MakerUserID != t.TakerUserID {
			continue
		}
		r.checkUser(ctx, uint32(t.MakerUserID), t.MatchedAt)
	}

	// trades는 matched_at 오름차순이므로 마지막 항목이 가장 최근 거래다.
	return trades[len(trades)-1].MatchedAt, nil
}

func (r *CrossTradingRule) checkUser(ctx context.Context, userID uint32, matchedAt time.Time) {
	date := time.Date(matchedAt.Year(), matchedAt.Month(), matchedAt.Day(), 0, 0, 0, 0, matchedAt.Location())

	result, err := r.queries.UpsertCrossTradingCount(ctx, db.UpsertCrossTradingCountParams{
		Userid: userID,
		Date:   date,
	})
	if err != nil {
		log.Printf("cross trading rule: upsert count failed: %v", err)
		return
	}

	count, err := result.LastInsertId()
	if err != nil {
		log.Printf("cross trading rule: get count failed: %v", err)
		return
	}
	if count <= crossTradingAlertThreshold {
		return
	}

	claimed, err := r.queries.ClaimCrossTradingAlert(ctx, db.ClaimCrossTradingAlertParams{
		Userid: userID,
		Date:   date,
	})
	if err != nil {
		log.Printf("cross trading rule: claim alert failed: %v", err)
		return
	}
	rows, err := claimed.RowsAffected()
	if err != nil {
		log.Printf("cross trading rule: claim alert rows affected failed: %v", err)
		return
	}
	if rows == 0 {
		// 이미 해당 날짜에 대해 Alert가 발생한 상태다.
		return
	}

	selfTrades, err := r.sourceQueries.ListSelfTradesByUserAndDate(ctx, sourcedb.ListSelfTradesByUserAndDateParams{
		UserID: int32(userID),
		Date:   date,
	})
	if err != nil {
		log.Printf("cross trading rule: list self trades failed: %v", err)
		return
	}

	reason := fmt.Sprintf("%s 하루 자전거래 횟수가 %d회로 임계값(%d회)을 초과했습니다.", date.Format("2006-01-02"), count, crossTradingAlertThreshold)

	alertResult, err := r.queries.CreateAlert(ctx, db.CreateAlertParams{
		Userid: userID,
		Type:   db.AlertsTypeCROSSTRADING,
		Reason: reason,
	})
	if err != nil {
		log.Printf("cross trading rule: create alert failed: %v", err)
		return
	}

	alertID, err := alertResult.LastInsertId()
	if err != nil {
		log.Printf("cross trading rule: get alert id failed: %v", err)
		return
	}

	for _, t := range selfTrades {
		if err := r.queries.LinkAlertTrade(ctx, db.LinkAlertTradeParams{
			Alertid: uint64(alertID),
			Tradeid: uint64(t.TradeID),
		}); err != nil {
			log.Printf("cross trading rule: link alert trade failed: %v", err)
		}
	}
}
