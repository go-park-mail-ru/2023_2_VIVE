package usecase

import (
	"HnH/pkg/serverErrors"
	pb "HnH/services/csat/csatPB"
	"HnH/services/csat/internal/repository/psql"
	"context"
	"time"
)

const (
	NO_SHOW_CSAT_TIME = time.Second * 5 // TODO: change
)

type ICsatUsecase interface {
	GetQuestions(ctx context.Context, userID *pb.UserID) (*pb.QuestionList, error)
	RegisterAnswer(ctx context.Context, answer *pb.Answer) (*pb.Empty, error)
	GetStatistics(ctx context.Context) (*pb.Statistics, error)
}

type CsatUsecase struct {
	csatRepo psql.ICsatRepository
}

func NewCsatUscase(csatRepo psql.ICsatRepository) ICsatUsecase {
	return &CsatUsecase{
		csatRepo: csatRepo,
	}
}

func (u *CsatUsecase) checkLastUpdate(lastUpdate time.Time) bool {
	now := time.Now()

	delta := now.Sub(lastUpdate)
	return delta >= NO_SHOW_CSAT_TIME
}

func (u *CsatUsecase) GetQuestions(ctx context.Context, userID *pb.UserID) (*pb.QuestionList, error) {
	res := pb.QuestionList{
		Questions: []*pb.Question{},
	}
	lastUpdate, err := u.csatRepo.GetLastUpdate(ctx, userID.UserID)
	// if err == psql.ErrEntityNotFound {
	// 	return &res, nil
	// }
	if err != nil && err != serverErrors.ErrEntityNotFound {
		return &res, err
	}

	if !u.checkLastUpdate(lastUpdate) {
		return &res, nil
	}

	questions, err := u.csatRepo.GetQuestions(ctx)
	if err != nil {
		return &res, err
	}

	res.Questions = append(res.Questions, questions...)
	// for _, question := range questions {
	// 	res.Questions = append(res.Questions, question)
	// }

	return &res, nil
}

func (u *CsatUsecase) RegisterAnswer(ctx context.Context, answer *pb.Answer) (*pb.Empty, error) {
	res := pb.Empty{}
	err := u.csatRepo.RegisterAnswer(ctx, answer)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

func (u *CsatUsecase) GetStatistics(ctx context.Context) (*pb.Statistics, error) {
	res := pb.Statistics{}
	statistics, err := u.csatRepo.GetStatistics(ctx)
	if err != nil {
		return &res, err
	}

	return statistics, nil
}
