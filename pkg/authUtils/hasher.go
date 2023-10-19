package authUtils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHash(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPass := hasher.Sum(nil)

	return hex.EncodeToString(hashedPass)
}
