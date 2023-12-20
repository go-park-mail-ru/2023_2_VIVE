package usecase

import (
	"HnH/internal/domain"
	psqlmock "HnH/internal/repository/mock"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"

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

func TestSignUpUserSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	defer mockCtrl.Finish()

	user := &domain.ApiUser{
		Email:     "vive342@mail.ru",
		Password:  "Vive2023top",
		FirstName: "Vladimir",
		LastName:  "Borozenets",
		Type:      domain.Applicant,
	}

	var sessID string
	userRepo.EXPECT().AddUser(ctx, user).Return(nil)
	userRepo.EXPECT().GetUserIdByEmail(ctx, "vive342@mail.ru").Return(123, nil)
	sessionRepo.EXPECT().AddSession(ctx, gomock.Any(), 123, int64(600)).Return(nil).DoAndReturn(
		func(ctx context.Context, sessionID string, userID int, expiryUnixSeconds int64) error {
			sessID = sessionID
			return nil
		})

	sID, err := userUsecase.SignUp(ctx, user, 600)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}

	assert.Equal(t, sID, sessID)
}

func TestGetInfoSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	defer mockCtrl.Finish()

	userDB := &domain.DbUser{
		ID:        10,
		Email:     "vive342@mail.ru",
		Password:  "Vive2023top",
		FirstName: "Vladimir",
		LastName:  "Borozenets",
		Type:      domain.Applicant,
	}

	id := 12
	user := &domain.ApiUser{
		ID:          10,
		ApplicantID: &id,
		Email:       "vive342@mail.ru",
		Password:    "Vive2023top",
		FirstName:   "Vladimir",
		LastName:    "Borozenets",
		Type:        domain.Applicant,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, user.ID)
	userRepo.EXPECT().GetUserInfo(ctxWithID, user.ID).Return(userDB, user.ApplicantID, nil, nil)

	u, err := userUsecase.GetInfo(ctxWithID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}

	assert.Equal(t, *u, *user)
}

func TestUpdateInfoSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	defer mockCtrl.Finish()

	id := 12
	birthday := "11.07.2002"
	phoneNumber := "+79992281437"
	location := "Москва"
	userUPD := &domain.UserUpdate{
		Email:       "vive342@mail.ru",
		FirstName:   "Vladimir",
		LastName:    "Borozenets",
		Birthday:    &birthday,
		PhoneNumber: &phoneNumber,
		Location:    &location,
		Password:    "Vive2023top",
		NewPassword: "Vive2023supertop",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, id)
	userRepo.EXPECT().CheckPasswordById(ctxWithID, id, userUPD.Password).Return(nil)
	userRepo.EXPECT().UpdateUserInfo(ctxWithID, id, userUPD).Return(nil)

	err := userUsecase.UpdateInfo(ctxWithID, userUPD)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestUploadAvatarSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	defer mockCtrl.Finish()

	h := textproto.MIMEHeader{}
	h.Add("Content-Type", "image/png")

	header := &multipart.FileHeader{
		Size:   10,
		Header: h,
	}

	id := 12

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, id)
	userRepo.EXPECT().GetAvatarByUserID(ctxWithID, id).Return("/ava/12.png", nil).Times(2)

	type impl struct {
		io.Reader
		io.ReaderAt
		io.Seeker
		io.Closer
	}

	file := impl{}
	err := userUsecase.UploadAvatar(ctxWithID, file, header)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetUserAvatarFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	defer mockCtrl.Finish()

	id := 12

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, id)
	userRepo.EXPECT().GetAvatarByUserID(ctxWithID, id).Return("/ava/12.png", nil)

	_, err := userUsecase.GetUserAvatar(ctxWithID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}
}

func TestGetUserEmptyAvatar(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	defer mockCtrl.Finish()

	id := 12

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, id)
	userRepo.EXPECT().GetAvatarByUserID(ctxWithID, id).Return("", nil)

	ava, err := userUsecase.GetUserAvatar(ctxWithID)
	if ava != nil || err != nil {
		fmt.Println("Error and response must be nil")
		t.Fail()
	}
}

func TestGetImageFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	defer mockCtrl.Finish()

	id := 12
	imageID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, id)
	userRepo.EXPECT().GetAvatarByUserID(ctxWithID, imageID).Return("avatar/12.jpg", nil)

	_, err := userUsecase.GetImage(ctxWithID, imageID)
	if err == nil {
		fmt.Println("Error and must be not nil")
		t.Fail()
	}
}

func TestGetEmptyImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)

	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	defer mockCtrl.Finish()

	id := 12
	imageID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, id)
	userRepo.EXPECT().GetAvatarByUserID(ctxWithID, imageID).Return("", nil)

	ava, err := userUsecase.GetImage(ctxWithID, imageID)
	if ava != nil || err != nil {
		fmt.Println("Error and response must be nil")
		t.Fail()
	}
}
