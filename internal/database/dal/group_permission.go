package dal

import (
	"git.zyner.org/meta/xyter/internal/database"
	"gorm.io/gorm"
)

// GetPermissionsForGroup returns all permissions of a group
func GetPermissionsForGroup(db *gorm.DB, groupID string) ([]database.GroupPermission, error) {
	var perms []database.GroupPermission
	err := db.Where("group_id = ?", groupID).Find(&perms).Error
	return perms, err
}

// SetPermissionsForGroup sets permissions for a group
func SetPermissionsForGroup(db *gorm.DB, groupID string, perms []database.GroupPermission) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", groupID).Delete(&database.GroupPermission{}).Error; err != nil {
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
