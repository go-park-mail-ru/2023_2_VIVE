package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"

	"net/http"
)

func SignUp(user *domain.User) (*http.Cookie, error) {
	addStatus := repository.AddUser(user)
	if addStatus != nil {
		return nil, addStatus
	}

	cookie, addErr := repository.AddSession(user)
	if addErr != nil {
		return nil, addErr
	}

	return cookie, nil
}

func GetInfo(cookie *http.Cookie) (*domain.User, error) {
	validStatus := repository.ValidateSession(cookie)
	if validStatus != nil {
		return nil, validStatus
	}

	user, getErr := repository.GetUserInfo(cookie)
	if getErr != nil {
		return nil, getErr
	}

	return user, nil
}
