package routes

import (
	v1 "git.zyner.org/meta/xyter/internal/routes/v1"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1.Register(api)
}
