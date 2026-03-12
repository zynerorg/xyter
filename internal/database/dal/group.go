package dal

import (
	"git.zyner.org/meta/xyter/internal/database"
	"gorm.io/gorm"
)

// GetGroupsForToken returns all groups a token belongs to
func GetGroupsForToken(db *gorm.DB, tokenID string) ([]database.Group, error) {
	var groups []database.Group
	err := db.Table("groups").
		Joins("JOIN token_groups tg ON tg.group_id = groups.id").
		Where("tg.token_id = ?", tokenID).
		Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// InsertGroup inserts a new group
func InsertGroup(db *gorm.DB, group database.Group) error {
	return db.Create(&group).Error
}
