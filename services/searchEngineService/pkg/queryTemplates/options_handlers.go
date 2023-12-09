package queryTemplates

import (
	"HnH/pkg/castUtils"
	"HnH/pkg/queryUtils"
	"HnH/services/searchEngineService/pkg/searchOptions"
	"HnH/services/searchEngineService/searchEnginePB"
	"fmt"
	"strings"
)

var (
	SearchOption = searchOptions.Option{
		Name: searchOptions.SearchQuery,
		Type: searchOptions.Search,
		Handler: handleSearchOption,
	}
	CityOption = searchOptions.Option{
		Name:         searchOptions.City,
		Type:         searchOptions.CheckBox,
		DBColumnName: `"location"`,
		Handler:      handleCheckBoxOption,
	}
	SalaryOption = searchOptions.Option{
		Name:                    searchOptions.Salary,
		Type:                    searchOptions.DoubleRange,
		DBColumnName:            "salary_lower_bound",
		AdditionalDBColumnNames: []string{"salary_upper_bound"},
		Handler:                 handleDoubleRangeOption,
	}
	EmploymentOption = searchOptions.Option{
		Name:         searchOptions.Employment,
		Type:         searchOptions.CheckBox,
		DBColumnName: "employment",
		Handler:      handleCheckBoxOption,
	}
	ExperienceOption = searchOptions.Option{
		Name:         searchOptions.Experience,
		Type:         searchOptions.CheckBox,
		DBColumnName: "experience",
		Handler:      handleCheckBoxOption,
	}
	EducationTypeOption = searchOptions.Option{
		Name:         searchOptions.EducationType,
		Type:         searchOptions.CheckBox,
		DBColumnName: "education_type",
		Handler:      handleCheckBoxOption,
	}
)

func handleSearchOption(
	option *searchOptions.Option,
	optionValues *searchEnginePB.SearchOptionValues,
	placeHolderStartIndex *int,
	args []interface{},
) (string, []interface{}) {
	if len(optionValues.GetValues()) > 0 && strings.TrimSpace(optionValues.GetValues()[0]) != "" {
		whereTerm := fmt.Sprintf("plainto_tsquery($%d) @@ tbl.fts", *placeHolderStartIndex)
		*placeHolderStartIndex++
		args = append(args, optionValues.GetValues()[0])

		return whereTerm, args
	}
	return "", args
}

func handleCheckBoxOption(
	option *searchOptions.Option,
	optionValues *searchEnginePB.SearchOptionValues,
	placeHolderStartIndex *int,
	args []interface{},
) (string, []interface{}) {
	if len(optionValues.GetValues()) > 0 {
		values := strings.Split(optionValues.Values[0], ",")
		if len(values) == 0 {
			return "", args
		}

		placeholders := queryUtils.QueryPlaceHolders(*placeHolderStartIndex, len(values))
		*placeHolderStartIndex += len(values)
		whereTerm := fmt.Sprintf(`tbl.%s IN (%s)`, option.DBColumnName, placeholders)
		args = append(args, castUtils.StringToAnySlice(values)...)

		return whereTerm, args
	}
	return "", args
}

func handleDoubleRangeOption(
	option *searchOptions.Option,
	optionValues *searchEnginePB.SearchOptionValues,
	placeHolderStartIndex *int,
	args []interface{},
) (string, []interface{}) {
	if len(option.AdditionalDBColumnNames) > 0 && len(optionValues.Values) > 0 {
		rangeBounds := strings.Split(optionValues.Values[0], ",")
		if len(rangeBounds) != 2 {
			return "", args
		}

		rangeBounds = append(rangeBounds, rangeBounds...)
		placeholders := queryUtils.GetPlaceholdersSliceFrom(*placeHolderStartIndex, 4)
		// placeholders = append(placeholders, placeholders...)
		*placeHolderStartIndex += 4
		whereTerm := fmt.Sprintf(
			`(tbl.%s BETWEEN %s AND %s OR tbl.%s BETWEEN %s AND %s)`,
			option.DBColumnName,
			placeholders[0],
			placeholders[1],
			option.AdditionalDBColumnNames[0],
			placeholders[2],
			placeholders[3],
		)
		args = append(args, castUtils.StringToAnySlice(rangeBounds)...)

		return whereTerm, args
	}
	return "", args
}
