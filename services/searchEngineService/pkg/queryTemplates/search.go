package queryTemplates

import (
	"HnH/pkg/castUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/queryUtils"
	"HnH/services/searchEngineService/pkg/searchOptions"
	"HnH/services/searchEngineService/searchEnginePB"
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	VacanciesSearchQueryTemplate = &SearchQueryTemplates{
		table_name: "hnh_data.vacancy",
		// limitClause:  "LIMIT $1",
		// offsetClause: "OFFSET $2",
	}

	CVsSearchQueryTemplate = &SearchQueryTemplates{
		table_name: "hnh_data.cv",
		// limitClause:  "LIMIT $1",
		// offsetClause: "OFFSET $2",
	}
)

type SearchQueryTemplates struct {
	table_name      string
	whereClause     string
	paginatorClause string
}

func (sqt *SearchQueryTemplates) setPaginatorClause(
	options *searchEnginePB.SearchOptions,
	placeHolderStartIndex *int,
	args []interface{},
) ([]interface{}, error) {
	pageNum, err := searchOptions.GetPageNum(options)
	if err != nil {
		return args, err
	}
	resultsPerPage, err := searchOptions.GetResultsPerPage(options)
	if err != nil {
		return args, err
	}

	limit := resultsPerPage
	offset := (pageNum - 1) * resultsPerPage

	args = append(args, limit, offset)
	placeholders := queryUtils.GetPlaceholdersSliceFrom(*placeHolderStartIndex, 2)
	sqt.paginatorClause = fmt.Sprintf("LIMIT %s OFFSET %s", castUtils.StringToAnySlice(placeholders)...)
	*placeHolderStartIndex += 2

	return args, nil
}

func (sqt *SearchQueryTemplates) setWhereCaluse(
	options *searchEnginePB.SearchOptions,
	placeHolderStartIndex *int,
	args []interface{},
) ([]interface{}, error) {
	whereTerms := []string{}

	// TODO: make universal methods for each filter type
	for optionName, optionValues := range options.Options {
		switch optionName {
		case string(searchOptions.SearchQuery):
			if strings.TrimSpace(optionValues.Values[0]) != "" {
				whereTerms = append(whereTerms, fmt.Sprintf("plainto_tsquery($%d) @@ tbl.fts", *placeHolderStartIndex))
				*placeHolderStartIndex++
				args = append(args, optionValues.Values[0])
			}

		case string(searchOptions.CityFilter):
			if len(optionValues.Values) > 0 {
				values := strings.Split(optionValues.Values[0], ",")
				placeholders := queryUtils.QueryPlaceHolders(*placeHolderStartIndex, len(values))
				*placeHolderStartIndex += len(values)
				whereTerms = append(whereTerms, fmt.Sprintf(`tbl."location" IN (%s)`, placeholders))
				args = append(args, castUtils.StringToAnySlice(values)...)
			}

		case string(searchOptions.SalaryFilter):
			if len(optionValues.Values) > 0 {
				rangeBounds := strings.Split(optionValues.Values[0], ",")
				if len(rangeBounds) != 2 {
					continue
				}
				rangeBounds = append(rangeBounds, rangeBounds...)
				placeholders := queryUtils.GetPlaceholdersSliceFrom(*placeHolderStartIndex, 4)
				// placeholders = append(placeholders, placeholders...)
				*placeHolderStartIndex += 4
				whereTerms = append(whereTerms, fmt.Sprintf(
					`(tbl.salary_lower_bound BETWEEN %s AND %s OR tbl.salary_upper_bound BETWEEN %s AND %s)`,
					castUtils.StringToAnySlice(placeholders)...,
				))
				args = append(args, castUtils.StringToAnySlice(rangeBounds)...)
			}

		case string(searchOptions.EmploymentFilter):
			if len(optionValues.Values) > 0 {
				values := strings.Split(optionValues.Values[0], ",")
				placeholders := queryUtils.QueryPlaceHolders(*placeHolderStartIndex, len(values))
				*placeHolderStartIndex += len(values)
				whereTerms = append(whereTerms, fmt.Sprintf(`tbl.employment IN (%s)`, placeholders))
				args = append(args, castUtils.StringToAnySlice(values)...)
			}

		case string(searchOptions.ExperienceFilter):
			if len(optionValues.Values) > 0 {
				values := strings.Split(optionValues.Values[0], ",")
				placeholders := queryUtils.QueryPlaceHolders(*placeHolderStartIndex, len(values))
				*placeHolderStartIndex += len(values)
				whereTerms = append(whereTerms, fmt.Sprintf(`tbl.experience IN (%s)`, placeholders))
				args = append(args, castUtils.StringToAnySlice(values)...)
			}

		case string(searchOptions.EducationTypeFilter):
			if len(optionValues.Values) > 0 {
				values := strings.Split(optionValues.Values[0], ",")
				placeholders := queryUtils.QueryPlaceHolders(*placeHolderStartIndex, len(values))
				*placeHolderStartIndex += len(values)
				whereTerms = append(whereTerms, fmt.Sprintf(`tbl.education_type IN (%s)`, placeholders))
				args = append(args, castUtils.StringToAnySlice(values)...)
			}
		}

	}

	if len(whereTerms) > 0 {
		sqt.whereClause = fmt.Sprintf("WHERE %s", strings.Join(whereTerms, " AND "))
	}

	return args, nil
}

func (sqt SearchQueryTemplates) BuildTemplate(ctx context.Context, options *searchEnginePB.SearchOptions) (string, []interface{}) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	var args []interface{}
	placeHolderIndex := 1
	args, paginatorErr := sqt.setPaginatorClause(options, &placeHolderIndex, args)
	if paginatorErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"error_msg": paginatorErr,
			"options":   options.Options,
		}).
			Warn("warning while parsing for paginator options")
	}

	args, whereErr := sqt.setWhereCaluse(options, &placeHolderIndex, args)
	if whereErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"error_msg": whereErr,
			"options":   options,
		}).
			Warn("warning while parsing for where options")
	}

	return fmt.Sprintf(
		`WITH filtered_items AS (
		SELECT
			tbl.id,
			tbl.fts
		FROM
			%s tbl
		%s
	),
	count_total AS (
		SELECT
			COUNT(*) AS total
		FROM
			filtered_items
	)
	SELECT
		fi.id,
		ct.total
	FROM
		filtered_items fi,
		count_total ct
	%s`,
		sqt.table_name,
		sqt.whereClause,
		sqt.paginatorClause,
	), args
}
