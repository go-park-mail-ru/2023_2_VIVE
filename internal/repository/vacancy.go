package repository

import (
	"HnH/internal/domain"
	"HnH/internal/repository/mock"
)

func GetVacancies() ([]domain.Vacancy, error) {
	mock.VacancyDB.Mu.RLock()

	defer mock.VacancyDB.Mu.RUnlock()

	listToReturn := mock.VacancyDB.VacancyList

	return listToReturn, nil
}
