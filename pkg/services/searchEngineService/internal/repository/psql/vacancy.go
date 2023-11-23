package psql

import (
	"HnH/internal/domain"
	"database/sql"
)

type IVacancyRepository interface {
	SearchVacancies(words []string, pageNumber, resultsPerPage int32) ([]domain.ApiVacancy, error)
}

type psqlVacancyRepository struct {
	DB *sql.DB
}

func NewPsqlVacancyRepository(db *sql.DB) IVacancyRepository {
	return &psqlVacancyRepository{
		DB: db,
	}
}

func (repo *psqlVacancyRepository) SearchVacancies(words []string, pageNumber, resultsPerPage int32) ([]domain.ApiVacancy, error) {
	return []domain.ApiVacancy{}, nil
}
