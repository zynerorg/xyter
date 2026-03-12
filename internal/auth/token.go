package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"git.zyner.org/meta/xyter/internal/database"
	"git.zyner.org/meta/xyter/internal/database/dal"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateToken(db *gorm.DB, ttl time.Duration) (string, *database.Token, error) {
	b := make([]byte, 32)
	rand.Read(b)
	tokenStr := "tk_" + hex.EncodeToString(b)
	hash := HashToken(tokenStr)
	expires := time.Now().Add(ttl)
	if ttl == 0 {
		expires = time.Now().AddDate(100, 0, 0) // effectively no expiration
	}
	tok := &database.Token{
		ID:        uuid.NewString(),
		Hash:      hash,
		ExpiresAt: expires,
		Revoked:   false,
	}

	if err := dal.InsertToken(db, *tok); err != nil {
		return "", nil, err
	}

	return tokenStr, tok, nil
}
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
