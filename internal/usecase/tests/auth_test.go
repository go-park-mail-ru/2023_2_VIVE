package usecase

import (
	"HnH/internal/domain"
	psqlmock "HnH/internal/repository/mock"
	"HnH/internal/usecase"
	"HnH/pkg/authUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"context"
	"fmt"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20

	dbUser := &domain.DbUser{
		ID:        20,
		Email:     "vive20000@mail.ru",
		Password:  "Vive2023top!",
		FirstName: "Vladimir",
		LastName:  "Borozenets",
		Type:      domain.Applicant,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().CheckUser(ctxWithID, dbUser).Return(nil)
	userRepo.EXPECT().GetUserIdByEmail(ctxWithID, dbUser.Email).Return(dbUser.ID, nil)
	sessionRepo.EXPECT().AddSession(ctxWithID, gomock.Any(), dbUser.ID, int64(600)).Return(nil)

	_, err := sessUsecase.Login(ctxWithID, dbUser, 600)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestLoginFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20

	dbUser := &domain.DbUser{
		ID:        20,
		Email:     "",
		Password:  "Vive2023top!",
		FirstName: "Vladimir",
		LastName:  "Borozenets",
		Type:      domain.Applicant,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)

	_, err := sessUsecase.Login(ctxWithID, dbUser, 600)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, authUtils.EMPTY_EMAIL.Error())
}

func TestLoginFail2(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20

	dbUser := &domain.DbUser{
		ID:        20,
		Email:     "vive20000@mail.ru",
		Password:  "",
		FirstName: "Vladimir",
		LastName:  "Borozenets",
		Type:      domain.Applicant,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)

	_, err := sessUsecase.Login(ctxWithID, dbUser, 600)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, authUtils.EMPTY_PASSWORD.Error())
}

func TestLoginFail3(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20

	dbUser := &domain.DbUser{
		ID:        20,
		Email:     "vive20000@mail.ru",
		Password:  "Vive2023top!",
		FirstName: "Vladimir",
		LastName:  "Borozenets",
		Type:      domain.Applicant,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().CheckUser(ctxWithID, dbUser).Return(serverErrors.INCORRECT_CREDENTIALS)

	_, err := sessUsecase.Login(ctxWithID, dbUser, 600)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INCORRECT_CREDENTIALS.Error())
}

func TestLoginFail4(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20

	dbUser := &domain.DbUser{
		ID:        20,
		Email:     "vive20000@mail.ru",
		Password:  "Vive2023top!",
		FirstName: "Vladimir",
		LastName:  "Borozenets",
		Type:      domain.Applicant,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().CheckUser(ctxWithID, dbUser).Return(nil)
	userRepo.EXPECT().GetUserIdByEmail(ctxWithID, dbUser.Email).Return(0, serverErrors.INVALID_EMAIL)

	_, err := sessUsecase.Login(ctxWithID, dbUser, 600)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INVALID_EMAIL.Error())
}

func TestLoginFail5(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20

	dbUser := &domain.DbUser{
		ID:        20,
		Email:     "vive20000@mail.ru",
		Password:  "Vive2023top!",
		FirstName: "Vladimir",
		LastName:  "Borozenets",
		Type:      domain.Applicant,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().CheckUser(ctxWithID, dbUser).Return(nil)
	userRepo.EXPECT().GetUserIdByEmail(ctxWithID, dbUser.Email).Return(dbUser.ID, nil)
	sessionRepo.EXPECT().AddSession(ctxWithID, gomock.Any(), dbUser.ID, int64(600)).Return(serverErrors.INTERNAL_SERVER_ERROR)

	_, err := sessUsecase.Login(ctxWithID, dbUser, 600)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR.Error())
}

func TestLogoutSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20
	sessionID := "asdzxcegb"

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	ctxWithSessID := context.WithValue(ctxWithID, contextUtils.SESSION_ID_KEY, sessionID)
	sessionRepo.EXPECT().DeleteSession(ctxWithSessID, gomock.Any()).Return(nil)

	err := sessUsecase.Logout(ctxWithSessID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestLogoutFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20
	sessionID := "asdzxcegb"

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	ctxWithSessID := context.WithValue(ctxWithID, contextUtils.SESSION_ID_KEY, sessionID)
	sessionRepo.EXPECT().DeleteSession(ctxWithSessID, gomock.Any()).Return(serverErrors.NO_COOKIE)

	err := sessUsecase.Logout(ctxWithSessID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.NO_COOKIE.Error())
}

func TestCheckLoginSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20
	sessionID := "asdzxcegb"

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	ctxWithSessID := context.WithValue(ctxWithID, contextUtils.SESSION_ID_KEY, sessionID)
	sessionRepo.EXPECT().GetUserIdBySession(ctxWithSessID, sessionID).Return(userID, nil)

	_, err := sessUsecase.CheckLogin(ctxWithSessID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestCheckLoginFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	sessUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)

	defer mockCtrl.Finish()

	userID := 20
	sessionID := "asdzxcegb"

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	ctxWithSessID := context.WithValue(ctxWithID, contextUtils.SESSION_ID_KEY, sessionID)
	sessionRepo.EXPECT().GetUserIdBySession(ctxWithSessID, sessionID).Return(0, serverErrors.AUTH_REQUIRED)

	_, err := sessUsecase.CheckLogin(ctxWithSessID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.AUTH_REQUIRED)
}
