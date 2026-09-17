package worker

import (
	"context"
	"log"
	"time"

	db "github.com/kangbaek324/AML/db/sqlc"
)

const (
	// 유효한(ABNORMAL) Alert 발생 횟수 기준. docs 참고 - LOW 0~2회 / MEDIUM 3~5회 / HIGH 6회 이상
	userRiskLowMax    = 2
	userRiskMediumMax = 5
)

var userRiskRank = map[db.UsersRiskLevel]int{
	db.UsersRiskLevelLOW:    0,
	db.UsersRiskLevelMEDIUM: 1,
	db.UsersRiskLevelHIGH:   2,
}

var userRiskByRank = []db.UsersRiskLevel{
	db.UsersRiskLevelLOW,
	db.UsersRiskLevelMEDIUM,
	db.UsersRiskLevelHIGH,
}

// UserRiskWorker는 지난 한 달간의 유효 Alert 발생 횟수를 집계해 유저의 risk_level을
// 갱신한다. 등급 상승은 즉시 반영하고, 하락은 한 단계씩만 진행한다.
type UserRiskWorker struct {
	queries *db.Queries
}

func NewUserRiskWorker(queries *db.Queries) *UserRiskWorker {
	return &UserRiskWorker{queries: queries}
}

// Start는 부팅 시 조건에 맞춰 즉시 한 번 실행하고, 이후에는 매달 1일 00:00(UTC)
// 트리거로 재실행한다.
func (w *UserRiskWorker) Start(ctx context.Context) {
	w.run(ctx)

	for {
		timer := time.NewTimer(time.Until(nextMonthStart(time.Now().UTC())))

		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			w.run(ctx)
		}
	}
}

func (w *UserRiskWorker) run(ctx context.Context) {
	now := time.Now().UTC()
	periodEnd := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodStart := periodEnd.AddDate(0, -1, 0)

	counts, err := w.queries.CountValidAlertsByUser(ctx, db.CountValidAlertsByUserParams{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	})
	if err != nil {
		log.Printf("user risk worker: count alerts failed: %v", err)
		return
	}

	alertCount := make(map[uint32]int64, len(counts))
	for _, c := range counts {
		alertCount[c.Userid] = c.Count
	}

	users, err := w.queries.ListUserRiskLevels(ctx)
	if err != nil {
		log.Printf("user risk worker: list users failed: %v", err)
		return
	}

	for _, u := range users {
		computed := classifyUserRisk(alertCount[u.ID])
		next := nextUserRisk(u.RiskLevel, computed)
		if next == u.RiskLevel {
			continue
		}

		if err := w.queries.UpdateUserRiskLevel(ctx, db.UpdateUserRiskLevelParams{
			ID:        u.ID,
			RiskLevel: next,
		}); err != nil {
			log.Printf("user risk worker: update user %d failed: %v", u.ID, err)
		}
	}
}

func classifyUserRisk(alertCount int64) db.UsersRiskLevel {
	switch {
	case alertCount > userRiskMediumMax:
		return db.UsersRiskLevelHIGH
	case alertCount > userRiskLowMax:
		return db.UsersRiskLevelMEDIUM
	default:
		return db.UsersRiskLevelLOW
	}
}

// nextUserRisk는 등급 상승은 즉시, 하락은 한 단계씩만 반영한다.
func nextUserRisk(current, computed db.UsersRiskLevel) db.UsersRiskLevel {
	currentRank := userRiskRank[current]
	computedRank := userRiskRank[computed]

	switch {
	case computedRank > currentRank:
		return computed
	case computedRank < currentRank:
		return userRiskByRank[currentRank-1]
	default:
		return current
	}
}

func nextMonthStart(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
}
