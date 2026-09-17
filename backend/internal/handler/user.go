package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type accountSummary struct {
	ID               int32  `json:"id"`
	AccountNumber    int32  `json:"account_number"`
	Balance          uint64 `json:"balance"`
	AvailableBalance uint64 `json:"available_balance"`
	Status           string `json:"status"`
}

type userDetail struct {
	ID           uint32           `json:"id"`
	AverageAsset string           `json:"average_asset"`
	RiskLevel    string           `json:"risk_level"`
	AssetTier    string           `json:"asset_tier"`
	Accounts     []accountSummary `json:"accounts"`
}

// ListUsers godoc
// GET /api/v1/users
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.Queries.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// GET /api/v1/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.Queries.GetUser(c.Request.Context(), uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accounts, err := h.SourceQueries.ListAccountsByUser(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	summaries := make([]accountSummary, len(accounts))
	for i, a := range accounts {
		summaries[i] = accountSummary{
			ID:               a.ID,
			AccountNumber:    a.AccountNumber,
			Balance:          a.Balance,
			AvailableBalance: a.AvailableBalance,
			Status:           string(a.Status),
		}
	}

	c.JSON(http.StatusOK, userDetail{
		ID:           user.ID,
		AverageAsset: user.AverageAsset,
		RiskLevel:    string(user.RiskLevel),
		AssetTier:    string(user.AssetTier),
		Accounts:     summaries,
	})
}
