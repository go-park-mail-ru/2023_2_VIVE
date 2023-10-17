package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"

	"net/http"
	"time"
)

func Login(user *domain.User) (*http.Cookie, error) {
	loginErr := repository.CheckUser(user)
	if loginErr != nil {
		return nil, loginErr
	}

	cookie, addErr := repository.AddSession(user)
	if addErr != nil {
		return nil, addErr
	}

	return cookie, nil
}

func Logout(cookie *http.Cookie) error {
	deleteErr := repository.DeleteSession(cookie)
	if deleteErr != nil {
		return deleteErr
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	return nil
}

func CheckLogin(cookie *http.Cookie) error {
	sessionErr := repository.ValidateSession(cookie)
	if sessionErr != nil {
		return sessionErr
	}

	return nil
}
