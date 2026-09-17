package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefreshUsers godoc
// POST /api/v1/users/refresh
// Asset Tier / User Risk 워커를 스케줄과 무관하게 즉시 실행한다.
func (h *Handler) RefreshUsers(c *gin.Context) {
	h.Worker.RunNow(c.Request.Context())
	c.Status(http.StatusNoContent)
}
