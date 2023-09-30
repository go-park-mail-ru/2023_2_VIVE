package modelHandlers

import (
	"models/models"
	"models/serverErrors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func AddSession(user *models.User) http.Cookie {
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

	if !exist {
		return serverErrors.AUTH_REQUIRED
	}

	models.Sessions.Delete(cookie.Value)
	return nil
}

func ValidateSession(cookie *http.Cookie) error {
	if time.Now().After(cookie.Expires) {
		return serverErrors.COOKIE_HAS_EXPIRED
	}

	storedID, ok := models.Sessions.Load(cookie.Value)

	if !ok {
		return serverErrors.AUTH_REQUIRED
	}

	if cookie.Value != storedID {
		return serverErrors.INVALID_COOKIE
	}

	return nil
}
