package api

import (
	"net/http"

	"git.zyner.org/meta/xyter/internal/api/middleware"
	routes "git.zyner.org/meta/xyter/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
)

// Start initializes the API, registers routers etc
func Start(k *koanf.Koanf, db *gorm.DB) {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Register DB middleware
	r.Use(middleware.InjectDB(db))

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	routes.RegisterRoutes(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
