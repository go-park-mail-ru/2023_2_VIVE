package usecase

import (
	"context"
	"errors"

	"HnH/pkg/serverErrors"
	pb "HnH/services/auth/authPB"
	"HnH/services/auth/internal/repository/redisRepo"
)

type IAuthUsecase interface {
	AddSession(context.Context, *pb.AuthData) (*pb.Empty, error)
	DeleteSession(context.Context, *pb.SessionID) (*pb.Empty, error)
	ValidateSession(context.Context, *pb.SessionID) (*pb.Empty, error)
	GetUserIdBySession(context.Context, *pb.SessionID) (*pb.UserID, error)
}

type AuthUsecase struct {
	authRepo redisRepo.IAuthRepository
}

func NewAuthUscase(authRepo redisRepo.IAuthRepository) IAuthUsecase {
	return &AuthUsecase{
		authRepo: authRepo,
	}
}

func (u *AuthUsecase) AddSession(ctx context.Context, authData *pb.AuthData) (*pb.Empty, error) {
	sessID := authData.SessionId.SessionId
	userID := authData.UserId.UserId
	expTime := authData.ExpiryTime

	err := u.authRepo.AddSession(ctx, sessID, userID, expTime)
	if errors.Is(err, redisRepo.ERROR_WHILE_WRITING) {
		return &pb.Empty{}, serverErrors.INTERNAL_SERVER_ERROR
	} else if err != nil {
		return &pb.Empty{}, err
	}

	return &pb.Empty{}, nil
}

func (u *AuthUsecase) DeleteSession(ctx context.Context, sessionID *pb.SessionID) (*pb.Empty, error) {
	sessID := sessionID.SessionId

	err := u.authRepo.DeleteSession(ctx, sessID)
	if err != nil {
		return &pb.Empty{}, nil
	}

	return &pb.Empty{}, nil
}

func (u *AuthUsecase) ValidateSession(ctx context.Context, sessionID *pb.SessionID) (*pb.Empty, error) {
	sessID := sessionID.SessionId

	err := u.authRepo.ValidateSession(ctx, sessID)
	if err != nil {
		return &pb.Empty{}, err
	}

	return &pb.Empty{}, nil
}

func (u *AuthUsecase) GetUserIdBySession(ctx context.Context, sessionID *pb.SessionID) (*pb.UserID, error) {
	sessID := sessionID.SessionId

	userID, err := u.authRepo.GetUserIdBySession(ctx, sessID)
	if err != nil {
		return &pb.UserID{}, err
	}

	return userID, nil
}
