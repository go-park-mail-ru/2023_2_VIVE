package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	"HnH/pkg/contextUtils"
	pb "HnH/services/csat/csatPB"
	"context"
)

type ICsatUsecase interface {
	GetQuestions(ctx context.Context) (*domain.QuestionList, error)
	RegisterAnswer(ctx context.Context, answer *domain.Answer) error
	GetStatistic(ctx context.Context) (*domain.Statistics, error)
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

func (u *CvUsecase) convertQuestionList(qList *pb.QuestionList) *domain.QuestionList {
	toReturn := &domain.QuestionList{}

	list := qList.Questions
	for _, question := range list {
		toAppend := domain.Question{}

		toAppend.Name = question.Name
		toAppend.Question = question.Question
		toAppend.QuestionId = question.QuestionId

		toReturn.Questions = append(toReturn.Questions, toAppend)
	}

	return toReturn
}

func (u *CvUsecase) convertAnswer(answer *domain.Answer) *pb.Answer {
	toReturn := &pb.Answer{}

	toReturn.QuestionId = answer.QuestionId
	toReturn.Starts = answer.Starts
	toReturn.Comment = answer.Comment

	return toReturn
}

func (u *CvUsecase) convertQuestionCommentList(list []*pb.QuestionComment) domain.QuestionCommentSlice {
	toReturn := domain.QuestionCommentSlice([]domain.QuestionComment{})

	for _, comment := range list {
		toAppend := domain.QuestionComment{}

		toAppend.Comment = comment.Comment

		toReturn = append(toReturn, toAppend)
	}

	return toReturn
}

func (u *CvUsecase) convertStarsNumList(list []*pb.StarsNum) domain.StarsNumSlice {
	toReturn := domain.StarsNumSlice([]domain.StarsNum{})

	for _, mark := range list {
		toAppend := domain.StarsNum{}

		toAppend.Count = mark.Count
		toAppend.StarsNum = mark.StarsNum

		toReturn = append(toReturn, toAppend)
	}

	return toReturn
}

func (u *CvUsecase) convertStatisctics(stat *pb.Statistics) *domain.Statistics {
	toReturn := &domain.Statistics{}

	list := stat.StatisticsList
	for _, questionStat := range list {
		toAppend := domain.QuestionStatistics{}

		toAppend.AvgStars = questionStat.AvgStars
		toAppend.QuestionText = questionStat.QuestionText
		toAppend.QuestionCommentList = u.convertQuestionCommentList(questionStat.QuestionCommentList)
		toAppend.StarsNumList = u.convertStarsNumList(questionStat.StarsNumList)

		toReturn.StatisticsList = append(toReturn.StatisticsList, toAppend)
	}

	return toReturn
}

func (u *CvUsecase) GetQuestions(ctx context.Context) (*domain.QuestionList, error) {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	questions, err := u.csatRepo.GetQuestions(ctx, userID)
	if err != nil {
		return nil, err
	}

	toReturn := u.convertQuestionList(questions)

	return toReturn, nil
}

func (u *CvUsecase) RegisterAnswer(ctx context.Context, answer *domain.Answer) error {
	ans := u.convertAnswer(answer)

	err := u.csatRepo.RegisterAnswer(ctx, ans)
	if err != nil {
		return err
	}

	return nil
}

func (u *CvUsecase) GetStatistic(ctx context.Context) (*domain.Statistics, error) {
	statistics, err := u.csatRepo.GetStatistic(ctx)
	if err != nil {
		return nil, err
	}

	toReturn := u.convertStatisctics(statistics)

	return toReturn, nil
}
