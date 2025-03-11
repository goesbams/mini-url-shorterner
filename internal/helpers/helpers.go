package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateShortCode(originalURL string, length int) string {
	hash := md5.Sum([]byte(originalURL))

	hashString := hex.EncodeToString(hash[:])

	shortCode := hashString[:length]

	return shortCode
}
