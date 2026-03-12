package guilds

import (
	"net/http"

	"git.zyner.org/meta/xyter/internal/api/middleware"
	"git.zyner.org/meta/xyter/internal/database"
	"git.zyner.org/meta/xyter/internal/routes/v1/guilds/id"
	"github.com/gin-gonic/gin"
)

func PostHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	var g database.Guild
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	db.FirstOrCreate(&g, database.Guild{ID: g.ID})
	c.JSON(http.StatusOK, g)

}

func GetHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	var guilds []database.Guild
	db.Find(&guilds)
	c.JSON(http.StatusOK, guilds)
}
func Register(r *gin.RouterGroup) {
	guildsGroup := r.Group("/guilds")
	guildsGroup.POST("", PostHandler)
	guildsGroup.GET("", GetHandler)
	id.Register(r)
}
