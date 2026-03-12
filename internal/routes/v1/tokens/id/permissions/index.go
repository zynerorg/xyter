package permissions

import (
	"net/http"

	"git.zyner.org/meta/xyter/internal/api/middleware"
	"git.zyner.org/meta/xyter/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// POST /api/tokens/:id/permissions
func PostHandler(c *gin.Context) {
	db := middleware.GetDB(c)
	tokenID := c.Param("id")

	var req struct {
		Endpoint string `json:"endpoint"`
		Method   string `json:"method"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	perm := database.TokenPermission{
		ID:       uuid.NewString(),
		TokenID:  tokenID,
		Endpoint: req.Endpoint,
		Method:   req.Method,
	}

	if err := db.Create(&perm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add permission"})
		return
	}

	c.JSON(http.StatusOK, perm)
}

// DELETE /api/tokens/:id/permissions/:perm_id
func DeleteHandler(c *gin.Context) {
	db := middleware.GetDB(c)
	permID := c.Param("perm_id")

	if err := db.Delete(&database.TokenPermission{}, "id = ?", permID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Register routes
func Register(r *gin.RouterGroup) {
	permissionsGroup := r.Group("/permissions")
	permissionsGroup.POST("", PostHandler)
	permissionsGroup.DELETE("/:perm_id", DeleteHandler)
}
