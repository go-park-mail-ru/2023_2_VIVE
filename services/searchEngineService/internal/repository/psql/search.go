package psql

import (
	"HnH/internal/repository/psql"
	"HnH/pkg/contextUtils"
	"HnH/services/searchEngineService/pkg/queryTemplates"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type ISearchRepository interface {
	SearchVacanciesIDs(ctx context.Context /* , searchQuery string, limit, offset int64 */, options *pb.SearchOptions) ([]int64, int64, error)
	FilterCitiesVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)
	FilterSalaryVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)
	FilterExperienceVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)
	FilterEmploymentVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)
	FilterEducationTypeVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)

	// GetAllCVsIDs(ctx context.Context, limit, offset int64) ([]int64, int64, error)
	SearchCVsIDs(ctx context.Context /* query string, limit, offset int64,  */, options *pb.SearchOptions) ([]int64, int64, error)
	FilterCitiesCVs(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)
	FilterEducationTypeCVs(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)
	FilterGenderCVs(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)
	// TODO: FilterExperienceCVs(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error)
}

type psqlSearchRepository struct {
	DB *sql.DB
}

func NewPsqlSearchRepository(db *sql.DB) ISearchRepository {
	return &psqlSearchRepository{
		DB: db,
	}
}

// commonFilterItems executes a query built from a template and returns filter values.
func (repo *psqlSearchRepository) commonFilterItems(ctx context.Context, qt *queryTemplates.CommonFilterQueryTemplate, args ...interface{}) ([]*pb.FilterValue, error) {
	return repo.executeCommonFilterQuery(ctx, qt.BuildQuery(len(args) == 1), args...)
}

// executeCommonFilterQuery is a helper function to execute a SQL query and scan results into filter values.
func (repo *psqlSearchRepository) executeCommonFilterQuery(ctx context.Context, query string, args ...interface{}) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"db_query": query,
		"args":     args,
	}).
		Debug("executing db query with args")

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var filterValues []*pb.FilterValue
	for rows.Next() {
		var filterValue pb.FilterValue
		if err := rows.Scan(&filterValue.Value, &filterValue.Count); err != nil {
			return nil, err
		}
		filterValues = append(filterValues, &filterValue)
	}

	return filterValues, nil
}

func (repo *psqlSearchRepository) salaryFilterItems(ctx context.Context, qt *queryTemplates.CommonFilterQueryTemplate, args ...interface{}) ([]*pb.FilterValue, error) {
	return repo.executeSalaryFilterQuery(ctx, qt.BuildQuery(len(args) == 1), args...)
}

// func compareStartEndRanges(range_start, range_end int64) (int64, int64) {
// 	if range_end < range_start {
// 		return range_start, range_start
// 	} else if
// }

func (repo *psqlSearchRepository) executeSalaryFilterQuery(ctx context.Context, query string, args ...interface{}) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"db_query": query,
		"args":     args,
	}).
		Debug("executing db query with args")

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var filterValues []*pb.FilterValue
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

		if range_end < range_start {
			range_end = range_start
		}

		filterValues = append(filterValues, &pb.FilterValue{
			Value: fmt.Sprintf("%d:%d", range_start, range_end),
			Count: count,
		})

	}

	return filterValues, nil
}

func (repo *psqlSearchRepository) FilterCitiesVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting city filters")
	if strings.TrimSpace(searchQuery) != "" {
		return repo.commonFilterItems(ctx, queryTemplates.VacCitiesQueryTemplate, searchQuery)
	}
	return repo.commonFilterItems(ctx, queryTemplates.VacCitiesQueryTemplate)
}

func (repo *psqlSearchRepository) FilterSalaryVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting salary filters")
	if strings.TrimSpace(searchQuery) != "" {
		return repo.salaryFilterItems(ctx, queryTemplates.VacSalaryQueryTemplate, searchQuery)
	}
	return repo.salaryFilterItems(ctx, queryTemplates.VacSalaryQueryTemplate)
}

