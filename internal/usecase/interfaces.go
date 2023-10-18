package usecase

import (
	"HnH/internal/domain"
)

type SessionRepository interface {
	AddSession(sessionID string, userID int) error
	DeleteSession(sessionID string) error
	ValidateSession(sessionID string) error
	GetUserIdBySession(sessionID string) (int, error)
}

type UserRepository interface {
	CheckUser(user *domain.User) error
	AddUser(user *domain.User) error
	GetUserInfo(userID int) (*domain.User, error)
	GetUserIdByEmail(email string) (int, error)
}

type VacancyRepository interface {
	GetVacancies() ([]domain.Vacancy, error)
}
