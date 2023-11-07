package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/authUtils"

	"github.com/google/uuid"
)

type ISessionUsecase interface {
	Login(user *domain.User, expiryUnixSeconds int64) (string, error)
	Logout(sessionID string) error
	CheckLogin(sessionID string) error
}

type SessionUsecase struct {
	sessionRepo redisRepo.ISessionRepository
	userRepo    psql.IUserRepository
}

func NewSessionUsecase(sessionRepository redisRepo.ISessionRepository, userRepository psql.IUserRepository) ISessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
	}
}

func (sessionUsecase *SessionUsecase) Login(user *domain.User, expiryUnixSeconds int64) (string, error) {
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

	addErr := sessionUsecase.sessionRepo.AddSession(sessionID, userID, expiryUnixSeconds)
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
