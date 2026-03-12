package tokens

import (
	"net/http"
	"time"

	"git.zyner.org/meta/xyter/internal/api/middleware"
	"git.zyner.org/meta/xyter/internal/auth"
	"git.zyner.org/meta/xyter/internal/database"
	"git.zyner.org/meta/xyter/internal/routes/v1/guilds/id"
	"github.com/gin-gonic/gin"
)

type CreateTokenRequest struct {
	TTLSeconds int64 `json:"ttl_seconds"` // optional
}

// POST /api/tokens
func PostHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	ttl := time.Duration(req.TTLSeconds) * time.Second
	tokenStr, tok, err := auth.GenerateToken(db, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenStr,
		"id":         tok.ID,
		"expires_at": tok.ExpiresAt,
	})
}

// GET /api/tokens
func GetHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	var tokens []database.Token
	db.Find(&tokens)
	c.JSON(http.StatusOK, tokens)
}

// Register routes
func Register(r *gin.RouterGroup) {
	tokens := r.Group("/tokens")
	tokens.POST("", PostHandler)
	tokens.GET("", GetHandler)
	id.Register(tokens)
}
