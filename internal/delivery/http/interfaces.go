package http

import (
	"HnH/internal/domain"
)

type SessionUsecase interface {
	Login(user *domain.User) (string, error)
	Logout(sessionID string) error
	CheckLogin(sessionID string) error
}

type UserUsecase interface {
	SignUp(user *domain.User) (string, error)
	GetInfo(sessionID string) (*domain.User, error)
}

type VacancyUsecase interface {
	GetVacancies() ([]domain.Vacancy, error)
}
