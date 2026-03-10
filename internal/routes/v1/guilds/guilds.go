package guilds

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Handler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost:
			var g db.Guild
			if err := c.ShouldBindJSON(&g); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
				return
			}
			database.FirstOrCreate(&g, db.Guild{ID: g.ID, Name: g.Name})
			c.JSON(http.StatusOK, g)
			return

		case http.MethodGet:
			id := c.Query("id")
			if id == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
				return
			}
			var g db.Guild
			if err := database.First(&g, "id = ?", id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "guild not found"})
				return
			}
			c.JSON(http.StatusOK, g)
			return

		case http.MethodPut:
			var g db.Guild
			if err := c.ShouldBindJSON(&g); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
				return
			}
			database.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"name"}),
			}).Create(&g)
			c.JSON(http.StatusOK, g)
			return
		}
	}
}
