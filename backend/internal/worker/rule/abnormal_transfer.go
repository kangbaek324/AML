package rule

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
)

const abnormalTransferInterval = 10 * time.Second

const abnormalTransferThreshold = 10_000_000

// AbnormalTransferRule은 완료된 송금 중 받은 금액이 임계값 이상인 단일 송금과,
// 정시(1시간) 단위로 버킷팅해 합산한 금액이 임계값을 초과하는 누적 송금을 탐지한다.
type AbnormalTransferRule struct {
	queries       *db.Queries
	sourceQueries *sourcedb.Queries
}

func NewAbnormalTransferRule(queries *db.Queries, sourceQueries *sourcedb.Queries) *AbnormalTransferRule {
	return &AbnormalTransferRule{queries: queries, sourceQueries: sourceQueries}
}

func (r *AbnormalTransferRule) Name() string { return "abnormal_transfer" }

func (r *AbnormalTransferRule) Interval() time.Duration { return abnormalTransferInterval }

func (r *AbnormalTransferRule) Run(ctx context.Context, since time.Time) (time.Time, error) {
	transfers, err := r.sourceQueries.ListTransfersSince(ctx, sql.NullTime{Time: since, Valid: true})
	if err != nil {
		return since, fmt.Errorf("list transfers: %w", err)
	}
	if len(transfers) == 0 {
		return since, nil
	}

	for _, t := range transfers {
		r.checkSingle(ctx, uint32(t.RecipientUserID), t.Amount, t.TransferID)
		r.checkCumulative(ctx, uint32(t.RecipientUserID), t.Amount, t.CompletedAt.Time)
	}

	// transfers는 completed_at 오름차순이므로 마지막 항목이 가장 최근 완료 송금이다.
	return transfers[len(transfers)-1].CompletedAt.Time, nil
}

func (r *AbnormalTransferRule) checkSingle(ctx context.Context, recipientUserID uint32, amount uint64, transferID int64) {
	if amount < abnormalTransferThreshold {
		return
	}

	reason := fmt.Sprintf("단일 송금 수신액 %d원이 임계값(%d원) 이상입니다.", amount, abnormalTransferThreshold)
	r.createAlert(ctx, recipientUserID, reason, []int64{transferID})
}

func (r *AbnormalTransferRule) checkCumulative(ctx context.Context, recipientUserID uint32, amount uint64, completedAt time.Time) {
	hour := completedAt.Truncate(time.Hour)

	if err := r.queries.UpsertTransferHourSum(ctx, db.UpsertTransferHourSumParams{
		Userid: recipientUserID,
		Hour:   hour,
		Amount: amount,
	}); err != nil {
		log.Printf("abnormal transfer rule: upsert hour sum failed: %v", err)
		return
	}

	sum, err := r.queries.GetTransferHourSum(ctx, db.GetTransferHourSumParams{
		Userid: recipientUserID,
		Hour:   hour,
	})
	if err != nil {
		log.Printf("abnormal transfer rule: get hour sum failed: %v", err)
		return
	}
	if sum <= abnormalTransferThreshold {
		return
	}

	claimed, err := r.queries.ClaimTransferAlert(ctx, db.ClaimTransferAlertParams{
		Userid: recipientUserID,
		Hour:   hour,
	})
	if err != nil {
		log.Printf("abnormal transfer rule: claim alert failed: %v", err)
		return
	}
	rows, err := claimed.RowsAffected()
	if err != nil {
		log.Printf("abnormal transfer rule: claim alert rows affected failed: %v", err)
		return
	}
	if rows == 0 {
		return
	}

	transfers, err := r.sourceQueries.ListReceivedTransfersByUserAndHour(ctx, sourcedb.ListReceivedTransfersByUserAndHourParams{
		UserID:    int32(recipientUserID),
		HourStart: sql.NullTime{Time: hour, Valid: true},
		HourEnd:   sql.NullTime{Time: hour.Add(time.Hour), Valid: true},
	})
	if err != nil {
		log.Printf("abnormal transfer rule: list hour transfers failed: %v", err)
		return
	}

	transferIDs := make([]int64, 0, len(transfers))
	for _, t := range transfers {
		transferIDs = append(transferIDs, t.TransferID)
	}

	reason := fmt.Sprintf(
		"%s 1시간 누적 송금 수신액이 %d원으로 임계값(%d원)을 초과했습니다.",
		hour.Format("2006-01-02 15:00"), sum, abnormalTransferThreshold,
	)
	r.createAlert(ctx, recipientUserID, reason, transferIDs)
}

func (r *AbnormalTransferRule) createAlert(ctx context.Context, userID uint32, reason string, transferIDs []int64) {
	result, err := r.queries.CreateAlert(ctx, db.CreateAlertParams{
		Userid: userID,
		Type:   db.AlertsTypeABNORMALTRANSFER,
		Reason: reason,
	})
	if err != nil {
		log.Printf("abnormal transfer rule: create alert failed: %v", err)
		return
	}

	alertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("abnormal transfer rule: get alert id failed: %v", err)
		return
	}

	for _, transferID := range transferIDs {
		if err := r.queries.LinkAlertTransfer(ctx, db.LinkAlertTransferParams{
			Alertid:    uint64(alertID),
			Transferid: uint64(transferID),
		}); err != nil {
			log.Printf("abnormal transfer rule: link alert transfer failed: %v", err)
		}
	}
}
