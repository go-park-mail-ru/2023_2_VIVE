package repository

import (
	"HnH/internal/domain"
	"HnH/internal/repository/mock"
	"HnH/pkg/serverErrors"

	"net/http"
	"time"

	"github.com/google/uuid"
)

type psqlSessionRepository struct {
	connection *mock.Sessions
}

func NewPsqlSessionRepository(conn *mock.Sessions) {

}

func AddSession(user *domain.User) (*http.Cookie, error) {
	uniqueID := uuid.NewString()

	cookie := &http.Cookie{
		Name:     "session",
		Value:    uniqueID,
		Expires:  time.Now().Add(10 * time.Hour),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	userIndex, exist := mock.UserDB.EmailToUser.Load(user.Email)
	if !exist {
		return nil, serverErrors.INVALID_EMAIL
	}

	userToAdd := mock.UserDB.UsersList[userIndex.(int)]

	mock.SessionDB.SessionsList.Store(uniqueID, userToAdd.ID)
	return cookie, nil
}

func DeleteSession(cookie *http.Cookie) error {
	_, exist := mock.SessionDB.SessionsList.Load(cookie.Value)

	if !exist {
		return serverErrors.AUTH_REQUIRED
	}

	mock.SessionDB.SessionsList.Delete(cookie.Value)
	return nil
}

func ValidateSession(cookie *http.Cookie) error {
	_, ok := mock.SessionDB.SessionsList.Load(cookie.Value)

	if !ok {
		return serverErrors.INVALID_COOKIE
	}

	return nil
}
