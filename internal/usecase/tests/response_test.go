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

func TestGetApplicantsListSuccess(t *testing.T) {
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

	dbCV := []domain.DbCV{
		{
			ID:             20,
			ApplicantID:    40,
			ProfessionName: "Programmer",
			FirstName:      "Vladimir",
			LastName:       "Borozenets",
		},
		{
			ID:             21,
			ApplicantID:    40,
			ProfessionName: "Programmer",
			FirstName:      "Vladimir",
			LastName:       "Borozenets",
		}}

	expList := []domain.DbExperience{
		{
			ID:               1,
			CvID:             20,
			OrganizationName: "VK",
			Position:         "Senior",
			Description:      "Senior backender",
		},
		{
			ID:               2,
			CvID:             20,
			OrganizationName: "Yandex",
			Position:         "CTO",
			Description:      "Technical director",
		}}

	eduList := []domain.DbEducationInstitution{
		{
			ID:             3,
			CvID:           20,
			Name:           "BMSTU",
			MajorField:     "Math and CS",
			GraduationYear: "2024",
		},
		{
			ID:             4,
			CvID:           20,
			Name:           "MSU",
			MajorField:     "Mathematical analysis",
			GraduationYear: "2025",
		}}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Employer, nil)
	userRepo.EXPECT().GetUserEmpId(ctxWithID, userID).Return(21, nil)
	vacancyRepo.EXPECT().GetEmpId(ctxWithID, vacID).Return(21, nil)
	respRepo.EXPECT().GetAttachedCVs(ctxWithID, vacID).Return([]int{1, 2, 3, 4}, nil)
	cvRepo.EXPECT().GetCVsByIds(ctxWithID, gomock.Any()).Return(dbCV, expList, eduList, nil)

	_, err := respUsecase.GetApplicantsList(ctxWithID, vacID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetApplicantsListFail1(t *testing.T) {
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
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)

	_, err := respUsecase.GetApplicantsList(ctxWithID, vacID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, usecase.ErrInapropriateRole.Error())
}

func TestGetApplicantsListFail2(t *testing.T) {
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
	userRepo.EXPECT().GetUserEmpId(ctxWithID, userID).Return(21, serverErrors.INTERNAL_SERVER_ERROR)

	_, err := respUsecase.GetApplicantsList(ctxWithID, vacID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR.Error())
}

func TestGetApplicantsListFail3(t *testing.T) {
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
	vacancyRepo.EXPECT().GetEmpId(ctxWithID, vacID).Return(9, nil)

	_, err := respUsecase.GetApplicantsList(ctxWithID, vacID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.FORBIDDEN.Error())
}

func TestGetApplicantsListFail4(t *testing.T) {
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
	respRepo.EXPECT().GetAttachedCVs(ctxWithID, vacID).Return([]int{1, 2, 3, 4}, serverErrors.NO_DATA_FOUND)

	_, err := respUsecase.GetApplicantsList(ctxWithID, vacID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.NO_DATA_FOUND)
}

func TestGetApplicantsListFail5(t *testing.T) {
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

	dbCV := []domain.DbCV{
		{
			ID:             20,
			ApplicantID:    40,
			ProfessionName: "Programmer",
			FirstName:      "Vladimir",
			LastName:       "Borozenets",
		},
		{
			ID:             21,
			ApplicantID:    40,
			ProfessionName: "Programmer",
			FirstName:      "Vladimir",
			LastName:       "Borozenets",
		}}

	expList := []domain.DbExperience{
		{
			ID:               1,
			CvID:             20,
			OrganizationName: "VK",
			Position:         "Senior",
			Description:      "Senior backender",
		},
		{
			ID:               2,
			CvID:             20,
			OrganizationName: "Yandex",
			Position:         "CTO",
			Description:      "Technical director",
		}}

	eduList := []domain.DbEducationInstitution{
		{
			ID:             3,
			CvID:           20,
			Name:           "BMSTU",
			MajorField:     "Math and CS",
			GraduationYear: "2024",
		},
		{
			ID:             4,
			CvID:           20,
			Name:           "MSU",
			MajorField:     "Mathematical analysis",
			GraduationYear: "2025",
		}}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Employer, nil)
	userRepo.EXPECT().GetUserEmpId(ctxWithID, userID).Return(21, nil)
	vacancyRepo.EXPECT().GetEmpId(ctxWithID, vacID).Return(21, nil)
	respRepo.EXPECT().GetAttachedCVs(ctxWithID, vacID).Return([]int{1, 2, 3, 4}, nil)
	cvRepo.EXPECT().GetCVsByIds(ctxWithID, gomock.Any()).Return(dbCV, expList, eduList, serverErrors.INTERNAL_SERVER_ERROR)

	_, err := respUsecase.GetApplicantsList(ctxWithID, vacID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR.Error())
}

func TestGetUserResponsesSuccess(t *testing.T) {
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

	respList := []domain.ApiResponse{
		{
			Id:               1,
			VacancyName:      "C++",
			VacancyID:        1,
			OrganizationName: "VK",
			EmployerID:       1,
		},
		{
			Id:               2,
			VacancyName:      "Golang",
			VacancyID:        2,
			OrganizationName: "Yandex",
			EmployerID:       2,
		},
		{
			Id:               3,
			VacancyName:      "Java",
			VacancyID:        3,
			OrganizationName: "Croc",
			EmployerID:       3,
		}}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)
	respRepo.EXPECT().GetUserResponses(ctxWithID, userID).Return(respList, nil)

	_, err := respUsecase.GetUserResponses(ctxWithID, userID)
	if err != nil {
		fmt.Printf("Error must be nil: %v\n", err)
		t.Fail()
	}
}

func TestGetUserResponsesFail1(t *testing.T) {
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

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Employer, nil)

	_, err := respUsecase.GetUserResponses(ctxWithID, userID)
	if err == nil {
		fmt.Printf("Error must be not nil: %v\n", err)
		t.Fail()
	}

	assert.Error(t, err, usecase.ErrInapropriateRole.Error())
}

func TestGetUserResponsesFail2(t *testing.T) {
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

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, serverErrors.AUTH_REQUIRED)

	_, err := respUsecase.GetUserResponses(ctxWithID, userID)
	if err == nil {
		fmt.Printf("Error must be not nil: %v\n", err)
		t.Fail()
	}

	assert.Error(t, err, serverErrors.AUTH_REQUIRED)
}

func TestGetUserResponsesFail3(t *testing.T) {
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

	respList := []domain.ApiResponse{
		{
			Id:               1,
			VacancyName:      "C++",
			VacancyID:        1,
			OrganizationName: "VK",
			EmployerID:       1,
		},
		{
			Id:               2,
			VacancyName:      "Golang",
			VacancyID:        2,
			OrganizationName: "Yandex",
			EmployerID:       2,
		},
		{
			Id:               3,
			VacancyName:      "Java",
			VacancyID:        3,
			OrganizationName: "Croc",
			EmployerID:       3,
		}}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)
	respRepo.EXPECT().GetUserResponses(ctxWithID, userID).Return(respList, serverErrors.INTERNAL_SERVER_ERROR)

	_, err := respUsecase.GetUserResponses(ctxWithID, userID)
	if err == nil {
		fmt.Printf("Error must be not nil: %v\n", err)
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR)
}

func TestGetUserResponsesFail4(t *testing.T) {
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

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)

	_, err := respUsecase.GetUserResponses(ctxWithID, 15)
	if err == nil {
		fmt.Printf("Error must be not nil: %v\n", err)
		t.Fail()
	}

	assert.Error(t, err, usecase.ErrForbidden)
}
