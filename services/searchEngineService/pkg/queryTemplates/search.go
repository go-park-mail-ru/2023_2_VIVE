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
		allowedOptions: []Option{
			SearchOption,
			CityOption,
			SalaryOption,
			EmploymentOption,
			ExperienceOption,
			VacEducationTypeOption,
		},
	}

	CVsSearchQueryTemplate = &SearchQueryTemplates{
		table_name: "hnh_data.cv",
		allowedOptions: []Option{
			SearchOption,
			CityOption,
			CvEducationTypeOption,
			// ExperienceOption, TODO: denormilize cv table
			GenderOption,
		},
	}
)

type SearchQueryTemplates struct {
	table_name      string
	allowedOptions  []Option
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

	for _, option := range sqt.allowedOptions {
		optionValues, exists := options.Options[string(option.Name)]
		if !exists {
			continue
		}

		var whereTerm string
		whereTerm, args = option.Handle(optionValues, placeHolderStartIndex, args)
		if whereTerm != "" {
			whereTerms = append(whereTerms, whereTerm)
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
