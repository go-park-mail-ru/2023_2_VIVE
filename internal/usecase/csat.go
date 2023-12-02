package usecase

import (
	"HnH/internal/repository/grpc"
	"HnH/pkg/contextUtils"
	pb "HnH/services/csat/csatPB"
	"context"
)

type ICsatUsecase interface {
	GetQuestions(ctx context.Context) (*pb.QuestionList, error)
	RegisterAnswer(ctx context.Context, answer *pb.Answer) error
	GetStatistic(ctx context.Context) (*pb.Statistics, error)
}

type CvUsecase struct {
	csatRepo    grpc.ICsatRepository
	sessionRepo grpc.IAuthRepository
}

func NewCsatUsecase(
	csatRepository grpc.ICsatRepository,
	sessionRepository grpc.IAuthRepository,
) ICsatUsecase {
	return &CvUsecase{
		csatRepo:    csatRepository,
		sessionRepo: sessionRepository,
	}
}

func (u *CvUsecase) GetQuestions(ctx context.Context) (*pb.QuestionList, error) {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	questions, err := u.csatRepo.GetQuestions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (u *CvUsecase) RegisterAnswer(ctx context.Context, answer *pb.Answer) error {
	err := u.csatRepo.RegisterAnswer(ctx, answer)
	if err != nil {
		return err
	}
	return nil
}

func (u *CvUsecase) GetStatistic(ctx context.Context) (*pb.Statistics, error) {
	statistics, err := u.csatRepo.GetStatistic(ctx)
	if err != nil {
		return nil, err
	}
	return statistics, nil
}
