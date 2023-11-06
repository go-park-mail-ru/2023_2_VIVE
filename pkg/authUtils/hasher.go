package authUtils

import (
	"crypto/rand"
	"crypto/subtle"

	"golang.org/x/crypto/argon2"
)

type hashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var p = hashParams{
	memory:      64 * 1024,
	iterations:  1,
	parallelism: 4,
	saltLength:  8,
	keyLength:   32,
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, p.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func GenerateHash(password string) (hash []byte, salt []byte, err error) {
	salt, err = generateSalt()
	if err != nil {
		return nil, nil, err
	}

	hash = argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	return hash, salt, nil
}

func ComparePasswordAndHash(password string, salt, hashedPass []byte) bool {
	hashToCheck := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hashedPass, hashToCheck) == 1 {
		return true
	}

	return false
}
