package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/pkg/castUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"HnH/services/searchEngineService/searchEnginePB"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type IVacancyUsecase interface {
	GetAllVacancies(ctx context.Context) ([]domain.ApiVacancy, error)
	GetVacancy(ctx context.Context, vacancyID int) (*domain.ApiVacancy, error)
	GetVacancyWithCompanyName(ctx context.Context, vacancyID int) (*domain.CompanyVacancy, error)
	GetUserVacancies(ctx context.Context) ([]domain.ApiVacancy, error)
	GetEmployerInfo(ctx context.Context, employerID int) (*domain.EmployerInfo, error)
	AddVacancy(ctx context.Context, vacancy *domain.ApiVacancy) (int, error)
	UpdateVacancy(ctx context.Context, vacancyID int, vacancy *domain.ApiVacancy) error
	DeleteVacancy(ctx context.Context, vacancyID int) error
	SearchVacancies(ctx context.Context, options *searchEnginePB.SearchOptions) (domain.ApiMetaVacancy, error)
	AddToFavourite(ctx context.Context, vacancyID int) error
	DeleteFromFavourite(ctx context.Context, vacancyID int) error
	GetFavourite(ctx context.Context) ([]domain.ApiVacancy, error)
}

type VacancyUsecase struct {
	vacancyRepo      psql.IVacancyRepository
	sessionRepo      grpc.IAuthRepository
	userRepo         psql.IUserRepository
	searchEngineRepo grpc.ISearchEngineRepository
	skillRepo        psql.ISkillRepository
	statRepo         psql.IStatRepository
}

func NewVacancyUsecase(
	vacancyRepository psql.IVacancyRepository,
	sessionRepository grpc.IAuthRepository,
	userRepository psql.IUserRepository,
	searchEngineRepository grpc.ISearchEngineRepository,
	skillRepository psql.ISkillRepository,
	statRepo psql.IStatRepository,
) IVacancyUsecase {
	return &VacancyUsecase{
		vacancyRepo:      vacancyRepository,
		sessionRepo:      sessionRepository,
		userRepo:         userRepository,
		searchEngineRepo: searchEngineRepository,
		skillRepo:        skillRepository,
		statRepo:         statRepo,
	}
}

func (vacancyUsecase *VacancyUsecase) validateEmployerAndGetEmpId(ctx context.Context, userID int) (int, error) {
	userRole, err := vacancyUsecase.userRepo.GetRoleById(ctx, userID)
	if err != nil {
		return 0, err
	} else if userRole != domain.Employer {
		return 0, ErrInapropriateRole
	}

	userEmpID, err := vacancyUsecase.userRepo.GetUserEmpId(ctx, userID)
	if err != nil {
		return 0, err
	}

	return userEmpID, nil
}

func (vacancyUsecase *VacancyUsecase) validateEmployer(ctx context.Context, userID int, vacancyID int) (int, error) {
	userEmpID, validStatus := vacancyUsecase.validateEmployerAndGetEmpId(ctx, userID)
	if validStatus != nil {
		return 0, validStatus
	}

	empID, err := vacancyUsecase.vacancyRepo.GetEmpId(ctx, vacancyID)
	if err != nil {
		return 0, err
	}

	if userEmpID != empID {
		return 0, serverErrors.FORBIDDEN
	}

	return userEmpID, nil
}

func (vacancyUsecase *VacancyUsecase) collectApiVacs(vacs []domain.DbVacancy) []domain.ApiVacancy {
	res := []domain.ApiVacancy{}
	for _, vac := range vacs {
		res = append(res, *vac.ToAPI())
	}
	return res
}

func (vacancyUsecase *VacancyUsecase) setFavouriteFlags(ctx context.Context, vacs ...domain.ApiVacancy) ([]domain.ApiVacancy, error) {
	userID, loggedIn := contextUtils.IsLoggedIn(ctx)
	if !loggedIn {
		return vacs, nil
	}

	vacIDs := []int{}
	for _, vac := range vacs {
		vacIDs = append(vacIDs, vac.ID)
	}

	vacIDToFav, err := vacancyUsecase.vacancyRepo.GetFavouriteFlags(ctx, userID, vacIDs...)
	if err != nil {
		return nil, err
	}

	vacsToReturn := []domain.ApiVacancy{}
	for _, vac := range vacs {
		isFav, found := vacIDToFav[vac.ID]
		if found {
			vac.Favourite = isFav
		} else {
			vac.Favourite = false
		}

		vacsToReturn = append(vacsToReturn, vac)
	}

	return vacsToReturn, nil
}

