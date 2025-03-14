package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/pkg/audit"
)

func AuditMiddleware(auditLogger *audit.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		path := c.Request.URL.Path
		if path == "/api/v1/health" {
			return
		}

		action := c.Request.Method

		resourceType := "api"
		if len(path) > 8 && path[:8] == "/api/v1/" {
			parts := []rune(path[8:])
			for i, char := range parts {
				if char == '/' {
					resourceType = string(parts[:i])
					break
				}
			}

			if resourceType == path[8:] {
				resourceType = path[8:]
			}
		}

		resourceID := c.Param("id")
		
		details := map[string]interface{}{
			"path": path,
			"method": action,
			"status": c.Writer.Status(),
		}

		if err := auditLogger.LogFromGin(c, action, resourceType, resourceID, details); err != nil {
			c.Error(err)
		}
	}
}