package id

import (
	"net/http"

	"git.zyner.org/meta/xyter/internal/api/middleware"
	"git.zyner.org/meta/xyter/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	var g database.Guild
	if err := db.First(&g, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "guild not found"})
		return
	}
	c.JSON(http.StatusOK, g)

}

func PutHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	var g database.Guild
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&g)
	c.JSON(http.StatusOK, g)
}

func Register(r *gin.RouterGroup) {
	r.GET("/:id", GetHandler)
	r.PUT("/:id", PutHandler)
}
