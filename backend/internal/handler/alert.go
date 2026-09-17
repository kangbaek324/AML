package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	db "github.com/kangbaek324/AML/db/sqlc"
)

type alertSummary struct {
	ID          uint64     `json:"id"`
	UserID      uint32     `json:"user_id"`
	Type        string     `json:"type"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	AlertedAt   time.Time  `json:"alerted_at"`
	ProcessedAt *time.Time `json:"processed_at"`
}

type tradeSummary struct {
	TradeID     int64     `json:"trade_id"`
	StockID     int32     `json:"stock_id"`
	Quantity    uint64    `json:"quantity"`
	Price       uint64    `json:"price"`
	MatchedAt   time.Time `json:"matched_at"`
	MakerUserID int32     `json:"maker_user_id"`
	TakerUserID int32     `json:"taker_user_id"`
}

type transferSummary struct {
	TransferID      int64      `json:"transfer_id"`
	SenderUserID    int32      `json:"sender_user_id"`
	RecipientUserID int32      `json:"recipient_user_id"`
	Amount          uint64     `json:"amount"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

type alertDetail struct {
	alertSummary
	Trades    []tradeSummary    `json:"trades"`
	Transfers []transferSummary `json:"transfers"`
}

// ListAlerts godoc
// GET /api/v1/alerts
func (h *Handler) ListAlerts(c *gin.Context) {
	alerts, err := h.Queries.ListAlerts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	summaries := make([]alertSummary, len(alerts))
	for i, a := range alerts {
		summaries[i] = toAlertSummary(a)
	}

	c.JSON(http.StatusOK, summaries)
}

// GetAlert godoc
// GET /api/v1/alerts/:id
func (h *Handler) GetAlert(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert id"})
		return
	}

	alert, err := h.Queries.GetAlert(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detail := alertDetail{
		alertSummary: toAlertSummary(alert),
		Trades:       []tradeSummary{},
		Transfers:    []transferSummary{},
	}

	if alert.Type == db.AlertsTypeABNORMALTRANSFER {
		transferIDs, err := h.Queries.ListAlertTransferIDs(ctx, alert.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, transferID := range transferIDs {
			t, err := h.SourceQueries.GetTransferDetail(ctx, int64(transferID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			s := transferSummary{
				TransferID:      t.TransferID,
				SenderUserID:    t.SenderUserID,
				RecipientUserID: t.RecipientUserID,
				Amount:          t.Amount,
				Status:          string(t.Status),
				CreatedAt:       t.CreatedAt,
			}
			if t.CompletedAt.Valid {
				s.CompletedAt = &t.CompletedAt.Time
			}
			detail.Transfers = append(detail.Transfers, s)
		}
	} else {
		tradeIDs, err := h.Queries.ListAlertTradeIDs(ctx, alert.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, tradeID := range tradeIDs {
			t, err := h.SourceQueries.GetTradeDetail(ctx, int64(tradeID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			detail.Trades = append(detail.Trades, tradeSummary{
				TradeID:     t.TradeID,
				StockID:     t.StockID,
				Quantity:    t.Quantity,
				Price:       t.Price,
				MatchedAt:   t.MatchedAt,
				MakerUserID: t.MakerUserID,
				TakerUserID: t.TakerUserID,
			})
		}
	}

	c.JSON(http.StatusOK, detail)
}

type updateAlertStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=NORMAL ABNORMAL"`
}

// UpdateAlertStatus godoc
// PATCH /api/v1/alerts/:id/status
// PENDING 상태인 Alert만 NORMAL/ABNORMAL로 처리할 수 있다.
func (h *Handler) UpdateAlertStatus(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert id"})
		return
	}

	var req updateAlertStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.Queries.UpdateAlertStatus(ctx, db.UpdateAlertStatusParams{
		Status: db.AlertsStatus(req.Status),
		ID:     id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if affected == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "alert not found or already processed"})
		return
	}

	c.Status(http.StatusNoContent)
}

func toAlertSummary(a db.Alert) alertSummary {
	s := alertSummary{
		ID:        a.ID,
		UserID:    a.Userid,
		Type:      string(a.Type),
		Reason:    a.Reason,
		Status:    string(a.Status),
		AlertedAt: a.AlertedAt,
	}
	if a.ProcessedAt.Valid {
		s.ProcessedAt = &a.ProcessedAt.Time
	}
	return s
}