func (repo *psqlSearchRepository) FilterExperienceVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting experience filters")
	if strings.TrimSpace(searchQuery) != "" {
		return repo.commonFilterItems(ctx, queryTemplates.VacExperienceQueryTemplate, searchQuery)
	}
	return repo.commonFilterItems(ctx, queryTemplates.VacExperienceQueryTemplate)
}

func (repo *psqlSearchRepository) FilterEmploymentVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting employment filters")
	if strings.TrimSpace(searchQuery) != "" {
		return repo.commonFilterItems(ctx, queryTemplates.VacEmploymentQueryTemplate, searchQuery)
	}
	return repo.commonFilterItems(ctx, queryTemplates.VacEmploymentQueryTemplate)
}

func (repo *psqlSearchRepository) FilterEducationTypeVacancies(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting EducationType filters")
	if strings.TrimSpace(searchQuery) != "" {
		return repo.commonFilterItems(ctx, queryTemplates.VacEducationTypeQueryTemplate, searchQuery)
	}
	return repo.commonFilterItems(ctx, queryTemplates.VacEducationTypeQueryTemplate)
}

func (repo *psqlSearchRepository) executeSearchQuery(ctx context.Context, query string, args ...interface{}) ([]int64, int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"db_query": query,
		"args":     args,
	}).
		Debug("executing query with args")

	rows, err := repo.DB.Query(query, args...)
	if err == sql.ErrNoRows {
		return nil, 0, psql.ErrEntityNotFound
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	ids := []int64{}
	var count int64

	for rows.Next() {
		var id int64
		err := rows.Scan(&id, &count)
		if err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	contextLogger.WithFields(logrus.Fields{
		"ids":   ids,
		"count": count,
	}).
		Debug("got results")

	return ids, count, nil
}

func (repo *psqlSearchRepository) searchItems(
	ctx context.Context,
	qt *queryTemplates.SearchQueryTemplates,
	options *pb.SearchOptions,
) ([]int64, int64, error) {
	query, args := qt.BuildTemplate(ctx, options)
	return repo.executeSearchQuery(ctx, query, args...)
}

func (repo *psqlSearchRepository) SearchVacanciesIDs(
	ctx context.Context,
	options *pb.SearchOptions,
) ([]int64, int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("searching vacancies")
	contextLogger.WithFields(logrus.Fields{
		"options": options.Options,
	}).
		Debug("search options")

	return repo.searchItems(ctx, queryTemplates.VacanciesSearchQueryTemplate, options)
}

func (repo *psqlSearchRepository) SearchCVsIDs(
	ctx context.Context,
	options *pb.SearchOptions,
) ([]int64, int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("searching cvs")
	contextLogger.WithFields(logrus.Fields{
		"options": options.Options,
	}).
		Debug("search params")

	return repo.searchItems(ctx, queryTemplates.CVsSearchQueryTemplate, options)
}

func (repo *psqlSearchRepository) FilterCitiesCVs(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting city filters")

	query := queryTemplates.CvCitiesQueryTemplate.BuildQuery(strings.TrimSpace(searchQuery) != "")

	if strings.TrimSpace(searchQuery) != "" {
		return repo.executeCommonFilterQuery(ctx, query, searchQuery)
	}
	return repo.executeCommonFilterQuery(ctx, query)
}

func (repo *psqlSearchRepository) FilterEducationTypeCVs(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting city filters")

	query := queryTemplates.CvEducationTypeQueryTemplate.BuildQuery(strings.TrimSpace(searchQuery) != "")

	if strings.TrimSpace(searchQuery) != "" {
		return repo.executeCommonFilterQuery(ctx, query, searchQuery)
	}
	return repo.executeCommonFilterQuery(ctx, query)
}

func (repo *psqlSearchRepository) FilterGenderCVs(ctx context.Context, searchQuery string) ([]*pb.FilterValue, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting city filters")

	query := queryTemplates.CvGenderQueryTemplate.BuildQuery(strings.TrimSpace(searchQuery) != "")

	if strings.TrimSpace(searchQuery) != "" {
		return repo.executeCommonFilterQuery(ctx, query, searchQuery)
	}
	return repo.executeCommonFilterQuery(ctx, query)
}
