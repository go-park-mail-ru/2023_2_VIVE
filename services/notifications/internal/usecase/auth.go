package usecase

import (
	repository "HnH/services/notifications/internal/repository/grpc"
	"context"
)

type IAuthUsecase interface {
	ValidateAndGetUserID(ctx context.Context, sessionID string) (int64, error)
}

type AuthUsecase struct {
	authRepo repository.IAuthRepository
}

func NewAuthUsecase(authRepo repository.IAuthRepository) IAuthUsecase {
	return &AuthUsecase{
		authRepo: authRepo,
	}
}

func (u *AuthUsecase) ValidateAndGetUserID(ctx context.Context, sessionID string) (int64, error) {
	sessErr := u.authRepo.ValidateSession(ctx, sessionID)
	if sessErr != nil {
		return 0, sessErr
	}

	userID, userErr := u.authRepo.GetUserIdBySession(ctx, sessionID)
	if userErr != nil {
		return 0, userErr
	}
	return userID, nil
}
