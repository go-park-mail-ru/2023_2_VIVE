package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"
	"HnH/pkg/serverErrors"
)

type IVacancyUsecase interface {
	GetAllVacancies() ([]domain.Vacancy, error)

	GetVacancy(vacancyID int) (*domain.Vacancy, error)
	AddVacancy(sessionID string, vacancy *domain.Vacancy) (int, error)
	UpdateVacancy(sessionID string, vacancyID int, vacancy *domain.Vacancy) error
	DeleteVacancy(sessionID string, vacancyID int) error
}

type VacancyUsecase struct {
	vacancyRepo repository.IVacancyRepository
	sessionRepo repository.ISessionRepository
	userRepo    repository.IUserRepository
}

func NewVacancyUsecase(vacancyRepository repository.IVacancyRepository,
	sessionRepository repository.ISessionRepository,
	userRepository repository.IUserRepository) IVacancyUsecase {
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

func (vacancyUsecase *VacancyUsecase) validateEmployer(sessionID string, vacancyID int) error {
	userOrgID, validStatus := vacancyUsecase.validateEmployerAndGetOrgId(sessionID)
	if validStatus != nil {
		return validStatus
	}

	orgID, err := vacancyUsecase.vacancyRepo.GetOrgId(vacancyID)
	if err != nil {
		return err
	}

	if userOrgID != orgID {
		return serverErrors.FORBIDDEN
	}

	return nil
}

func (vacancyUsecase *VacancyUsecase) GetAllVacancies() ([]domain.Vacancy, error) {
	vacancies, getErr := vacancyUsecase.vacancyRepo.GetAllVacancies()
	if getErr != nil {
		return nil, getErr
	}

	return vacancies, nil
}

func (vacancyUsecase *VacancyUsecase) GetVacancy(vacancyID int) (*domain.Vacancy, error) {
	vacancy, err := vacancyUsecase.vacancyRepo.GetVacancy(vacancyID)
	if err != nil {
		return nil, err
	}

	return vacancy, nil
}

func (vacancyUsecase *VacancyUsecase) AddVacancy(sessionID string, vacancy *domain.Vacancy) (int, error) {
	userOrgID, validStatus := vacancyUsecase.validateEmployerAndGetOrgId(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	if userOrgID != vacancy.CompanyID {
		return 0, serverErrors.FORBIDDEN
	}

	vacancyID, addStatus := vacancyUsecase.vacancyRepo.AddVacancy(vacancy)
	if addStatus != nil {
		return 0, addStatus
	}

	return vacancyID, nil
}

func (vacancyUsecase *VacancyUsecase) UpdateVacancy(sessionID string, vacancyID int, vacancy *domain.Vacancy) error {
	validStatus := vacancyUsecase.validateEmployer(sessionID, vacancyID)
	if validStatus != nil {
		return validStatus
	}

	updStatus := vacancyUsecase.vacancyRepo.UpdateVacancy(vacancy)
	if updStatus != nil {
		return updStatus
	}

	return nil
}

func (vacancyUsecase *VacancyUsecase) DeleteVacancy(sessionID string, vacancyID int) error {
	validStatus := vacancyUsecase.validateEmployer(sessionID, vacancyID)
	if validStatus != nil {
		return validStatus
	}

	delStatus := vacancyUsecase.vacancyRepo.DeleteVacancy(vacancyID)
	if delStatus != nil {
		return delStatus
	}

	return nil
}
