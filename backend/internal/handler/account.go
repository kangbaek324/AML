package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type stockHolding struct {
	StockID   int32  `json:"stock_id"`
	Name      string `json:"name"`
	Quantity  uint64 `json:"quantity"`
	Average   uint64 `json:"average"`
	Price     uint64 `json:"price"`
	Valuation uint64 `json:"valuation"`
}

type accountDetail struct {
	ID               int32          `json:"id"`
	UserID           int32          `json:"user_id"`
	AccountNumber    int32          `json:"account_number"`
	Balance          uint64         `json:"balance"`
	AvailableBalance uint64         `json:"available_balance"`
	Status           string         `json:"status"`
	Stocks           []stockHolding `json:"stocks"`
}

// GetAccount godoc
// GET /api/v1/accounts/:id
func (h *Handler) GetAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	account, err := h.SourceQueries.GetAccount(c.Request.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	holdings, err := h.SourceQueries.ListStockHoldingsByAccount(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stocks := make([]stockHolding, len(holdings))
	for i, s := range holdings {
		stocks[i] = stockHolding{
			StockID:   s.StockID,
			Name:      s.Name,
			Quantity:  s.Quantity,
			Average:   s.Average,
			Price:     s.Price,
			Valuation: s.Quantity * s.Price,
		}
	}

	c.JSON(http.StatusOK, accountDetail{
		ID:               account.ID,
		UserID:           account.UserID,
		AccountNumber:    account.AccountNumber,
		Balance:          account.Balance,
		AvailableBalance: account.AvailableBalance,
		Status:           string(account.Status),
		Stocks:           stocks,
	})
}
