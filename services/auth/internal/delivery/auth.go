package delivery

import (
	pb "HnH/services/auth/authPB"
	"HnH/services/auth/internal/repository/redisRepo"
	"HnH/services/auth/internal/usecase"
	"context"

	"github.com/gomodule/redigo/redis"
)

type AuthServer struct {
	pb.UnimplementedAuthServer

	authUscase usecase.IAuthUsecase
}

func NewServer(conn *redis.Pool) (*AuthServer, error) {
	authRepo := redisRepo.NewRedisAuthRepository(conn)
	authUsecase := usecase.NewAuthUscase(authRepo)

	return &AuthServer{
		authUscase: authUsecase,
	}, nil
}

func (authServer *AuthServer) AddSession(ctx context.Context, authData *pb.AuthData) (*pb.Empty, error) {
	return authServer.authUscase.AddSession(ctx, authData)
}

func (authServer *AuthServer) DeleteSession(ctx context.Context, sessionID *pb.SessionID) (*pb.Empty, error) {
	return authServer.authUscase.DeleteSession(ctx, sessionID)
}

func (authServer *AuthServer) ValidateSession(ctx context.Context, sessionID *pb.SessionID) (*pb.Empty, error) {
	return authServer.authUscase.ValidateSession(ctx, sessionID)
}

func (authServer *AuthServer) GetUserIdBySession(ctx context.Context, sessionID *pb.SessionID) (*pb.UserID, error) {
	return authServer.authUscase.GetUserIdBySession(ctx, sessionID)
}
