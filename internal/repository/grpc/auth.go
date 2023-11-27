package grpc

import (
	"HnH/pkg/contextUtils"
	pb "HnH/services/auth/authPB"
	"context"

	"google.golang.org/grpc/metadata"
)

type IAuthRepository interface {
	AddSession(ctx context.Context, sessionID string, userID int, expiryUnixSeconds int64) error
	DeleteSession(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, sessionID string) error
	GetUserIdBySession(ctx context.Context, sessionID string) (int, error)
}

type grpcAuthRepository struct {
	client pb.AuthClient
}

func NewGrpcAuthRepository(client pb.AuthClient) IAuthRepository {
	return &grpcAuthRepository{
		client: client,
	}
}

func (repo *grpcAuthRepository) AddSession(ctx context.Context, sessionID string, userID int, expiryUnixSeconds int64) error {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	sessID := pb.SessionID{SessionId: sessionID}
	uID := pb.UserID{UserId: int64(userID)}
	authData := pb.AuthData{
		SessionId:  &sessID,
		UserId:     &uID,
		ExpiryTime: expiryUnixSeconds,
	}

	_, err := repo.client.AddSession(ctx, &authData)
	if err != nil {
		return err
	}

	return nil
}

func (repo *grpcAuthRepository) DeleteSession(ctx context.Context, sessionID string) error {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	sessID := pb.SessionID{SessionId: sessionID}

	_, err := repo.client.DeleteSession(ctx, &sessID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *grpcAuthRepository) ValidateSession(ctx context.Context, sessionID string) error {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	sessID := pb.SessionID{SessionId: sessionID}

	_, err := repo.client.ValidateSession(ctx, &sessID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *grpcAuthRepository) GetUserIdBySession(ctx context.Context, sessionID string) (int, error) {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	sessID := pb.SessionID{SessionId: sessionID}

	userID, err := repo.client.GetUserIdBySession(ctx, &sessID)
	if err != nil {
		return 0, err
	}

	return int(userID.UserId), nil
}
