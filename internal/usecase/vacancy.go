package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/castUtils"
	"HnH/pkg/serverErrors"
	"context"
)

type IVacancyUsecase interface {
	GetAllVacancies(ctx context.Context) ([]domain.ApiVacancy, error)
	GetVacancy(ctx context.Context, vacancyID int) (*domain.ApiVacancy, error)
	GetUserVacancies(ctx context.Context, sessionID string) ([]domain.ApiVacancy, error)
	AddVacancy(ctx context.Context, sessionID string, vacancy *domain.DbVacancy) (int, error)
	UpdateVacancy(ctx context.Context, sessionID string, vacancyID int, vacancy *domain.ApiVacancy) error
	DeleteVacancy(ctx context.Context, sessionID string, vacancyID int) error
	SearchVacancies(ctx context.Context, query string, pageNumber, resultsPerPage int64) (domain.ApiMetaVacancy, error)
}

type VacancyUsecase struct {
	vacancyRepo      psql.IVacancyRepository
	sessionRepo      redisRepo.ISessionRepository
	userRepo         psql.IUserRepository
	searchEngineRepo grpc.ISearchEngineRepository
}

func NewVacancyUsecase(
	vacancyRepository psql.IVacancyRepository,
	sessionRepository redisRepo.ISessionRepository,
	userRepository psql.IUserRepository,
	searchEngineRepository grpc.ISearchEngineRepository,
) IVacancyUsecase {
	return &VacancyUsecase{
		vacancyRepo:      vacancyRepository,
		sessionRepo:      sessionRepository,
		userRepo:         userRepository,
		searchEngineRepo: searchEngineRepository,
	}
}

func (vacancyUsecase *VacancyUsecase) validateEmployerAndGetEmpId(ctx context.Context, sessionID string) (int, error) {
	validStatus := vacancyUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userID, err := vacancyUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return 0, err
	}

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

func (vacancyUsecase *VacancyUsecase) validateEmployer(ctx context.Context, sessionID string, vacancyID int) (int, error) {
	userEmpID, validStatus := vacancyUsecase.validateEmployerAndGetEmpId(ctx, sessionID)
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

func (vacancyUsecase *VacancyUsecase) GetAllVacancies(ctx context.Context) ([]domain.ApiVacancy, error) {
	vacancies, getErr := vacancyUsecase.vacancyRepo.GetAllVacancies(ctx)
	if getErr != nil {
		return nil, getErr
	}

	return vacancyUsecase.collectApiVacs(vacancies), nil
}

func (vacancyUsecase *VacancyUsecase) GetVacancy(ctx context.Context, vacancyID int) (*domain.ApiVacancy, error) {
	vacancy, err := vacancyUsecase.vacancyRepo.GetVacancy(ctx, vacancyID)
	if err != nil {
		return nil, err
	}

	return vacancy.ToAPI(), nil
}

func (vacancyUsecase *VacancyUsecase) AddVacancy(ctx context.Context, sessionID string, vacancy *domain.DbVacancy) (int, error) {
	userOrgID, validStatus := vacancyUsecase.validateEmployerAndGetEmpId(ctx, sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	// fmt.Printf("vacancy: %v\n", vacancy)

	vacancyID, addStatus := vacancyUsecase.vacancyRepo.AddVacancy(ctx, userOrgID, vacancy)
	if addStatus != nil {
		return 0, addStatus
	}

	return vacancyID, nil
}

func (vacancyUsecase *VacancyUsecase) UpdateVacancy(ctx context.Context, sessionID string, vacancyID int, vacancy *domain.ApiVacancy) error {
	empID, validStatus := vacancyUsecase.validateEmployer(ctx, sessionID, vacancyID)
	if validStatus != nil {
		return validStatus
	}

	updStatus := vacancyUsecase.vacancyRepo.UpdateEmpVacancy(ctx, empID, vacancyID, vacancy.ToDb())
	if updStatus != nil {
		return updStatus
	}
	return nil
}

func (vacancyUsecase *VacancyUsecase) DeleteVacancy(ctx context.Context, sessionID string, vacancyID int) error {
	empID, validStatus := vacancyUsecase.validateEmployer(ctx, sessionID, vacancyID)
	if validStatus != nil {
		return validStatus
	}

	delStatus := vacancyUsecase.vacancyRepo.DeleteEmpVacancy(ctx, empID, vacancyID)
	if delStatus != nil {
		return delStatus
	}
	return nil
}

func (vacancyUsecase *VacancyUsecase) GetUserVacancies(ctx context.Context, sessionID string) ([]domain.ApiVacancy, error) {
	userID, err := vacancyUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return nil, serverErrors.AUTH_REQUIRED
	}

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

	return vacancyUsecase.collectApiVacs(vacanciesList), nil
}

func (vacancyUsecase *VacancyUsecase) SearchVacancies(
	ctx context.Context,
	query string,
	pageNumber, resultsPerPage int64,
) (domain.ApiMetaVacancy, error) {
	vacancyIDs, count, err := vacancyUsecase.searchEngineRepo.SearchVacancyIDs(ctx, query, pageNumber, resultsPerPage)
	if err != nil {
		return domain.ApiMetaVacancy{
			Count:     0,
			Vacancies: nil,
		}, nil
	}

	vacancies, vacErr := vacancyUsecase.vacancyRepo.GetVacanciesByIds(ctx, castUtils.Int64SliceToIntSlice(vacancyIDs))
	if vacErr == psql.ErrEntityNotFound {
		return domain.ApiMetaVacancy{
			Count:     0,
			Vacancies: nil,
		}, nil
	}
	if vacErr != nil {
		return domain.ApiMetaVacancy{}, vacErr
	}

	vacanciesToReturn := vacancyUsecase.collectApiVacs(vacancies)
	result := domain.ApiMetaVacancy{
		Count:     count,
		Vacancies: vacanciesToReturn,
	}

	return result, nil
}
