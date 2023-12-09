package queryTemplates

import (
	"HnH/pkg/castUtils"
	"HnH/pkg/queryUtils"
	"HnH/services/searchEngineService/pkg/searchOptions"
	"HnH/services/searchEngineService/searchEnginePB"
	"fmt"
	"strings"
)

type OptionHandler func(
	option *Option,
	optionValues *searchEnginePB.SearchOptionValues,
	placeHolderStartIndex *int,
	args []interface{},
) (string, []interface{})

type Option struct {
	Name          searchOptions.OptionName
	Type          searchOptions.OptionType
	DBColumnNames []string
	Handler       OptionHandler
	// AdditionalDBColumnNames []string
}

func (op *Option) Handle(
	optionValues *searchEnginePB.SearchOptionValues,
	placeHolderStartIndex *int,
	args []interface{},
) (string, []interface{}) {
	if op.Handler != nil {
		return op.Handler(op, optionValues, placeHolderStartIndex, args)
	}

	return "", args
}

var (
	SearchOption = Option{
		Name:    searchOptions.SearchQuery,
		Type:    searchOptions.Search,
		Handler: handleSearchOption,
	}
	CityOption = Option{
		Name:          searchOptions.City,
		Type:          searchOptions.CheckBox,
		DBColumnNames: []string{`"location"`},
		Handler:       handleCheckBoxOption,
	}
	SalaryOption = Option{
		Name:          searchOptions.Salary,
		Type:          searchOptions.DoubleRange,
		DBColumnNames: []string{"salary_lower_bound", "salary_upper_bound"},
		Handler:       handleDoubleRangeOption,
		// AdditionalDBColumnNames: []string{"salary_upper_bound"},
	}
	EmploymentOption = Option{
		Name:          searchOptions.Employment,
		Type:          searchOptions.CheckBox,
		DBColumnNames: []string{"employment"},
		Handler:       handleCheckBoxOption,
	}
	ExperienceOption = Option{
		Name:          searchOptions.Experience,
		Type:          searchOptions.CheckBox,
		DBColumnNames: []string{"experience"},
		Handler:       handleCheckBoxOption,
	}
	VacEducationTypeOption = Option{
		Name:          searchOptions.EducationType,
		Type:          searchOptions.CheckBox,
		DBColumnNames: []string{"education_type"},
		Handler:       handleCheckBoxOption,
	}
	CvEducationTypeOption = Option{
		Name:          searchOptions.EducationType,
		Type:          searchOptions.CheckBox,
		DBColumnNames: []string{"education_level"},
		Handler:       handleCheckBoxOption,
	}
	GenderOption = Option{
		Name:          searchOptions.Gender,
		Type:          searchOptions.Radio,
		DBColumnNames: []string{"gender"},
		Handler:       handleRadioOption,
	}
)

func handleSearchOption(
	option *Option,
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
	option *Option,
	optionValues *searchEnginePB.SearchOptionValues,
	placeHolderStartIndex *int,
	args []interface{},
) (string, []interface{}) {
	if len(optionValues.GetValues()) > 0 && len(option.DBColumnNames) > 0 {
		values := strings.Split(optionValues.Values[0], ",")
		if len(values) == 0 {
			return "", args
		}

		placeholders := queryUtils.QueryPlaceHolders(*placeHolderStartIndex, len(values))
		*placeHolderStartIndex += len(values)
		whereTerm := fmt.Sprintf(`tbl.%s IN (%s)`, option.DBColumnNames[0], placeholders)
		args = append(args, castUtils.StringToAnySlice(values)...)

		return whereTerm, args
	}
	return "", args
}

func handleRadioOption(
	option *Option,
	optionValues *searchEnginePB.SearchOptionValues,
	placeHolderStartIndex *int,
	args []interface{},
) (string, []interface{}) {
	if len(optionValues.GetValues()) > 0 && len(option.DBColumnNames) > 0 {
		value := optionValues.GetValues()[0]
		if strings.TrimSpace(value) == "" {
			return "", args
		}

		placeholders := queryUtils.QueryPlaceHolders(*placeHolderStartIndex, 1)
		*placeHolderStartIndex++
		whereTerm := fmt.Sprintf(`tbl.%s = %s`, option.DBColumnNames[0], placeholders)
		args = append(args, value)

		return whereTerm, args
	}
	return "", args
}

func handleDoubleRangeOption(
	option *Option,
	optionValues *searchEnginePB.SearchOptionValues,
	placeHolderStartIndex *int,
	args []interface{},
) (string, []interface{}) {
	if len(optionValues.Values) > 0 && len(option.DBColumnNames) == 2 {
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
			option.DBColumnNames[0],
			placeholders[0],
			placeholders[1],
			option.DBColumnNames[1],
			placeholders[2],
			placeholders[3],
		)
		args = append(args, castUtils.StringToAnySlice(rangeBounds)...)

		return whereTerm, args
	}
	return "", args
}
