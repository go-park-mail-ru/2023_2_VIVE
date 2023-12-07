package searchOptions

import (
	"HnH/pkg/castUtils"
	"HnH/pkg/queryUtils"
	"HnH/services/searchEngineService/searchEnginePB"
	"fmt"
	"strconv"
	"strings"
)

type OptionName string

const (
	SearchQuery         OptionName = "q"
	PageNum             OptionName = "page_num"
	ResultsPerPage      OptionName = "results_per_page"
	CityFilter          OptionName = "city"
	SalaryFilter        OptionName = "salary"
	EmploymentFilter    OptionName = "employment"
	ExperienceFilter    OptionName = "experience"
	EducationTypeFilter OptionName = "education_type"
	GenderFilter        OptionName = "gender"
)

type FilterType string

const (
	CheckBox       FilterType = "checkbox"
	Radio          FilterType = "radio"
	CheckBoxSearch FilterType = "checkbox_search"
	DoubleRange    FilterType = "double_range"
)

func PopSearchQuery(options *searchEnginePB.SearchOptions) (string, error) {
	if len(options.GetOptions()) > 0 {
		searchQuery := options.GetOptions()[string(SearchQuery)].GetValues()[0]
		// delete(options.Options, string(SearchQuery))
		return searchQuery, nil
	}
	return "", ErrNoOption
}

func PopPageNum(options *searchEnginePB.SearchOptions) (int64, error) {
	if len(options.GetOptions()) > 0 {
		pageNumStr := options.GetOptions()[string(PageNum)].GetValues()[0]
		pageNum, convErr := strconv.ParseInt(pageNumStr, 10, 64)
		if convErr != nil {
			// responseTemplates.SendErrorMessage(w, ErrWrongQueryParam, http.StatusBadRequest)
			return 0, convErr
		}
		delete(options.Options, string(PageNum))
		return pageNum, nil
	}
	return 0, ErrNoOption
}

func PopResultsPerPage(options *searchEnginePB.SearchOptions) (int64, error) {
	if len(options.GetOptions()) > 0 {
		resultsPerPage := options.GetOptions()[string(ResultsPerPage)].GetValues()[0]
		pageNum, convErr := strconv.ParseInt(resultsPerPage, 10, 64)
		if convErr != nil {
			// responseTemplates.SendErrorMessage(w, ErrWrongQueryParam, http.StatusBadRequest)
			return 0, convErr
		}
		delete(options.Options, string(ResultsPerPage))
		return pageNum, nil
	}
	return 0, ErrNoOption
}

func GetWhereTerms(options *searchEnginePB.SearchOptions, placeHolderIndex int) ([]string, []interface{}) {
	queryElemets := []string{}
	var args []interface{}

	// placeHolderIndex := 1
	for optionName, optionValues := range options.Options {
		switch optionName {
		case string(SearchQuery):
			if strings.TrimSpace(optionValues.Values[0]) != "" {
				queryElemets = append(queryElemets, fmt.Sprintf("plainto_tsquery($%d) @@ tbl.fts", placeHolderIndex))
				placeHolderIndex++
				args = append(args, optionValues.Values[0])
			}

		case string(CityFilter):
			if len(optionValues.Values) > 0 {
				values := strings.Split(optionValues.Values[0], ",")
				placeholders := queryUtils.QueryPlaceHolders(placeHolderIndex, len(values))
				placeHolderIndex += len(values)
				queryElemets = append(queryElemets, fmt.Sprintf(`tbl."location" IN (%s)`, placeholders))
				args = append(args, castUtils.StringToAnySlice(values)...)
			}

		case string(SalaryFilter):
			if len(optionValues.Values) > 0 {
				rangeBounds := strings.Split(optionValues.Values[0], ",")
				if len(rangeBounds) != 2 {
					continue
				}
				rangeBounds = append(rangeBounds, rangeBounds...)
				placeholders := queryUtils.GetPlaceholdersSliceFrom(placeHolderIndex, 4)
				// placeholders = append(placeholders, placeholders...)
				placeHolderIndex += 4
				queryElemets = append(queryElemets, fmt.Sprintf(
					`(tbl.salary_lower_bound BETWEEN %s AND %s OR tbl.salary_upper_bound BETWEEN %s AND %s)`,
					castUtils.StringToAnySlice(placeholders)...,
				))
				args = append(args, castUtils.StringToAnySlice(rangeBounds)...)
			}

		case string(EmploymentFilter):
			if len(optionValues.Values) > 0 {
				values := strings.Split(optionValues.Values[0], ",")
				placeholders := queryUtils.QueryPlaceHolders(placeHolderIndex, len(values))
				placeHolderIndex += len(values)
				queryElemets = append(queryElemets, fmt.Sprintf(`tbl.employment IN (%s)`, placeholders))
				args = append(args, castUtils.StringToAnySlice(values)...)
			}

		case string(ExperienceFilter):
			if len(optionValues.Values) > 0 {
				values := strings.Split(optionValues.Values[0], ",")
				placeholders := queryUtils.QueryPlaceHolders(placeHolderIndex, len(values))
				placeHolderIndex += len(values)
				queryElemets = append(queryElemets, fmt.Sprintf(`tbl.experience IN (%s)`, placeholders))
				args = append(args, castUtils.StringToAnySlice(values)...)
			}

		case string(EducationTypeFilter):
			if len(optionValues.Values) > 0 {
				values := strings.Split(optionValues.Values[0], ",")
				placeholders := queryUtils.QueryPlaceHolders(placeHolderIndex, len(values))
				placeHolderIndex += len(values)
				queryElemets = append(queryElemets, fmt.Sprintf(`tbl.education_type IN (%s)`, placeholders))
				args = append(args, castUtils.StringToAnySlice(values)...)
			}
		}

	}

	return queryElemets, args
}
