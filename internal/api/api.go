package api

import (
	"log"
	"net/http"

	"git.zyner.org/meta/xyter/internal/api/middleware"
	"git.zyner.org/meta/xyter/internal/auth"
	routes "git.zyner.org/meta/xyter/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

// List Routes for clients
func ListRoutesHandler(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		routes := r.Routes()
		out := make([]gin.H, 0, len(routes))
		for _, rt := range routes {
			out = append(out, gin.H{
				"method":  rt.Method,
				"path":    rt.Path,
				"handler": rt.Handler,
			})
		}
		c.JSON(http.StatusOK, out)
	}
}

// Start initializes the API, registers routers etc
func Start(k *koanf.Koanf, db *gorm.DB) {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Ensure root token exists
	_, err := auth.EnsureRootToken(db)
	if err != nil {
		log.Fatal("Failed to create initial token:", err)
	}

	// Register DB middleware
	r.Use(middleware.InjectDB(db))
	// Register permission middleware
	r.Use(middleware.TokenMiddleware(db))

	r.GET("/routes", ListRoutesHandler(r))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	routes.RegisterRoutes(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
