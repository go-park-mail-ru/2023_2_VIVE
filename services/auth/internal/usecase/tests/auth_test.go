package usecase

import (
	"HnH/pkg/authUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	pb "HnH/services/auth/authPB"
	psqlmock "HnH/services/auth/internal/repository/mock"
	"HnH/services/auth/internal/usecase"
	"context"
	"fmt"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	file interface{}
}

func (w mockWriter) Write(p []byte) (n int, err error) {
	return 1, nil
}

var newWriter = mockWriter{
	file: 5,
}

var ctxBack = context.Background()
var logger = logrus.Logger{
	Out: newWriter,
}

var ctxLogger = logger.WithField("testing", "gomock")

var ctx = context.WithValue(ctxBack, contextUtils.LOGGER_KEY, ctxLogger)

func TestAddSessionSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	uID := &pb.UserID{
		UserId: int64(userID),
	}

	data := &pb.AuthData{
		SessionId:  sessID,
		UserId:     uID,
		ExpiryTime: int64(600),
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().AddSession(ctxWithID, data.SessionId.SessionId, data.UserId.UserId, data.ExpiryTime).Return(nil)

	_, err := sessUsecase.AddSession(ctxWithID, data)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestLoginFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	uID := &pb.UserID{
		UserId: int64(userID),
	}

	data := &pb.AuthData{
		SessionId:  sessID,
		UserId:     uID,
		ExpiryTime: int64(600),
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().AddSession(ctxWithID, data.SessionId.SessionId, data.UserId.UserId, data.ExpiryTime).Return(authUtils.ERROR_WHILE_WRITING)

	_, err := sessUsecase.AddSession(ctxWithID, data)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR)
}

func TestLoginFail2(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	uID := &pb.UserID{
		UserId: int64(userID),
	}

	data := &pb.AuthData{
		SessionId:  sessID,
		UserId:     uID,
		ExpiryTime: int64(600),
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().AddSession(ctxWithID, data.SessionId.SessionId, data.UserId.UserId, data.ExpiryTime).Return(authUtils.INCORRECT_CREDENTIALS)

	_, err := sessUsecase.AddSession(ctxWithID, data)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, authUtils.INCORRECT_CREDENTIALS)
}

func TestDeleteSessionSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().DeleteSession(ctxWithID, sessID.SessionId).Return(nil)

	_, err := sessUsecase.DeleteSession(ctxWithID, sessID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestDeleteSessionFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().DeleteSession(ctxWithID, sessID.SessionId).Return(authUtils.ERROR_WHILE_DELETING)

	_, err := sessUsecase.DeleteSession(ctxWithID, sessID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, authUtils.ERROR_WHILE_DELETING)
}

func TestValidateSessionSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().ValidateSession(ctxWithID, sessID.SessionId).Return(nil)

	_, err := sessUsecase.ValidateSession(ctxWithID, sessID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestValidateSessionFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().ValidateSession(ctxWithID, sessID.SessionId).Return(serverErrors.NO_SESSION)

	_, err := sessUsecase.ValidateSession(ctxWithID, sessID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.NO_SESSION)
}

func TestGetUserIDBySessionSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	uID := &pb.UserID{
		UserId: int64(userID),
	}

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().GetUserIdBySession(ctxWithID, sessID.SessionId).Return(uID, nil)

	_, err := sessUsecase.GetUserIdBySession(ctxWithID, sessID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetUserIDBySessionFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	authRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewAuthUscase(authRepo)

	defer mockCtrl.Finish()

	userID := 20

	sessID := &pb.SessionID{
		SessionId: "someID",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	authRepo.EXPECT().GetUserIdBySession(ctxWithID, sessID.SessionId).Return(nil, serverErrors.INTERNAL_SERVER_ERROR)

	_, err := sessUsecase.GetUserIdBySession(ctxWithID, sessID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR)
}
