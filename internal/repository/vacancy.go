package repository

import (
	"HnH/internal/domain"
	"HnH/internal/repository/mock"
)

type IVacancyRepository interface {
	GetVacancies() ([]domain.Vacancy, error)
}

type psqlVacancyRepository struct {
	vacancyStorage *mock.Vacancies
}

func NewPsqlVacancyRepository(vacancies *mock.Vacancies) IVacancyRepository {
	return &psqlVacancyRepository{
		vacancyStorage: vacancies,
	}
}

func (p *psqlVacancyRepository) GetVacancies() ([]domain.Vacancy, error) {
	p.vacancyStorage.Mu.RLock()

	defer p.vacancyStorage.Mu.RUnlock()

	listToReturn := p.vacancyStorage.VacancyList

	return listToReturn, nil
}
