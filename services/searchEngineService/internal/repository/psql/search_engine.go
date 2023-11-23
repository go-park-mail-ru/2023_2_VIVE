package psql

import (
	"HnH/internal/domain"
	"context"
	"database/sql"
)

type ISearchRepository interface {
	SearchVacancies(ctx context.Context, words []string, pageNumber, resultsPerPage int64) ([]domain.DbVacancy, error)
}

type psqlSearchRepository struct {
	DB *sql.DB
}

func NewPsqlSearchRepository(db *sql.DB) ISearchRepository {
	return &psqlSearchRepository{
		DB: db,
	}
}

func (repo *psqlSearchRepository) SearchVacancies(ctx context.Context, words []string, pageNumber, resultsPerPage int64) ([]domain.DbVacancy, error) {
	return []domain.DbVacancy{}, nil
}
