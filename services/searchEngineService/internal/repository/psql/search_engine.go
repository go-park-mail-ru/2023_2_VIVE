package psql

import (
	"context"
	"database/sql"
)

type ISearchRepository interface {
	SearchVacanciesIDs(ctx context.Context, query string, pageNumber, resultsPerPage int64) ([]int64, error)
}

type psqlSearchRepository struct {
	DB *sql.DB
}

func NewPsqlSearchRepository(db *sql.DB) ISearchRepository {
	return &psqlSearchRepository{
		DB: db,
	}
}

func (repo *psqlSearchRepository) SearchVacanciesIDs(
	ctx context.Context,
	query string,
	pageNumber,
	resultsPerPage int64,
) ([]int64, error) {
	return []int64{}, nil
}
