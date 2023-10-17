package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"
)

func GetVacancies() ([]domain.Vacancy, error) {
	vacancies, getErr := repository.GetVacancies()
	if getErr != nil {
		return nil, getErr
	}

	return vacancies, nil
}
