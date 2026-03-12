package auth

import (
	"fmt"

	"git.zyner.org/meta/xyter/internal/database"
	"git.zyner.org/meta/xyter/internal/database/dal"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnsureRootToken checks if there’s any token in the DB; if not, it creates a new root token
func EnsureRootToken(db *gorm.DB) (string, error) {
	var existingTokens int64
	if err := db.Model(&database.Token{}).Count(&existingTokens).Error; err != nil {
		return "", err
	}

	if existingTokens > 0 {
		// tokens already exist, no need to create root
		return "", nil
	}

	// No tokens exist → generate root token
	rootTokenTok, rootToken, err := GenerateToken(db, 0) // 0 = no expiration
	if err != nil {
		return "", err
	}

	// 2️⃣ Assign only token-management permissions
	rootPerms := []struct {
		Endpoint string
		Method   string
	}{
		{"/routes", "GET"},
		{"/api/v1/tokens", "POST"},
		{"/api/v1/tokens/:id", "DELETE"},
		{"/api/v1/tokens/:id/permissions", "POST"},
		{"/api/v1/tokens/:id/permissions/:perm_id", "DELETE"},
		{"/api/v1/tokens/:id/groups", "POST"}, // optional if you have groups
	}
	var perms []database.TokenPermission
	for _, p := range rootPerms {
		perms = append(perms, database.TokenPermission{
			ID:       uuid.NewString(),
			TokenID:  rootToken.ID,
			Endpoint: p.Endpoint,
			Method:   p.Method,
		})
	}

	if err := dal.SetTokenPermissions(db, rootToken.ID, perms); err != nil {
		return "", err
	}
	fmt.Println("Root token created! Store this safely:", rootTokenTok)
	return rootTokenTok, nil
}
