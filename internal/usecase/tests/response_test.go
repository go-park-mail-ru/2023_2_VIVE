package usecase

import (
	"HnH/internal/domain"
	psqlmock "HnH/internal/repository/mock"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/notificationMessages"
	"HnH/pkg/serverErrors"
	notificationsPB "HnH/services/notifications/api/proto"
	"context"
	"fmt"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRespondToVacancySuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)
	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	respUsecase := usecase.NewResponseUsecase(respRepo, sessionRepo, userRepo, vacancyRepo, cvRepo, notifRepo)

	defer mockCtrl.Finish()

	userID := 20
	var empID int64 = 21
	vacID := 99
	cvID := 12

	message := &notificationsPB.NotificationMessage{
		UserId:    empID,
		Message:   notificationMessages.NewVacancyResponse,
		VacancyId: int64(vacID),
		CvId:      int64(cvID),
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)
	respRepo.EXPECT().RespondToVacancy(ctxWithID, vacID, cvID).Return(nil)
	userRepo.EXPECT().GetUserIDByVacID(ctxWithID, vacID).Return(empID, nil)
	notifRepo.EXPECT().SendMessage(ctxWithID, message).Return(nil)

	err := respUsecase.RespondToVacancy(ctxWithID, vacID, cvID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestRespondToVacancyFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)
	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	respUsecase := usecase.NewResponseUsecase(respRepo, sessionRepo, userRepo, vacancyRepo, cvRepo, notifRepo)

	defer mockCtrl.Finish()

	userID := 20
	vacID := 99
	cvID := 12

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Employer, nil)

	err := respUsecase.RespondToVacancy(ctxWithID, vacID, cvID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}
}

func TestRespondToVacancyFail2(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)
	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	respUsecase := usecase.NewResponseUsecase(respRepo, sessionRepo, userRepo, vacancyRepo, cvRepo, notifRepo)

	defer mockCtrl.Finish()

	userID := 20
	vacID := 99
	cvID := 12

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, serverErrors.AUTH_REQUIRED)

	err := respUsecase.RespondToVacancy(ctxWithID, vacID, cvID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.AUTH_REQUIRED.Error())
}

func TestRespondToVacancyFail3(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)
	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	respUsecase := usecase.NewResponseUsecase(respRepo, sessionRepo, userRepo, vacancyRepo, cvRepo, notifRepo)

	defer mockCtrl.Finish()

	userID := 20
	vacID := 99
	cvID := 12

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)
	respRepo.EXPECT().RespondToVacancy(ctxWithID, vacID, cvID).Return(serverErrors.INTERNAL_SERVER_ERROR)

	err := respUsecase.RespondToVacancy(ctxWithID, vacID, cvID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR.Error())
}

func TestRespondToVacancyFail4(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)
	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	respUsecase := usecase.NewResponseUsecase(respRepo, sessionRepo, userRepo, vacancyRepo, cvRepo, notifRepo)

	defer mockCtrl.Finish()

	userID := 20
	vacID := 99
	cvID := 12

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)
	respRepo.EXPECT().RespondToVacancy(ctxWithID, vacID, cvID).Return(nil)
	userRepo.EXPECT().GetUserIDByVacID(ctxWithID, vacID).Return(int64(0), serverErrors.FORBIDDEN)

	err := respUsecase.RespondToVacancy(ctxWithID, vacID, cvID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.FORBIDDEN)
}

func TestRespondToVacancyFail5(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)
	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	respUsecase := usecase.NewResponseUsecase(respRepo, sessionRepo, userRepo, vacancyRepo, cvRepo, notifRepo)

	defer mockCtrl.Finish()

	userID := 20
	var empID int64 = 21
	vacID := 99
	cvID := 12

	message := &notificationsPB.NotificationMessage{
		UserId:    empID,
		Message:   notificationMessages.NewVacancyResponse,
		VacancyId: int64(vacID),
		CvId:      int64(cvID),
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)
	respRepo.EXPECT().RespondToVacancy(ctxWithID, vacID, cvID).Return(nil)
	userRepo.EXPECT().GetUserIDByVacID(ctxWithID, vacID).Return(empID, nil)
	notifRepo.EXPECT().SendMessage(ctxWithID, message).Return(serverErrors.INTERNAL_SERVER_ERROR)

	err := respUsecase.RespondToVacancy(ctxWithID, vacID, cvID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR.Error())
}

/*func TestGetApplicantsListSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)
	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	respUsecase := usecase.NewResponseUsecase(respRepo, sessionRepo, userRepo, vacancyRepo, cvRepo, notifRepo)

	defer mockCtrl.Finish()

	userID := 20
	vacID := 99

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Employer, nil)
	userRepo.EXPECT().GetUserEmpId(ctxWithID, userID).Return(21, nil)
	vacancyRepo.EXPECT().GetEmpId(ctxWithID, vacID).Return(21, nil)
	respRepo.EXPECT().GetAttachedCVs(ctxWithID, vacID).Return([]int{1, 2, 3, 4}, nil)
	cvRepo.EXPECT().GetCVsByIds(ctxWithID, gomock.Any()).Return(gomock.Any(), gomock.Any(), gomock.Any(), nil)

	_, err := respUsecase.GetApplicantsList(ctxWithID, vacID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}*/
