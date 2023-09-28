package modelHandlers

import (
	"models/errors"
	"models/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func AddSession(user models.User) http.Cookie {
	uniqueID := uuid.NewString()

	cookie := http.Cookie{
		Name:     "session",
		Value:    uniqueID,
		Expires:  time.Now().Add(10 * time.Hour),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	models.Sessions.Store(uniqueID, user)
	return cookie
}

func DeleteSession(cookie *http.Cookie) error {
	_, exist := models.Sessions.Load(cookie.Value)

	if exist {
		return errors.SESSION_ALREADY_EXISTS
	}

	models.Sessions.Delete(cookie.Value)
	return nil
}

func ValidateSession(cookie *http.Cookie) error {
	if time.Now().After(cookie.Expires) {
		return errors.COOKIE_HAS_EXPIRED
	}

	storedID, ok := models.Sessions.Load(cookie.Value)

	if !ok {
		return errors.AUTH_REQUIRED
	}

	if cookie.Value != storedID {
		return errors.INVALID_COOKIE
	}

	return nil
}
