package audit

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/pkg/audit"
)

type Handler struct {
	auditLogger *audit.Logger
}

func NewHandler(auditLogger *audit.Logger) *Handler {
	return &Handler{
		auditLogger: auditLogger,
	}
}

func (h *Handler) GetLogs(c *gin.Context) {
	userID := c.Query("user_id")
	action := c.Query("action")
	resourceType := c.Query("resource_type")

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	logs, err := h.auditLogger.GetLogs(c.Request.Context(), userID, action, resourceType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get audit logs: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"pagination": gin.H{
			"count": len(logs),
			"limit": limit,
			"offset": offset,
		},
	})
}
