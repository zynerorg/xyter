package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SaveToken saves a token to disk, optionally encrypting it with a password.
func SaveToken(token, password string, encrypt bool) error {
	if encrypt {
		var err error
		token, err = encryptToken(token, password)
		if err != nil {
			return fmt.Errorf("encrypting token: %w", err)
		}
	}

	path := filepath.Join(os.Getenv("HOME"), ".myapp_token")
	return os.WriteFile(path, []byte(token), 0600)
}

// encryptToken encrypts a string with AES-GCM using a password-derived key.
func encryptToken(token, password string) (string, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("creating cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("creating GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generating nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
