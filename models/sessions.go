package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type cookieData struct {
	Name    string
	Value   string
	Expires time.Time
}

func getNewRecord(cookie *http.Cookie) cookieData {
	return cookieData{
		Name:    cookie.Name,
		Value:   cookie.Value,
		Expires: cookie.Expires,
	}
}

var sessions = make(map[cookieData][]byte)

func AddSession(sessionName string) http.Cookie {
	hasher := sha256.New()
	hasher.Write([]byte(sessionName))
	hash := hasher.Sum(nil)

	cookie := http.Cookie{
		Name:     "authCookie",
		Value:    hex.EncodeToString(hash),
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}

	sessions[getNewRecord(&cookie)] = hash

	return cookie
}

func ValidateSession(cookie *http.Cookie) error {
	storedHash, ok := sessions[getNewRecord(cookie)]

	if !ok {
		return fmt.Errorf("AUTHENTICATION_REQUIRED")
	}

	hasher := sha256.New()
	hasher.Write([]byte(cookie.Value))
	hashToCheck := hasher.Sum(nil)

	if hex.EncodeToString(hashToCheck) != hex.EncodeToString(storedHash) {
		return fmt.Errorf("INVALID_COOKIE")
	}

	return nil
}
