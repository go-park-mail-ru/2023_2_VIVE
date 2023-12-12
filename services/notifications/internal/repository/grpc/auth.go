package repository

import (
	"HnH/pkg/contextUtils"
	"HnH/services/auth/authPB"
	"context"
)

type IAuthRepository interface {
	ValidateSession(ctx context.Context, sessionID string) error
	GetUserIdBySession(ctx context.Context, sessionID string) (int64, error)
}

type grpcAuthRepository struct {
	client authPB.AuthClient
}

func NewGrpcAuthRepository(client authPB.AuthClient) IAuthRepository {
	return &grpcAuthRepository{
		client: client,
	}
}

func (repo *grpcAuthRepository) ValidateSession(ctx context.Context, sessionID string) error {
	// md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDFromCtx(ctx))
	// ctx = metadata.NewOutgoingContext(ctx, md)
	ctx = contextUtils.PutRequestIDToMetaDataCtx(ctx)

	sessID := authPB.SessionID{SessionId: sessionID}

	_, err := repo.client.ValidateSession(ctx, &sessID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *grpcAuthRepository) GetUserIdBySession(ctx context.Context, sessionID string) (int64, error) {
	// md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDFromCtx(ctx))
	// ctx = metadata.NewOutgoingContext(ctx, md)
	ctx = contextUtils.PutRequestIDToMetaDataCtx(ctx)

	sessID := authPB.SessionID{SessionId: sessionID}

	userID, err := repo.client.GetUserIdBySession(ctx, &sessID)
	if err != nil {
		return 0, err
	}

	return userID.UserId, nil
}
