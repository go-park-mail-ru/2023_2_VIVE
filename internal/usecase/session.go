package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/pkg/authUtils"
	"HnH/pkg/contextUtils"
	"context"

	"github.com/google/uuid"
)

type ISessionUsecase interface {
	Login(ctx context.Context, user *domain.DbUser, expiryUnixSeconds int64) (string, error)
	Logout(ctx context.Context) error
	CheckLogin(ctx context.Context) (int, error)
}

type SessionUsecase struct {
	sessionRepo grpc.IAuthRepository
	userRepo    psql.IUserRepository
}

func NewSessionUsecase(sessionRepository grpc.IAuthRepository, userRepository psql.IUserRepository) ISessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
	}
}

func (sessionUsecase *SessionUsecase) Login(ctx context.Context, user *domain.DbUser, expiryUnixSeconds int64) (string, error) {
	validEmailStatus := authUtils.ValidateEmail(user.Email)
	if validEmailStatus != nil {
		return "", validEmailStatus
	}

	validPasswordStatus := authUtils.IsPasswordEmpty(user.Password)
	if validPasswordStatus != nil {
		return "", validPasswordStatus
	}

	loginErr := sessionUsecase.userRepo.CheckUser(ctx, user)
	if loginErr != nil {
		return "", loginErr
	}

	userID, err := sessionUsecase.userRepo.GetUserIdByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}

	sessionID := uuid.NewString()

	addErr := sessionUsecase.sessionRepo.AddSession(ctx, sessionID, userID, expiryUnixSeconds)
	if addErr != nil {
		return "", addErr
	}

	return sessionID, nil
}

func (sessionUsecase *SessionUsecase) Logout(ctx context.Context) error {
	sessionID := contextUtils.GetSessionIDFromCtx(ctx)

	deleteErr := sessionUsecase.sessionRepo.DeleteSession(ctx, sessionID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (sessionUsecase *SessionUsecase) CheckLogin(ctx context.Context) (int, error) {
	sessionID := contextUtils.GetSessionIDFromCtx(ctx)

	userID, sessionErr := sessionUsecase.sessionRepo.GetUserIdBySession(ctx, sessionID)
	if sessionErr != nil {
		return 0, sessionErr
	}

	return userID, nil
}
