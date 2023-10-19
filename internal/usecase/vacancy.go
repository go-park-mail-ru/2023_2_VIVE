package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"
)

type IVacancyUsecase interface {
	GetVacancies() ([]domain.Vacancy, error)
}

type VacancyUsecase struct {
	vacancyRepo repository.IVacancyRepository
}

func NewVacancyUsecase(vacancyRepository repository.IVacancyRepository) IVacancyUsecase {
	return &VacancyUsecase{
		vacancyRepo: vacancyRepository,
	}
}

func (vacancyUsecase *VacancyUsecase) GetVacancies() ([]domain.Vacancy, error) {
	vacancies, getErr := vacancyUsecase.vacancyRepo.GetVacancies()
	if getErr != nil {
		return nil, getErr
	}

	return vacancies, nil
}
