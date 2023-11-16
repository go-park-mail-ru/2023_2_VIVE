package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/serverErrors"
)

type IVacancyUsecase interface {
	GetAllVacancies() ([]domain.ApiVacancy, error)
	GetVacancy(vacancyID int) (*domain.ApiVacancy, error)
	GetUserVacancies(sessionID string) ([]domain.ApiVacancy, error)
	AddVacancy(sessionID string, vacancy *domain.DbVacancy) (int, error)
	UpdateVacancy(sessionID string, vacancyID int, vacancy *domain.ApiVacancy) error
	DeleteVacancy(sessionID string, vacancyID int) error
}

type VacancyUsecase struct {
	vacancyRepo psql.IVacancyRepository
	sessionRepo redisRepo.ISessionRepository
	userRepo    psql.IUserRepository
}

func NewVacancyUsecase(vacancyRepository psql.IVacancyRepository,
	sessionRepository redisRepo.ISessionRepository,
	userRepository psql.IUserRepository) IVacancyUsecase {
	return &VacancyUsecase{
		vacancyRepo: vacancyRepository,
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
	}
}

func (vacancyUsecase *VacancyUsecase) validateEmployerAndGetOrgId(sessionID string) (int, error) {
	validStatus := vacancyUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userID, err := vacancyUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return 0, err
	}

	userRole, err := vacancyUsecase.userRepo.GetRoleById(userID)
	if err != nil {
		return 0, err
	} else if userRole != domain.Employer {
		return 0, INAPPROPRIATE_ROLE
	}

	userOrgID, err := vacancyUsecase.userRepo.GetUserOrgId(userID)
	if err != nil {
		return 0, err
	}

	return userOrgID, nil
}

func (vacancyUsecase *VacancyUsecase) validateEmployer(sessionID string, vacancyID int) (int, error) {
	userOrgID, validStatus := vacancyUsecase.validateEmployerAndGetOrgId(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	orgID, err := vacancyUsecase.vacancyRepo.GetOrgId(vacancyID)
	if err != nil {
		return 0, err
	}

	if userOrgID != orgID {
		return 0, serverErrors.FORBIDDEN
	}

	return userOrgID, nil
}

func (vacancyUsecase *VacancyUsecase) collectApiVacs(vacs []domain.DbVacancy) []domain.ApiVacancy {
	res := []domain.ApiVacancy{}
	for _, vac := range vacs {
		res = append(res, *vac.ToAPI())
	}
	return res
}

func (vacancyUsecase *VacancyUsecase) GetAllVacancies() ([]domain.ApiVacancy, error) {
	vacancies, getErr := vacancyUsecase.vacancyRepo.GetAllVacancies()
	if getErr != nil {
		return nil, getErr
	}

	return vacancyUsecase.collectApiVacs(vacancies), nil
}

func (vacancyUsecase *VacancyUsecase) GetVacancy(vacancyID int) (*domain.ApiVacancy, error) {
	vacancy, err := vacancyUsecase.vacancyRepo.GetVacancy(vacancyID)
	if err != nil {
		return nil, err
	}

	return vacancy.ToAPI(), nil
}

func (vacancyUsecase *VacancyUsecase) AddVacancy(sessionID string, vacancy *domain.DbVacancy) (int, error) {
	userOrgID, validStatus := vacancyUsecase.validateEmployerAndGetOrgId(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	// fmt.Printf("vacancy: %v\n", vacancy)

	vacancyID, addStatus := vacancyUsecase.vacancyRepo.AddVacancy(userOrgID, vacancy)
	if addStatus != nil {
		return 0, addStatus
	}

	return vacancyID, nil
}

func (vacancyUsecase *VacancyUsecase) UpdateVacancy(sessionID string, vacancyID int, vacancy *domain.ApiVacancy) error {
	orgID, validStatus := vacancyUsecase.validateEmployer(sessionID, vacancyID)
	if validStatus != nil {
		return validStatus
	}

	updStatus := vacancyUsecase.vacancyRepo.UpdateOrgVacancy(orgID, vacancyID, vacancy.ToDb())
	if updStatus != nil {
		return updStatus
	}
	return nil
}

func (vacancyUsecase *VacancyUsecase) DeleteVacancy(sessionID string, vacancyID int) error {
	orgID, validStatus := vacancyUsecase.validateEmployer(sessionID, vacancyID)
	if validStatus != nil {
		return validStatus
	}

	delStatus := vacancyUsecase.vacancyRepo.DeleteOrgVacancy(orgID, vacancyID)
	if delStatus != nil {
		return delStatus
	}
	return nil
}

func (vacancyUsecase *VacancyUsecase) GetUserVacancies(sessionID string) ([]domain.ApiVacancy, error) {
	userID, err := vacancyUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return nil, serverErrors.AUTH_REQUIRED
	}

	role, err := vacancyUsecase.userRepo.GetRoleById(userID)
	if err != nil {
		return nil, err
	}

	if role != domain.Employer {
		return nil, INAPPROPRIATE_ROLE
	}

	vacanciesList, err := vacancyUsecase.vacancyRepo.GetUserVacancies(userID)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("vacancies: %v\n", vacanciesList)

	return vacancyUsecase.collectApiVacs(vacanciesList), nil
}