func (vacancyUsecase *VacancyUsecase) setLogoPath(ctx context.Context, vacs ...domain.ApiVacancy) ([]domain.ApiVacancy, error) {
	vacIDs := []int{}

	for _, vac := range vacs {
		vacIDs = append(vacIDs, vac.ID)
	}

	vacIDToPath, err := vacancyUsecase.userRepo.GetLogoPathesByVacancyIDList(ctx, vacIDs...)
	if err != nil {
		return nil, err
	}

	vacsToReturn := []domain.ApiVacancy{}
	for _, vac := range vacs {
		path, found := vacIDToPath[vac.ID]
		if !found || path == "" {
			vac.LogoURL = ""
		} else {
			vac.LogoURL = "/image" + path
		}

		vacsToReturn = append(vacsToReturn, vac)
	}

	return vacsToReturn, nil
}

func (vacancyUsecase *VacancyUsecase) GetAllVacancies(ctx context.Context) ([]domain.ApiVacancy, error) {
	vacancies, getErr := vacancyUsecase.vacancyRepo.GetAllVacancies(ctx)
	if getErr != nil {
		return nil, getErr
	}

	// TODO: optimize
	apiVacs := vacancyUsecase.collectApiVacs(vacancies)
	for i := range apiVacs {
		skills, err := vacancyUsecase.skillRepo.GetSkillsByVacID(ctx, apiVacs[i].ID)
		if err != nil {
			return nil, err
		}
		apiVacs[i].Skills = skills
	}

	apiVacs, err := vacancyUsecase.setFavouriteFlags(ctx, apiVacs...)
	if err != nil {
		return nil, err
	}

	apiVacs, err = vacancyUsecase.setLogoPath(ctx, apiVacs...)
	if err != nil {
		return nil, err
	}

	return apiVacs, nil
}

