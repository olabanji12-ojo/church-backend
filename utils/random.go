package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString generates a secure random string of a given length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
