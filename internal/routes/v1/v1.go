package routes

import (
	"git.zyner.org/meta/xyter/internal/routes/v1/guilds"
	"git.zyner.org/meta/xyter/internal/routes/v1/tokens"
	"git.zyner.org/meta/xyter/internal/routes/v1/users"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	guilds.Register(v1)
	users.Register(v1)
	tokens.Register(v1)
}
