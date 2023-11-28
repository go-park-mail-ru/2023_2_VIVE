package grpc

import (
	"HnH/pkg/contextUtils"
	"HnH/services/csat/csatPB"
	"context"

	"google.golang.org/grpc/metadata"
)

type ICsatRepository interface {
	GetQuestions(ctx context.Context, userID int) (*csatPB.QuestionList, error)
	RegisterAnswer(ctx context.Context, answer *csatPB.Answer) error
	GetStatistic(ctx context.Context) (*csatPB.Statistics, error)
}

type grpcCsatRepository struct {
	client csatPB.CsatClient
}

func NewGrpcCsatRepository(client csatPB.CsatClient) ICsatRepository {
	return &grpcCsatRepository{
		client: client,
	}
}

func (repo *grpcCsatRepository) GetQuestions(ctx context.Context, userID int) (*csatPB.QuestionList, error) {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDFromCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	userIDPB := csatPB.UserID{UserID: int64(userID)}
	questions, err := repo.client.GetQuestions(ctx, &userIDPB)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (repo *grpcCsatRepository) RegisterAnswer(ctx context.Context, answer *csatPB.Answer) error {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDFromCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := repo.client.RegisterAnswer(ctx, answer)
	if err != nil {
		return err
	}
	return nil
}

func (repo *grpcCsatRepository) GetStatistic(ctx context.Context) (*csatPB.Statistics, error) {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDFromCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)
	
	statistics, err := repo.client.GetStatistic(ctx, &csatPB.Empty{})
	if err != nil {
		return nil, err
	}

	return statistics, nil
}
