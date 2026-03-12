package dal

import (
	"git.zyner.org/meta/xyter/internal/database"
	"gorm.io/gorm"
)

// SetTokenGroups assigns a token to multiple groups
func SetTokenGroups(db *gorm.DB, tokenID string, groupIDs []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("token_id = ?", tokenID).Delete(&database.TokenGroup{}).Error; err != nil {
			return err
		}
		for _, gid := range groupIDs {
			if err := tx.Create(&database.TokenGroup{
				TokenID: tokenID,
				GroupID: gid,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
