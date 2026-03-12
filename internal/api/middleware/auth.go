package middleware

import (
	"net/http"
	"strings"

	"git.zyner.org/meta/xyter/internal/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TokenMiddleware checks if a token has permission for the endpoint+method
func TokenMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" {
			c.Next()
		}
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		println(c.FullPath())
		token := strings.TrimPrefix(header, "Bearer ")
		hash := auth.HashToken(token)
		allowed, err := auth.CheckPermission(db, hash, c.FullPath(), c.Request.Method)
		if err != nil || !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		c.Next()
	}
}
