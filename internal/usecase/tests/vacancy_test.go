package usecase

import (
	"HnH/internal/domain"
	psqlmock "HnH/internal/repository/mock"
	"HnH/internal/usecase"
	"HnH/pkg/castUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"HnH/services/searchEngineService/searchEnginePB"
	"context"
	"fmt"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllVacanciesSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  16,
			VacancyName: "C++",
			Description: "Very cool vac",
		}}

	vacIDToFav := map[int]bool{
		10: true,
		11: true,
	}

	vacIDToLogo := map[int]string{
		10: "/avas/coolest/12.png",
		11: "",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetAllVacancies(ctxWithID).Return(vacList, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 10).Return([]string{"sql", "github"}, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 11).Return([]string{"boost", "RAII"}, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10, 11).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11).Return(vacIDToLogo, nil)

	_, err := vacUsecase.GetAllVacancies(ctxWithID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetAllVacanciesFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  16,
			VacancyName: "C++",
			Description: "Very cool vac",
		}}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetAllVacancies(ctxWithID).Return(vacList, serverErrors.INTERNAL_SERVER_ERROR)
	/*skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 10).Return([]string{"sql", "github"}, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 11).Return([]string{"boost", "RAII"}, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10, 11).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11).Return(vacIDToLogo, nil)*/

	_, err := vacUsecase.GetAllVacancies(ctxWithID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}
}

func TestGetAllVacanciesFail2(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  16,
			VacancyName: "C++",
			Description: "Very cool vac",
		}}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetAllVacancies(ctxWithID).Return(vacList, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 10).Return([]string{"sql", "github"}, serverErrors.INCORRECT_CREDENTIALS)
	/*skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 11).Return([]string{"boost", "RAII"}, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10, 11).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11).Return(vacIDToLogo, nil)*/

	_, err := vacUsecase.GetAllVacancies(ctxWithID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}
}

func TestGetAllVacanciesFail3(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  16,
			VacancyName: "C++",
			Description: "Very cool vac",
		}}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetAllVacancies(ctxWithID).Return(vacList, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 10).Return([]string{"sql", "github"}, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 11).Return([]string{"boost", "RAII"}, serverErrors.INCORRECT_CREDENTIALS)
	/*vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10, 11).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11).Return(vacIDToLogo, nil)*/

	_, err := vacUsecase.GetAllVacancies(ctxWithID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}
}

func TestGetAllVacanciesFail4(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  16,
			VacancyName: "C++",
			Description: "Very cool vac",
		}}

	vacIDToFav := map[int]bool{
		10: true,
		11: true,
	}

	vacIDToLogo := map[int]string{
		10: "/avas/coolest/12.png",
		11: "",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetAllVacancies(ctxWithID).Return(vacList, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 10).Return([]string{"sql", "github"}, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 11).Return([]string{"boost", "RAII"}, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10, 11).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11).Return(vacIDToLogo, serverErrors.INTERNAL_SERVER_ERROR)

	_, err := vacUsecase.GetAllVacancies(ctxWithID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}
}

func TestGetVacancySuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	dbVac := &domain.DbVacancy{
		ID:          10,
		EmployerID:  15,
		VacancyName: "Gopher",
		Description: "Cool vac",
	}

	vacIDToFav := map[int]bool{
		10: true,
	}

	vacIDToLogo := map[int]string{
		10: "/avas/coolest/12.png",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetVacancy(ctxWithID, 10).Return(dbVac, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 10).Return([]string{"sql", "github"}, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10).Return(vacIDToLogo, nil)

	_, err := vacUsecase.GetVacancy(ctxWithID, 10)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetVacancyWithCompanyNameSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	dbVac := &domain.DbVacancy{
		ID:          10,
		EmployerID:  15,
		VacancyName: "Gopher",
		Description: "Cool vac",
	}

	vacIDToFav := map[int]bool{
		10: true,
	}

	vacIDToLogo := map[int]string{
		10: "/avas/coolest/12.png",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetVacancy(ctxWithID, 10).Return(dbVac, nil)
	vacancyRepo.EXPECT().GetCompanyName(ctxWithID, 10).Return("VK", nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 10).Return([]string{"sql", "github"}, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10).Return(vacIDToLogo, nil).Times(2)

	_, err := vacUsecase.GetVacancyWithCompanyName(ctxWithID, 10)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestAddVacancySuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	apiVac := &domain.ApiVacancy{
		ID:          10,
		EmployerID:  15,
		VacancyName: "Gopher",
		Description: "Cool vac",
	}

	dbVac := &domain.DbVacancy{
		ID:          10,
		EmployerID:  0,
		VacancyName: "Gopher",
		Description: "Cool vac",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().AddVacancy(ctxWithID, 15, dbVac).Return(10, nil)
	skillRepo.EXPECT().AddSkillsByVacID(ctxWithID, 10, gomock.Any()).Return(nil)
	userRepo.EXPECT().GetRoleById(ctxWithID, 12).Return(domain.Employer, nil)
	userRepo.EXPECT().GetUserEmpId(ctxWithID, 12).Return(15, nil)

	_, err := vacUsecase.AddVacancy(ctxWithID, apiVac)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestUpdateVacancySuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	apiVac := &domain.ApiVacancy{
		ID:          10,
		EmployerID:  15,
		VacancyName: "Gopher",
		Description: "Cool vac",
	}

	dbVac := &domain.DbVacancy{
		ID:          10,
		EmployerID:  0,
		VacancyName: "Gopher",
		Description: "Cool vac",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().UpdateEmpVacancy(ctxWithID, 15, 10, dbVac).Return(nil)
	userRepo.EXPECT().GetRoleById(ctxWithID, 12).Return(domain.Employer, nil)
	userRepo.EXPECT().GetUserEmpId(ctxWithID, 12).Return(15, nil)
	vacancyRepo.EXPECT().GetEmpId(ctxWithID, 10).Return(15, nil)

	err := vacUsecase.UpdateVacancy(ctxWithID, 10, apiVac)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestDeleteVacancySuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().DeleteEmpVacancy(ctxWithID, 15, 10).Return(nil)
	userRepo.EXPECT().GetRoleById(ctxWithID, 12).Return(domain.Employer, nil)
	userRepo.EXPECT().GetUserEmpId(ctxWithID, 12).Return(15, nil)
	vacancyRepo.EXPECT().GetEmpId(ctxWithID, 10).Return(15, nil)

	err := vacUsecase.DeleteVacancy(ctxWithID, 10)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetUserVacanciesSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  15,
			VacancyName: "C++",
			Description: "Very cool vac",
		}}

	vacIDToFav := map[int]bool{
		10: true,
		11: true,
	}

	vacIDToLogo := map[int]string{
		10: "/avas/coolest/12.png",
		11: "",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetUserVacancies(ctxWithID, 12).Return(vacList, nil)
	userRepo.EXPECT().GetRoleById(ctxWithID, 12).Return(domain.Employer, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10, 11).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11).Return(vacIDToLogo, nil)

	_, err := vacUsecase.GetUserVacancies(ctxWithID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetEmployerInfoSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  15,
			VacancyName: "C++",
			Description: "Very cool vac",
		}}

	vacIDToFav := map[int]bool{
		10: true,
		11: true,
	}

	vacIDToLogo := map[int]string{
		10: "/avas/coolest/12.png",
		11: "",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetEmployerInfo(ctxWithID, 15).Return("Vladimir", "Borozenets", "VK", vacList, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10, 11).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11).Return(vacIDToLogo, nil)

	_, err := vacUsecase.GetEmployerInfo(ctxWithID, 15)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestSearchVacanciesSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  15,
			VacancyName: "C++",
			Description: "Very cool vac",
		},
		{
			ID:          12,
			EmployerID:  15,
			VacancyName: "Python",
			Description: "Not so cool vac",
		}}

	vacIDToFav := map[int]bool{
		10: true,
		11: true,
		12: false,
	}

	vacIDToLogo := map[int]string{
		10: "/avas/coolest/12.png",
		11: "",
		12: "/avas/notsocoolest/14.png",
	}

	opts := &searchEnginePB.SearchOptions{}
	resp := &searchEnginePB.SearchResponse{
		Ids:   []int64{10, 11, 12},
		Count: 3,
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	searchRepo.EXPECT().SearchVacancyIDs(ctxWithID, opts).Return(resp, nil)
	vacancyRepo.EXPECT().GetVacanciesByIds(ctxWithID, castUtils.Int64SliceToIntSlice(resp.Ids)).Return(vacList, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 10).Return([]string{"golang", "git"}, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 11).Return([]string{"sql", "docker"}, nil)
	skillRepo.EXPECT().GetSkillsByVacID(ctxWithID, 12).Return([]string{"pytorch"}, nil)
	vacancyRepo.EXPECT().GetFavouriteFlags(ctxWithID, userID, 10, 11, 12).Return(vacIDToFav, nil).AnyTimes()
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11, 12).Return(vacIDToLogo, nil)

	_, err := vacUsecase.SearchVacancies(ctxWithID, opts)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestAddToFavouriteSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12
	vacID := 15

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().AddToFavourite(ctxWithID, userID, vacID).Return(nil)

	err := vacUsecase.AddToFavourite(ctxWithID, vacID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestAddToFavouriteFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12
	vacID := 15

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().AddToFavourite(ctxWithID, userID, vacID).Return(serverErrors.FORBIDDEN)

	err := vacUsecase.AddToFavourite(ctxWithID, vacID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.FORBIDDEN.Error())
}

func TestDeleteFromFavouriteSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12
	vacID := 15

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().DeleteFromFavourite(ctxWithID, userID, vacID).Return(nil)

	err := vacUsecase.DeleteFromFavourite(ctxWithID, vacID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestDeleteFromFavouriteFail(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12
	vacID := 15

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().DeleteFromFavourite(ctxWithID, userID, vacID).Return(serverErrors.FORBIDDEN)

	err := vacUsecase.DeleteFromFavourite(ctxWithID, vacID)
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.FORBIDDEN.Error())
}

func TestGetFacouriteSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	userRepo := psqlmock.NewMockIUserRepository(mockCtrl)
	sessionRepo := psqlmock.NewMockIAuthRepository(mockCtrl)
	vacancyRepo := psqlmock.NewMockIVacancyRepository(mockCtrl)
	skillRepo := psqlmock.NewMockISkillRepository(mockCtrl)
	searchRepo := psqlmock.NewMockISearchEngineRepository(mockCtrl)

	vacUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo, searchRepo, skillRepo)

	defer mockCtrl.Finish()

	userID := 12

	vacList := []domain.DbVacancy{
		{
			ID:          10,
			EmployerID:  15,
			VacancyName: "Gopher",
			Description: "Cool vac",
		},
		{
			ID:          11,
			EmployerID:  15,
			VacancyName: "C++",
			Description: "Very cool vac",
		},
		{
			ID:          12,
			EmployerID:  15,
			VacancyName: "Python",
			Description: "Not so cool vac",
		}}

	vacIDToLogo := map[int]string{
		10: "/avas/coolest/12.png",
		11: "",
		12: "/avas/notsocoolest/14.png",
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	vacancyRepo.EXPECT().GetFavourite(ctxWithID, userID).Return(vacList, nil)
	userRepo.EXPECT().GetLogoPathesByVacancyIDList(ctxWithID, 10, 11, 12).Return(vacIDToLogo, nil)

	_, err := vacUsecase.GetFavourite(ctxWithID)
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}
