package usecase

import (
	"HnH/internal/domain"
	"HnH/pkg/authUtils"

	"github.com/google/uuid"
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

type SessionUsecase struct {
	sessionRepo SessionRepository
	userRepo    UserRepository
}

func NewSessionUsecase(sessionRepository SessionRepository, userRepository UserRepository) *SessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
	}
}

func (sessionUsecase *SessionUsecase) Login(user *domain.User) (string, error) {
	validEmailStatus := authUtils.ValidateEmail(user.Email)
	if validEmailStatus != nil {
		return "", validEmailStatus
	}

	validPasswordStatus := authUtils.IsPasswordEmpty(user.Password)
	if validPasswordStatus != nil {
		return "", validPasswordStatus
	}

	loginErr := sessionUsecase.userRepo.CheckUser(user)
	if loginErr != nil {
		return "", loginErr
	}

	userID, err := sessionUsecase.userRepo.GetUserIdByEmail(user.Email)
	if err != nil {
		return "", err
	}

	sessionID := uuid.NewString()

	addErr := sessionUsecase.sessionRepo.AddSession(sessionID, userID)
	if addErr != nil {
		return "", addErr
	}

	return sessionID, nil
}

func (sessionUsecase *SessionUsecase) Logout(sessionID string) error {
	deleteErr := sessionUsecase.sessionRepo.DeleteSession(sessionID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (sessionUsecase *SessionUsecase) CheckLogin(sessionID string) error {
	sessionErr := sessionUsecase.sessionRepo.ValidateSession(sessionID)
	if sessionErr != nil {
		return sessionErr
	}

	return nil
}
