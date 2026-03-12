package auth

import (
	"errors"
	"time"

	"git.zyner.org/meta/xyter/internal/database/dal"
	"gorm.io/gorm"
)

// CheckPermission checks if a token can access an endpoint with a method
func CheckPermission(db *gorm.DB, tokenHash, endpoint, method string) (bool, error) {
	// 1. Lookup token
	tok, err := dal.GetTokenByHash(db, tokenHash)
	if err != nil || tok.Revoked || tok.ExpiresAt.Before(time.Now()) {
		return false, errors.New("invalid token")
	}

	// 2. Direct token permissions
	tokenPerms, err := dal.GetTokenPermissions(db, tok.ID)
	if err != nil {
		return false, err
	}
	for _, p := range tokenPerms {
		if p.Endpoint == endpoint && p.Method == method {
			return true, nil
		}
	}

	// 3. Group permissions
	groups, err := dal.GetGroupsForToken(db, tok.ID)
	if err != nil {
		return false, err
	}

	for _, g := range groups {
		groupPerms, err := dal.GetPermissionsForGroup(db, g.ID)
		if err != nil {
			return false, err
		}
		for _, gp := range groupPerms {
			if gp.Endpoint == endpoint && gp.Method == method {
				return true, nil
			}
		}
	}

	return false, nil
}
