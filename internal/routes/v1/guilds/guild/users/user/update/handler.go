package update

import (
	"net/http"
	"yourproject/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func Handler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		guildID := c.Param("guildID")
		userID := c.Param("userID")

		var body struct{ Balance int }
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		gu := db.GuildUser{GuildID: guildID, UserID: userID, Balance: body.Balance}
		database.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"balance"}),
		}).Create(&gu)

		c.JSON(http.StatusOK, gu)
	}
}
