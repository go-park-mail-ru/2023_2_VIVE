package psql

import (
	"HnH/internal/repository/psql"
	"HnH/pkg/contextUtils"
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

func (repo *psqlSearchRepository) SearchCVsIDs(
	ctx context.Context,
	searchQuery string,
	limit, offset int64,
) ([]int64, int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("searching vacancies")
	contextLogger.WithFields(logrus.Fields{
		"query":  searchQuery,
		"limit":  limit,
		"offset": offset,
	}).
		Debug("search params")

	query := `WITH filtered_cvs AS (
			SELECT
				cv.id,
				cv.fts
			FROM
				hnh_data.cv cv
			WHERE
				plainto_tsquery($1) @@ cv.fts
		),
		count_total AS (
			SELECT
				COUNT(*) AS total
			FROM
				filtered_cvs
		)
		SELECT
			fcv.id,
			ct.total
		FROM
			filtered_cvs fcv,
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
	defer rows.Close()

	cvsIDs := []int64{}
	var count int64

	for rows.Next() {
		var cvID int64
		err := rows.Scan(&cvID, &count)
		if err != nil {
			return nil, 0, err
		}
		cvsIDs = append(cvsIDs, cvID)
	}

	contextLogger.WithFields(logrus.Fields{
		"ids":   cvsIDs,
		"count": count,
	}).
		Debug("got results")

	return cvsIDs, count, nil
}
