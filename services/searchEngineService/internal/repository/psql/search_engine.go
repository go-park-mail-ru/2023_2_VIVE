package psql

import (
	"HnH/internal/repository/psql"
	"HnH/pkg/contextUtils"
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type ISearchRepository interface {
	SearchVacanciesIDs(ctx context.Context, query string, pageNumber, resultsPerPage int64) ([]int64, int64, error)
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
	searchQuery string,
	pageNumber,
	resultsPerPage int64,
) ([]int64, int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("searching vacancies")
	contextLogger.WithFields(logrus.Fields{
		"query":            searchQuery,
		"page_number":      pageNumber,
		"results_per_page": resultsPerPage,
	}).
		Debug("search params")

	limit := resultsPerPage
	offset := (pageNumber - 1) * resultsPerPage

	query := `WITH filtered_vacancies AS (
			SELECT
				v.id,
				v.fts
			FROM
				hnh_data.vacancy v
			WHERE
			plainto_tsquery($1) @@ v.fts
		),
		count_total AS (
			SELECT
				COUNT(*) AS total
			FROM
				filtered_vacancies
		)
		SELECT
			fv.id,
			ct.total
		FROM
			filtered_vacancies fv,
			count_total ct
		LIMIT $2
		OFFSET $3`

	rows, err := repo.DB.Query(query, searchQuery, limit, offset)
	if err == sql.ErrNoRows {
		return nil, 0, psql.ErrEntityNotFound
	}
	if err != nil {
		return nil, 0, err
	}

	vacIDs := []int64{}
	var count int64

	for rows.Next() {
		var vacID int64
		err := rows.Scan(&vacID, &count)
		if err != nil {
			return nil, 0, err
		}
		vacIDs = append(vacIDs, vacID)
	}

	contextLogger.WithFields(logrus.Fields{
		"ids":   vacIDs,
		"count": count,
	}).
		Debug("got results")

	return vacIDs, count, nil
}
