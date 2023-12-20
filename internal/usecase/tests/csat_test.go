package usecase

import (
	"HnH/internal/domain"
	psqlmock "HnH/internal/repository/mock"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	pb "HnH/services/csat/csatPB"
	"context"
	"fmt"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetQuestionsSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	csatRepo := psqlmock.NewMockICsatRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	csatUsecase := usecase.NewCsatUsecase(csatRepo, sessionRepo)

	defer mockCtrl.Finish()

	userID := 20

	q1 := &pb.Question{
		QuestionId: 1,
		Question:   "nice one",
		Name:       "first",
	}
	q2 := &pb.Question{
		QuestionId: 2,
		Question:   "nice two",
		Name:       "second",
	}

	list := &pb.QuestionList{
		Questions: []*pb.Question{q1, q2},
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	csatRepo.EXPECT().GetQuestions(ctxWithID, userID).Return(list, nil)

	_, err := csatUsecase.GetQuestions(ctxWithID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetQuestionsFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	csatRepo := psqlmock.NewMockICsatRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	csatUsecase := usecase.NewCsatUsecase(csatRepo, sessionRepo)

	defer mockCtrl.Finish()

	userID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	csatRepo.EXPECT().GetQuestions(ctxWithID, userID).Return(nil, serverErrors.INTERNAL_SERVER_ERROR)

	_, err := csatUsecase.GetQuestions(ctxWithID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR.Error())
}

func TestRegisterAnswerSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	csatRepo := psqlmock.NewMockICsatRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	csatUsecase := usecase.NewCsatUsecase(csatRepo, sessionRepo)

	defer mockCtrl.Finish()

	userID := 20

	ans := &domain.Answer{
		Starts:     5,
		Comment:    "WOW!",
		QuestionId: 1,
	}

	pbAns := &pb.Answer{
		Starts:     5,
		Comment:    "WOW!",
		QuestionId: 1,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	csatRepo.EXPECT().RegisterAnswer(ctxWithID, pbAns).Return(nil)

	err := csatUsecase.RegisterAnswer(ctxWithID, ans)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestRegisterAnswerFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	csatRepo := psqlmock.NewMockICsatRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	csatUsecase := usecase.NewCsatUsecase(csatRepo, sessionRepo)

	defer mockCtrl.Finish()

	userID := 20

	ans := &domain.Answer{
		Starts:     5,
		Comment:    "WOW!",
		QuestionId: 1,
	}

	pbAns := &pb.Answer{
		Starts:     5,
		Comment:    "WOW!",
		QuestionId: 1,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	csatRepo.EXPECT().RegisterAnswer(ctxWithID, pbAns).Return(serverErrors.INTERNAL_SERVER_ERROR)

	err := csatUsecase.RegisterAnswer(ctxWithID, ans)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR)
}

func TestGetStatisticsSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	csatRepo := psqlmock.NewMockICsatRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	csatUsecase := usecase.NewCsatUsecase(csatRepo, sessionRepo)

	defer mockCtrl.Finish()

	userID := 20

	stars := &pb.StarsNum{
		StarsNum: 5,
		Count:    10,
	}

	qCom := &pb.QuestionComment{
		Comment: "Very good!",
	}

	q1stat := &pb.QuestionStatistics{
		AvgStars:            3.5,
		QuestionText:        "Do you like it?",
		StarsNumList:        []*pb.StarsNum{stars},
		QuestionCommentList: []*pb.QuestionComment{qCom},
	}

	stat := &pb.Statistics{
		StatisticsList: []*pb.QuestionStatistics{q1stat},
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	csatRepo.EXPECT().GetStatistic(ctxWithID).Return(stat, nil)

	_, err := csatUsecase.GetStatistic(ctxWithID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetStatisticsFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	csatRepo := psqlmock.NewMockICsatRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	csatUsecase := usecase.NewCsatUsecase(csatRepo, sessionRepo)

	defer mockCtrl.Finish()

	userID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	csatRepo.EXPECT().GetStatistic(ctxWithID).Return(nil, serverErrors.NO_DATA_FOUND)

	_, err := csatUsecase.GetStatistic(ctxWithID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.NO_DATA_FOUND)
}
