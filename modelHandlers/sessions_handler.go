package modelHandlers

import (
	"HnH/models"
	"HnH/serverErrors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func AddSession(user *models.User) *http.Cookie {
	uniqueID := uuid.NewString()

	cookie := &http.Cookie{
		Name:     "session",
		Value:    uniqueID,
		Expires:  time.Now().Add(10 * time.Hour),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	userIndex, _ := models.EmailToUser.Load(user.Email)
	userToAdd := models.UserDB.UsersList[userIndex.(int)]

	models.Sessions.Store(uniqueID, userToAdd.ID)
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
	_, ok := models.Sessions.Load(cookie.Value)

	if !ok {
		return serverErrors.INVALID_COOKIE
	}

	return nil
}
