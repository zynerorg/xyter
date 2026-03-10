package ensure

import (
	"net/http"
	"yourproject/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Handler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var g db.Guild
		if err := c.ShouldBindJSON(&g); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		database.FirstOrCreate(&g, db.Guild{ID: g.ID, Name: g.Name})
		c.JSON(http.StatusOK, g)
	}
}
