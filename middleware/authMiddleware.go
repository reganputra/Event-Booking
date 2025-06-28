package middleware

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/utils"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			return
		}

		// If the token is in the format "Bearer <token>", remove the "Bearer" prefix
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Validate the token
		userId, role, err := utils.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set the user ID and role in the context for use in handlers
		c.Set("userId", userId)
		c.Set("userRole", role)

		// Continue to the next handler
		c.Next()
	}
}
