package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthorizeRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User role not found in context"})
			return
		}

		if userRole.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
			return
		}
		c.Next()
	}
}
