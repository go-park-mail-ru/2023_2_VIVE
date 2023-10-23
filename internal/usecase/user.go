package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"
	"HnH/pkg/authUtils"
	"HnH/pkg/serverErrors"

	"github.com/google/uuid"
)

type IUserUsecase interface {
	SignUp(user *domain.User) (string, error)
	GetInfo(sessionID string) (*domain.User, error)
	UpdateInfo(sessionID string, user *domain.UserUpdate) error
}

type UserUsecase struct {
	userRepo    repository.IUserRepository
	sessionRepo repository.ISessionRepository
}

func NewUserUsecase(userRepository repository.IUserRepository, sessionRepository repository.ISessionRepository) IUserUsecase {
	return &UserUsecase{
		userRepo:    userRepository,
		sessionRepo: sessionRepository,
	}
}

func (userUsecase *UserUsecase) SignUp(user *domain.User) (string, error) {
	validEmailStatus := authUtils.ValidateEmail(user.Email)
	if validEmailStatus != nil {
		return "", validEmailStatus
	}

	validPassStatus := authUtils.ValidatePassword(user.Password)
	if validPassStatus != nil {
		return "", validPassStatus
	}

	if !user.Type.IsRole() {
		return "", serverErrors.INVALID_ROLE
	}

	addStatus := userUsecase.userRepo.AddUser(user)
	if addStatus != nil {
		return "", addStatus
	}

	userID, err := userUsecase.userRepo.GetUserIdByEmail(user.Email)
	if err != nil {
		return "", err
	}

	sessionID := uuid.NewString()

	addErr := userUsecase.sessionRepo.AddSession(sessionID, userID)
	if addErr != nil {
		return "", addErr
	}

	return sessionID, nil
}

func (userUsecase *UserUsecase) GetInfo(sessionID string) (*domain.User, error) {
	validStatus := userUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	userID, err := userUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return nil, err
	}

	user, getErr := userUsecase.userRepo.GetUserInfo(userID)
	if getErr != nil {
		return nil, getErr
	}

	return user, nil
}

func (userUsecase *UserUsecase) UpdateInfo(sessionID string, user *domain.UserUpdate) error {
	validStatus := userUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return validStatus
	}

	userID, err := userUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return err
	}

	validPassStatus := userUsecase.userRepo.CheckPasswordById(userID, user.Password)
	if validPassStatus != nil {
		return validPassStatus
	}

	updStatus := userUsecase.userRepo.UpdateUserInfo(user)
	if updStatus != nil {
		return updStatus
	}

	return nil
}