func (vacancyUsecase *VacancyUsecase) GetVacancy(ctx context.Context, vacancyID int) (*domain.ApiVacancy, error) {
	sesssionID, err := contextUtils.GetSessionIDFromCtx(ctx)
	if err == nil {
		userID, err := vacancyUsecase.sessionRepo.GetUserIdBySession(ctx, sesssionID)
		if err != nil {
			return nil, err
		}
		applicantID, err := vacancyUsecase.userRepo.GetUserAppId(ctx, userID)
		if !errors.Is(err, psql.ErrEntityNotFound) {
			if err != nil {
				return nil, err
			} else {
				err = vacancyUsecase.statRepo.AddVacancyView(ctx, vacancyID, applicantID)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	vacancy, err := vacancyUsecase.vacancyRepo.GetVacancy(ctx, vacancyID)
	if err != nil {
		return nil, err
	}

	apiVac := vacancy.ToAPI()
	skills, err := vacancyUsecase.skillRepo.GetSkillsByVacID(ctx, apiVac.ID)
	if err != nil {
		return nil, err
	}
	apiVac.Skills = skills

	viewsCount, err := vacancyUsecase.statRepo.CountVacancyViews(ctx, vacancyID)
	if err != nil {
		return nil, err
	}
	apiVac.ViewsCount = viewsCount

	responsesCount, err := vacancyUsecase.statRepo.CountVacancyResponses(ctx, vacancyID)
	if err != nil {
		return nil, err
	}
	apiVac.ResponsesCount = responsesCount

	vacToReturn, err := vacancyUsecase.setFavouriteFlags(ctx, *apiVac)
	if err != nil {
		return nil, err
	}

	vacToReturn, err = vacancyUsecase.setLogoPath(ctx, vacToReturn...)
	if err != nil {
		return nil, err
	}

	return &vacToReturn[0], nil
}

func (vacancyUsecase *VacancyUsecase) GetVacancyWithCompanyName(ctx context.Context, vacancyID int) (*domain.CompanyVacancy, error) {
	vacancy, err := vacancyUsecase.GetVacancy(ctx, vacancyID)
	if err != nil {
		return nil, err
	}

	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy": vacancy,
	}).
		Debug("got vacancy")

	companyName, err := vacancyUsecase.vacancyRepo.GetCompanyName(ctx, vacancyID)
	if err != nil {
		return nil, err
	}

	vacToReturn, err := vacancyUsecase.setFavouriteFlags(ctx, *vacancy)
	if err != nil {
		return nil, err
	}

	viewsCount, err := vacancyUsecase.statRepo.CountVacancyViews(ctx, vacancyID)
	if err != nil {
		return nil, err
	}
	vacancy.ViewsCount = viewsCount

	responsesCount, err := vacancyUsecase.statRepo.CountVacancyResponses(ctx, vacancyID)
	if err != nil {
		return nil, err
	}
	vacancy.ResponsesCount = responsesCount

	vacToReturn, err = vacancyUsecase.setLogoPath(ctx, vacToReturn...)
	if err != nil {
		return nil, err
	}

	compVac := &domain.CompanyVacancy{
		CompanyName: companyName,
		Vacancy:     vacToReturn[0],
	}

	return compVac, nil
}

func (vacancyUsecase *VacancyUsecase) AddVacancy(ctx context.Context, vacancy *domain.ApiVacancy) (int, error) {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	userEmpID, validStatus := vacancyUsecase.validateEmployerAndGetEmpId(ctx, userID)
	if validStatus != nil {
		return 0, validStatus
	}

	dbVac := vacancy.ToDb()
	vacancyID, addStatus := vacancyUsecase.vacancyRepo.AddVacancy(ctx, userEmpID, dbVac)
	if addStatus != nil {
		return 0, addStatus
	}

	addSkillsErr := vacancyUsecase.skillRepo.AddSkillsByVacID(ctx, vacancyID, vacancy.Skills)
	if addSkillsErr != nil {
		return 0, addSkillsErr
	}

	return vacancyID, nil
}

// TODO: add skills
func (vacancyUsecase *VacancyUsecase) UpdateVacancy(ctx context.Context, vacancyID int, vacancy *domain.ApiVacancy) error {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	empID, validStatus := vacancyUsecase.validateEmployer(ctx, userID, vacancyID)
	if validStatus != nil {
		return validStatus
	}

	updStatus := vacancyUsecase.vacancyRepo.UpdateEmpVacancy(ctx, empID, vacancyID, vacancy.ToDb())
	if updStatus != nil {
		return updStatus
	}
	return nil
}

func (vacancyUsecase *VacancyUsecase) DeleteVacancy(ctx context.Context, vacancyID int) error {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	empID, validStatus := vacancyUsecase.validateEmployer(ctx, userID, vacancyID)
	if validStatus != nil {
		return validStatus
	}

	delStatus := vacancyUsecase.vacancyRepo.DeleteEmpVacancy(ctx, empID, vacancyID)
	if delStatus != nil {
		return delStatus
	}
	return nil
}

func (vacancyUsecase *VacancyUsecase) GetUserVacancies(ctx context.Context) ([]domain.ApiVacancy, error) {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	role, err := vacancyUsecase.userRepo.GetRoleById(ctx, userID)
	if err != nil {
		return nil, err
	}

	if role != domain.Employer {
		return nil, ErrInapropriateRole
	}

	vacanciesList, err := vacancyUsecase.vacancyRepo.GetUserVacancies(ctx, userID)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("vacancies: %v\n", vacanciesList)

	apiVacs := vacancyUsecase.collectApiVacs(vacanciesList)

	for i := range apiVacs {
		viewsCount, err := vacancyUsecase.statRepo.CountVacancyViews(ctx, apiVacs[i].ID)
		if err != nil {
			return nil, err
		}
		apiVacs[i].ViewsCount = viewsCount

		responsesCount, err := vacancyUsecase.statRepo.CountVacancyResponses(ctx, apiVacs[i].ID)
		if err != nil {
			return nil, err
		}
		apiVacs[i].ResponsesCount = responsesCount
	}

	apiVacs, err = vacancyUsecase.setFavouriteFlags(ctx, apiVacs...)
	if err != nil {
		return nil, err
	}

	apiVacs, err = vacancyUsecase.setLogoPath(ctx, apiVacs...)
	if err != nil {
		return nil, err
	}

	return apiVacs, nil
}

func (vacancyUsecase *VacancyUsecase) GetEmployerInfo(ctx context.Context, employerID int) (*domain.EmployerInfo, error) {
	first_name, last_name, compName, empVacs, err := vacancyUsecase.vacancyRepo.GetEmployerInfo(ctx, employerID)
	if err != nil {
		return nil, err
	}

	vacsToReturn := vacancyUsecase.collectApiVacs(empVacs)

	for i := range vacsToReturn {
		viewsCount, err := vacancyUsecase.statRepo.CountVacancyViews(ctx, vacsToReturn[i].ID)
		if err != nil {
			return nil, err
		}
		vacsToReturn[i].ViewsCount = viewsCount

		responsesCount, err := vacancyUsecase.statRepo.CountVacancyResponses(ctx, vacsToReturn[i].ID)
		if err != nil {
			return nil, err
		}
		vacsToReturn[i].ResponsesCount = responsesCount
	}

	info := &domain.EmployerInfo{
		FirstName:   first_name,
		LastName:    last_name,
		CompanyName: compName,
		Vacancies:   vacsToReturn,
	}

	info.Vacancies, err = vacancyUsecase.setFavouriteFlags(ctx, info.Vacancies...)
	if err != nil {
		return nil, err
	}

	info.Vacancies, err = vacancyUsecase.setLogoPath(ctx, info.Vacancies...)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (vacancyUsecase *VacancyUsecase) SearchVacancies(
	ctx context.Context,
	options *searchEnginePB.SearchOptions,
) (domain.ApiMetaVacancy, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	vacanciesSearchResponse, err := vacancyUsecase.searchEngineRepo.SearchVacancyIDs(ctx, options)
	if err != nil {
		return domain.ApiMetaVacancy{
			Filters:   nil,
			Vacancies: domain.ApiVacancyCount{},
		}, err
	}
	vacancyIDs, count := vacanciesSearchResponse.Ids, vacanciesSearchResponse.Count

	vacancies, vacErr := vacancyUsecase.vacancyRepo.GetVacanciesByIds(ctx, castUtils.Int64SliceToIntSlice(vacancyIDs))
	if vacErr == psql.ErrEntityNotFound {
		return domain.ApiMetaVacancy{
			Filters:   nil,
			Vacancies: domain.ApiVacancyCount{},
		}, nil
	}
	if vacErr != nil {
		return domain.ApiMetaVacancy{}, vacErr
	}

	vacanciesToReturn := vacancyUsecase.collectApiVacs(vacancies)

	for i := range vacanciesToReturn {
		contextLogger.WithFields(logrus.Fields{
			"i":                       i,
			"vacanciesToReturn[i]":    vacanciesToReturn[i],
			"vacanciesToReturn[i].ID": vacanciesToReturn[i].ID,
		}).
			Debug("about to count vacancy views")
		viewsCount, err := vacancyUsecase.statRepo.CountVacancyViews(ctx, vacanciesToReturn[i].ID)
		if err != nil {
			return domain.ApiMetaVacancy{}, err
		}
		vacanciesToReturn[i].ViewsCount = viewsCount
		contextLogger.WithFields(logrus.Fields{
			"views_count": viewsCount,
		}).
			Debug("got views count")

		responsesCount, err := vacancyUsecase.statRepo.CountVacancyResponses(ctx, vacanciesToReturn[i].ID)
		if err != nil {
			return domain.ApiMetaVacancy{}, err
		}
		vacanciesToReturn[i].ResponsesCount = responsesCount
		contextLogger.WithFields(logrus.Fields{
			"response_count": responsesCount,
		}).
			Debug("got response count")
	}

	// TODO: optimize
	for i := range vacanciesToReturn {
		skills, err := vacancyUsecase.skillRepo.GetSkillsByVacID(ctx, vacanciesToReturn[i].ID)
		if err != nil {
			return domain.ApiMetaVacancy{}, err
		}
		vacanciesToReturn[i].Skills = skills
	}

	vacanciesToReturn, err = vacancyUsecase.setFavouriteFlags(ctx, vacanciesToReturn...)
	if err != nil {
		return domain.ApiMetaVacancy{
			Filters:   nil,
			Vacancies: domain.ApiVacancyCount{},
		}, err
	}

	vacanciesToReturn, err = vacancyUsecase.setLogoPath(ctx, vacanciesToReturn...)
	if err != nil {
		return domain.ApiMetaVacancy{
			Filters:   nil,
			Vacancies: domain.ApiVacancyCount{},
		}, err
	}

	result := domain.ApiMetaVacancy{
		Filters: vacanciesSearchResponse.Filters,
		Vacancies: domain.ApiVacancyCount{
			Count:     count,
			Vacancies: vacanciesToReturn,
		},
	}

	return result, nil
}

func (vacancyUsecase *VacancyUsecase) AddToFavourite(ctx context.Context, vacancyID int) error {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	err := vacancyUsecase.vacancyRepo.AddToFavourite(ctx, userID, vacancyID)
	if err != nil {
		return err
	}

	return nil
}

func (vacancyUsecase *VacancyUsecase) DeleteFromFavourite(ctx context.Context, vacancyID int) error {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	err := vacancyUsecase.vacancyRepo.DeleteFromFavourite(ctx, userID, vacancyID)
	if err != nil {
		return err
	}

	return nil
}

func (vacancyUsecase *VacancyUsecase) GetFavourite(ctx context.Context) ([]domain.ApiVacancy, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	userID := contextUtils.GetUserIDFromCtx(ctx)

	favVacs, err := vacancyUsecase.vacancyRepo.GetFavourite(ctx, userID)
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"err_msg": err,
			"user_id": userID,
		}).
			Error("error while getting favourite vacancies")
		return nil, err
	}

	apiVacs := vacancyUsecase.collectApiVacs(favVacs)
	for i := range apiVacs {
		apiVacs[i].Favourite = true
	}

	apiVacs, err = vacancyUsecase.setLogoPath(ctx, apiVacs...)
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"err_msg": err,
			"user_id": userID,
		}).
			Error("error while getting favourite vacancies")
		return nil, err
	}

	return apiVacs, nil
}
