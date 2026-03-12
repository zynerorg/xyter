package dal

import (
	"errors"
	"git.zyner.org/meta/xyter/internal/database"
	"gorm.io/gorm"
)

// InsertToken inserts a new token
func InsertToken(db *gorm.DB, token database.Token) error {
	return db.Create(&token).Error
}

// GetTokenByHash returns a token by its hash
func GetTokenByHash(db *gorm.DB, hash string) (*database.Token, error) {
	var t database.Token
	err := db.First(&t, "hash = ?", hash).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("token not found")
	}
	return &t, err
}

// RevokeToken marks a token as revoked
func RevokeToken(db *gorm.DB, hash string) error {
	return db.Model(&database.Token{}).
		Where("hash = ?", hash).
		Update("revoked", true).Error
}
