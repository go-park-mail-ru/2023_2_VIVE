package psql

import (
	"HnH/pkg/contextUtils"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type ISearchRepository interface {
	GetAllVacanciesIDs(ctx context.Context, limit, offset int64) ([]int64, int64, error)
	FilterCitiesAllVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	FilterSalaryAllVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	FilterExperienceAllVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	FilterEmploymentAllVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	FilterEducationTypeAllVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterCitiesVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterSalaryVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterExperienceVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterEmploymentVacancies(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterEducationTypeVacancies(ctx context.Context) ([]*pb.FilterValue, error)

	GetAllCVsIDs(ctx context.Context, limit, offset int64) ([]int64, int64, error)
	SearchVacanciesIDs(ctx context.Context, query string, limit, offset int64) ([]int64, int64, error)
	SearchCVsIDs(ctx context.Context, query string, limit, offset int64) ([]int64, int64, error)
	// TODO: FilterCitiesAllCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterSalaryAllCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterExperienceAllCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterEmploymentAllCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterEducationTypeAllCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterCitiesCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterSalaryCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterExperienceCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterEmploymentCVs(ctx context.Context) ([]*pb.FilterValue, error)
	// TODO: FilterEducationTypeCVs(ctx context.Context) ([]*pb.FilterValue, error)
}

type psqlSearchRepository struct {
	DB *sql.DB
}

func NewPsqlSearchRepository(db *sql.DB) ISearchRepository {
	return &psqlSearchRepository{
		DB: db,
	}
}

func (repo *psqlSearchRepository) FilterCitiesAllVacancies(ctx context.Context) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting city filters")

	query := `SELECT
			v."location", count(*) AS cnt
		FROM
			hnh_data.vacancy v
		GROUP BY v."location"
		ORDER BY cnt`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filterValues := []*pb.FilterValue{}
	for rows.Next() {
		filterValue := pb.FilterValue{}
		err := rows.Scan(&filterValue.Value, &filterValue.Count)
		if err != nil {
			return nil, err
		}
		filterValues = append(filterValues, &filterValue)
	}

	return filterValues, nil
}

func (repo *psqlSearchRepository) FilterSalaryAllVacancies(ctx context.Context) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting salary filters")

	query := `WITH min_max_salary AS (
		SELECT
			MIN(v.salary_lower_bound) AS min_salary,
			MAX(v.salary_upper_bound) AS max_salary
		FROM
			hnh_data.vacancy v
	), salary_range AS (
		SELECT
			min_salary,
			max_salary,
			(max_salary - min_salary) / 5 AS range_size
		FROM
			min_max_salary
	), salary_buckets AS (
		SELECT
			min_salary + range_size * n AS range_start,
			min_salary + range_size * (n + 1) AS range_end
		FROM
			salary_range,
			generate_series(0, 4) AS n
	)
	SELECT
		sb.range_start,
		sb.range_end,
		COUNT(v.*) AS count
	FROM
		salary_buckets sb
	LEFT JOIN hnh_data.vacancy v ON v.salary_lower_bound >= sb.range_start AND v.salary_lower_bound < sb.range_end
	GROUP BY
		sb.range_start, sb.range_end
	ORDER BY
		sb.range_start`

	rows, err := repo.DB.Query(query)
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"query_error": err,
		}).
			Debug("got query error")
		return nil, err
	}
	defer rows.Close()

	filterValues := []*pb.FilterValue{}
	for rows.Next() {
		var range_start, range_end, count int64
		err := rows.Scan(&range_start, &range_end, &count)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"scan_error": err,
			}).
				Debug("got scan error")
			return nil, err
		}

		filterValues = append(filterValues, &pb.FilterValue{
			Value: fmt.Sprintf("%d", range_start),
			Count: count,
		})
	}

	contextLogger.WithFields(logrus.Fields{
		"filters": filterValues,
	}).
		Debug("got filters")

	return filterValues, nil
}

func (repo *psqlSearchRepository) FilterExperienceAllVacancies(ctx context.Context) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting experience filters")

	query := `SELECT
			v.experience,
			count(*)
		FROM
			hnh_data.vacancy v 
		GROUP BY v.experience`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filterValues := []*pb.FilterValue{}
	for rows.Next() {
		filterValue := pb.FilterValue{}
		err := rows.Scan(&filterValue.Value, &filterValue.Count)
		if err != nil {
			return nil, err
		}
		filterValues = append(filterValues, &filterValue)
	}

	return filterValues, nil
}

func (repo *psqlSearchRepository) FilterEmploymentAllVacancies(ctx context.Context) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting employment filters")

	query := `SELECT
			v.employment,
			count(*)
		FROM
			hnh_data.vacancy v 
		GROUP BY v.employment`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filterValues := []*pb.FilterValue{}
	for rows.Next() {
		filterValue := pb.FilterValue{}
		err := rows.Scan(&filterValue.Value, &filterValue.Count)
		if err != nil {
			return nil, err
		}
		filterValues = append(filterValues, &filterValue)
	}

	return filterValues, nil
}

func (repo *psqlSearchRepository) FilterEducationTypeAllVacancies(ctx context.Context) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting education type filters")

	query := `SELECT
			v.education_type,
			count(*)
		FROM
			hnh_data.vacancy v 
		GROUP BY v.education_type`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filterValues := []*pb.FilterValue{}
	for rows.Next() {
		filterValue := pb.FilterValue{}
		err := rows.Scan(&filterValue.Value, &filterValue.Count)
		if err != nil {
			return nil, err
		}
		filterValues = append(filterValues, &filterValue)
	}

	return filterValues, nil
}

func (repo *psqlSearchRepository) GetAllVacanciesIDs(ctx context.Context, limit, offset int64) ([]int64, int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting all vacancies ids from postgres")
	contextLogger.WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
	}).
		Debug("params")

	queryAllVacs := `WITH all_vacs AS (
			SELECT
				v.id
			FROM
				hnh_data.vacancy v
		),
		count_total AS (
			SELECT
				COUNT(*) AS total
			FROM
				all_vacs
		)
		SELECT
			av.id,
			ct.total
		FROM
			all_vacs av,
			count_total ct
		LIMIT $1
		OFFSET $2`

	rows, err := repo.DB.Query(queryAllVacs, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

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

func (repo *psqlSearchRepository) SearchVacanciesIDs(
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
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	vacIDs := []int64{}
	var count int64

	for rows.Next() {
		var vacID int64
		err := rows.Scan(&vacID, &count)
		if err != nil {
			return nil, 0, err
		}
		contextLogger.WithFields(logrus.Fields{
			"vac_id": vacID,
			"count":  count,
		}).
			Debug("rows.Next()")
		vacIDs = append(vacIDs, vacID)
	}

	contextLogger.WithFields(logrus.Fields{
		"ids":   vacIDs,
		"count": count,
	}).
		Debug("got results")

	return vacIDs, count, nil
}
