package usecase

import (
	"HnH/internal/domain"
)

type VacancyUsecase struct {
	vacancyRepo VacancyRepository
}

func NewVacancyUsecase(vacancyRepository VacancyRepository) *VacancyUsecase {
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
