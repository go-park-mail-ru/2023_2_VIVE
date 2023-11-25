package usecase

import (
	pb "HnH/services/csat/csatPB"
	"HnH/services/csat/internal/repository/psql"
	"context"
	"time"
)

const (
	NO_SHOW_CSAT_TIME = time.Second * 5 // TODO: change
)

type ICsatUsecase interface {
	GetQuestions(ctx context.Context, request *pb.UserID) (*pb.QuestionList, error)
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
	if delta < NO_SHOW_CSAT_TIME {
		return false
	}
	return true
}

func (u *CsatUsecase) GetQuestions(ctx context.Context, userID *pb.UserID) (*pb.QuestionList, error) {
	res := pb.QuestionList{
		Questions: []*pb.Question{},
	}
	lastUpdate, err := u.csatRepo.GetLastUpdate(ctx, userID.UserID)
	if err == psql.ErrNoLastUpdate {
		return &res, nil
	}
	if err != nil {
		return &res, err
	}

	if !u.checkLastUpdate(lastUpdate) {
		return &res, nil
	}

	questions, err := u.csatRepo.GetQuestions(ctx)
	if err != nil {
		return &res, err
	}

	for _, question := range questions {
		res.Questions = append(res.Questions, &pb.Question{Question: question})
	}

	return &res, nil
}
