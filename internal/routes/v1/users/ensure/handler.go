package ensure

import (
	"net/http"
	"yourproject/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Handler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u db.User
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		database.FirstOrCreate(&u, db.User{ID: u.ID, Username: u.Username})
		c.JSON(http.StatusOK, u)
	}
}
