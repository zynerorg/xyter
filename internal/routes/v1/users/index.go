package users

import (
	"net/http"

	"git.zyner.org/meta/xyter/internal/api/middleware"
	"git.zyner.org/meta/xyter/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func PostHandler(c *gin.Context) {
	db := middleware.GetDB(c)
	var u database.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	db.FirstOrCreate(&u, database.User{ID: u.ID})
	c.JSON(http.StatusOK, u)
	return

}
func GetHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	id := c.Query("id")
	var u database.User
	if err := db.First(&u, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
	return
}
func PutHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	var body struct {
		UserID  string `json:"user_id"`
		GuildID string `json:"guild_id"`
		Balance int    `json:"balance"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	gu := database.GuildUser{
		GuildID: body.GuildID,
		UserID:  body.UserID,
		Balance: body.Balance,
	}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "guild_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"balance"}),
	}).Create(&gu)
	c.JSON(http.StatusOK, gu)
	return

}

func Register(r *gin.RouterGroup) {
	usersGroup := r.Group("/users")
	usersGroup.POST("", PostHandler)
	usersGroup.GET("/:id", GetHandler)
	usersGroup.PUT(":id", PutHandler)
}
