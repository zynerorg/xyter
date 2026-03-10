package get

import (
	"net/http"
	"yourproject/db"

	"github.com/gin-gonic/gin"
)

func Handler(database *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		guildID := c.Param("guildID")
		userID := c.Param("userID")

		var gu db.GuildUser
		if err := database.First(&gu, "guild_id = ? AND user_id = ?", guildID, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		c.JSON(http.StatusOK, gu)
	}
}
