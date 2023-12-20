package usecase

import (
	"HnH/internal/domain"
	psqlmock "HnH/internal/repository/mock"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"HnH/services/searchEngineService/searchEnginePB"
	"context"
	"fmt"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetCVByIDSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	cvID := 20

	dbCV := &domain.DbCV{
		ID:             20,
		ApplicantID:    15,
		ProfessionName: "Programmer",
		FirstName:      "Vladimir",
		LastName:       "Borozenets",
	}

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

	cvIDToLogo := map[int]string{
		20: "/avas/coolest/12.png",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	cvRepo.EXPECT().GetCVById(ctxWithID, cvID).Return(dbCV, expList, eduList, nil)
	skillRepo.EXPECT().GetSkillsByCvID(ctxWithID, cvID).Return([]string{"go", "git", "sql"}, nil)
	userRepo.EXPECT().GetAvatarPathesByCVIDList(ctxWithID, cvID).Return(cvIDToLogo, nil)

	_, err := cvUsecase.GetCVById(ctxWithID, cvID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetCVListSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	cvID := 20

	dbCV := []domain.DbCV{
		{
			ID:             20,
			ApplicantID:    15,
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

	cvIDToLogo := map[int]string{
		20: "/avas/coolest/12.png",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	userRepo.EXPECT().GetRoleById(ctxWithID, userID).Return(domain.Applicant, nil)
	cvRepo.EXPECT().GetCVsByUserId(ctxWithID, userID).Return(dbCV, expList, eduList, nil)
	skillRepo.EXPECT().GetSkillsByCvID(ctxWithID, cvID).Return([]string{"go", "git", "sql"}, nil)
	userRepo.EXPECT().GetAvatarPathesByCVIDList(ctxWithID, cvID).Return(cvIDToLogo, nil)

	_, err := cvUsecase.GetCVList(ctxWithID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestAddNewCVSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	cvID := 20

	dbCV := &domain.DbCV{
		ID:             20,
		ApplicantID:    15,
		ProfessionName: "Programmer",
		FirstName:      "Vladimir",
		LastName:       "Borozenets",
	}

	apiCV := &domain.ApiCV{
		ID:             20,
		ApplicantID:    15,
		ProfessionName: "Programmer",
		FirstName:      "Vladimir",
		LastName:       "Borozenets",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	cvRepo.EXPECT().AddCV(ctxWithID, userID, dbCV, gomock.Any(), gomock.Any()).Return(cvID, nil)
	skillRepo.EXPECT().AddSkillsByCvID(ctxWithID, cvID, gomock.Any()).Return(nil)

	_, err := cvUsecase.AddNewCV(ctxWithID, apiCV)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetCVOfUserByIdSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	cvID := 20

	dbCV := &domain.DbCV{
		ID:             20,
		ApplicantID:    15,
		ProfessionName: "Programmer",
		FirstName:      "Vladimir",
		LastName:       "Borozenets",
	}

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

	cvIDToLogo := map[int]string{
		20: "/avas/coolest/12.png",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	cvRepo.EXPECT().GetOneOfUsersCV(ctxWithID, userID, cvID).Return(dbCV, expList, eduList, nil)
	skillRepo.EXPECT().GetSkillsByCvID(ctxWithID, cvID).Return([]string{"go", "git", "sql"}, nil)
	userRepo.EXPECT().GetAvatarPathesByCVIDList(ctxWithID, cvID).Return(cvIDToLogo, nil)

	_, err := cvUsecase.GetCVOfUserById(ctxWithID, cvID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetApplicantInfoSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12
	applicantID := 40

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
			ProfessionName: "Data Satanist",
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

	cvIDToLogo := map[int]string{
		20: "/avas/coolest/12.png",
		21: "",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	cvRepo.EXPECT().GetApplicantInfo(ctxWithID, applicantID).Return("Vladimir", "Borozenets", dbCV, expList, eduList, nil)
	userRepo.EXPECT().GetAvatarPathesByCVIDList(ctxWithID, 20, 21).Return(cvIDToLogo, nil)

	_, err := cvUsecase.GetApplicantInfo(ctxWithID, applicantID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestUpdateCVOfUserByIDSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12
	cvID := 20

	dbCV := &domain.DbCV{

		ID:             20,
		ApplicantID:    40,
		ProfessionName: "Programmer",
		FirstName:      "Vladimir",
		LastName:       "Borozenets",
	}

	apiCV := &domain.ApiCV{

		ID:             20,
		ApplicantID:    40,
		ProfessionName: "Programmer",
		FirstName:      "Vladimir",
		LastName:       "Borozenets",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	expRepo.EXPECT().GetCVExperiencesIDs(ctxWithID, cvID).Return([]int{1, 2}, nil)
	instRepo.EXPECT().GetCVInstitutionsIDs(ctxWithID, cvID).Return([]int{3, 4}, nil)
	cvRepo.EXPECT().UpdateOneOfUsersCV(ctxWithID, userID, cvID, dbCV, []int{1, 2}, gomock.Any(), gomock.Any(), []int{3, 4}, gomock.Any(), gomock.Any()).Return(nil)
	skillRepo.EXPECT().UpdateSkillsByCvID(ctxWithID, cvID, gomock.Any()).Return(nil)

	err := cvUsecase.UpdateCVOfUserById(ctxWithID, cvID, apiCV)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestDeleteCVOfUserByIDSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12
	cvID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	cvRepo.EXPECT().DeleteOneOfUsersCV(ctxWithID, userID, cvID).Return(nil)

	err := cvUsecase.DeleteCVOfUserById(ctxWithID, cvID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestDeleteCVOfUserByIDFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12
	cvID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	cvRepo.EXPECT().DeleteOneOfUsersCV(ctxWithID, userID, cvID).Return(serverErrors.FORBIDDEN)

	err := cvUsecase.DeleteCVOfUserById(ctxWithID, cvID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.FORBIDDEN.Error())
}

func TestSearchCVSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

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

	opts := &searchEnginePB.SearchOptions{}
	resp := &searchEnginePB.SearchResponse{
		Ids:   []int64{10, 11, 12},
		Count: 3,
	}

	cvIDToLogo := map[int]string{
		20: "/avas/coolest/12.png",
		21: "",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	searchRepo.EXPECT().SearchCVsIDs(ctxWithID, opts).Return(resp, nil)
	cvRepo.EXPECT().GetCVsByIds(ctxWithID, gomock.Any()).Return(dbCV, expList, eduList, nil)
	skillRepo.EXPECT().GetSkillsByCvID(ctxWithID, 20).Return([]string{"go", "git"}, nil)
	skillRepo.EXPECT().GetSkillsByCvID(ctxWithID, 21).Return([]string{"c++", "sql"}, nil)
	userRepo.EXPECT().GetAvatarPathesByCVIDList(ctxWithID, 20, 21).Return(cvIDToLogo, nil)

	_, err := cvUsecase.SearchCVs(ctxWithID, opts)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

/*func TestGenerateCVsPDFSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)
	cvRepo := psqlmock.NewMockICVRepository(mockCtrl)
	expRepo := psqlmock.NewMockIExperienceRepository(mockCtrl)
	instRepo := psqlmock.NewMockIEducationInstitutionRepository(mockCtrl)
	respRepo := psqlmock.NewMockIResponseRepository(mockCtrl)

	cvUsecase := usecase.NewCVUsecase(cvRepo, expRepo, instRepo, sessionRepo, userRepo, respRepo, vacancyRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	cvID := 20

	dbCV := &domain.DbCV{
		ID:             20,
		ApplicantID:    15,
		ProfessionName: "Programmer",
		FirstName:      "Vladimir",
		LastName:       "Borozenets",
	}

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

	cvIDToLogo := map[int]string{
		20: "/avas/coolest/12.png",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	cvRepo.EXPECT().GetCVById(ctxWithID, cvID).Return(dbCV, expList, eduList, nil)
	skillRepo.EXPECT().GetSkillsByCvID(ctxWithID, cvID).Return([]string{"go", "git", "sql"}, nil)
	userRepo.EXPECT().GetAvatarPathesByCVIDList(ctxWithID, cvID).Return(cvIDToLogo, nil)
	cvRepo.EXPECT().DeleteOneOfUsersCV(ctxWithID, userID, cvID).Return(nil)

	_, err := cvUsecase.GenerateCVsPDF(ctxWithID, cvID)
	if err != nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}
}*/
