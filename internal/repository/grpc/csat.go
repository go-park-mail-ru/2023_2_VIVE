package grpc

import (
	"HnH/pkg/contextUtils"
	pb "HnH/services/csat/csatPB"
	"context"

	"google.golang.org/grpc/metadata"
)

type ICsatRepository interface {
	GetQuestions(ctx context.Context, userID int) (*pb.QuestionList, error)
	RegisterAnswer(ctx context.Context, answer *pb.Answer) error
	GetStatistic(ctx context.Context) (*pb.Statistics, error)
}

type grpcCsatRepository struct {
	client pb.CsatClient
}

func NewGrpcSearchEngineRepository(client pb.CsatClient) ICsatRepository {
	return &grpcCsatRepository{
		client: client,
	}
}

func (repo *grpcCsatRepository) GetQuestions(ctx context.Context, userID int) (*pb.QuestionList, error) {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	userIDPB := pb.UserID{UserID: int64(userID)}
	questions, err := repo.client.GetQuestions(ctx, &userIDPB)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (repo *grpcCsatRepository) RegisterAnswer(ctx context.Context, answer *pb.Answer) error {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := repo.client.RegisterAnswer(ctx, answer)
	if err != nil {
		return err
	}
	return nil
}

func (repo *grpcCsatRepository) GetStatistic(ctx context.Context) (*pb.Statistics, error) {
	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)
	
	statistics, err := repo.client.GetStatistic(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	return statistics, nil
}
