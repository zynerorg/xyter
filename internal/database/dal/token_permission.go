package dal

import (
	"git.zyner.org/meta/xyter/internal/database"
	"gorm.io/gorm"
)

// GetTokenPermissions returns all direct permissions of a token
func GetTokenPermissions(db *gorm.DB, tokenID string) ([]database.TokenPermission, error) {
	var perms []database.TokenPermission
	if err := db.Where("token_id = ?", tokenID).Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// SetTokenPermissions sets direct permissions for a token
func SetTokenPermissions(db *gorm.DB, tokenID string, perms []database.TokenPermission) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("token_id = ?", tokenID).Delete(&database.TokenPermission{}).Error; err != nil {
			return err
		}
		for _, p := range perms {
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
