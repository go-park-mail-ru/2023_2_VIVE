package usecase

import (
	"HnH/internal/domain"
	"HnH/pkg/authUtils"
	"HnH/pkg/serverErrors"

	"github.com/google/uuid"
)

type UserUsecase struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewUserUsecase(userRepository UserRepository, sessionRepository SessionRepository) *UserUsecase {
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
