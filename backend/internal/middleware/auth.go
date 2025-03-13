package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/pkg/jwt"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header from the request
		authHeader := c.GetHeader("Authorization")

		// Check if the header is present and starts with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required and must be a Bearer token",
			})
			c.Abort()
			return
		}

		// Extract the token string from the header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Set user information in the context
		// This makes user data available to handlers
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks if the user has the required role
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user role from the context (set by AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authenticated",
			})
			c.Abort()
			return
		}

		// Check if the user's role is in the list of required roles
		userRole := role.(string)
		hasRequiredRole := false
		
		for _, r := range requiredRoles {
			if userRole == r {
				hasRequiredRole = true
				break
			}
		}

		// If the user doesn't have any of the required roles, return forbidden
		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}