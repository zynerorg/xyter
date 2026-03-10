package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const DBKey = "db"

func InjectDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DBKey, db)
		c.Next()
	}
}

func GetDB(c *gin.Context) *gorm.DB {
	db, exists := c.Get(DBKey)
	if !exists {
		panic("database not found in gin context")
	}

	return db.(*gorm.DB)
}
