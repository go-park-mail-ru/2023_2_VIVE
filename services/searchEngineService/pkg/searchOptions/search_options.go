package searchOptions

import (
	"HnH/services/searchEngineService/searchEnginePB"
	"strconv"
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

func GetSearchQuery(options *searchEnginePB.SearchOptions) (string, error) {
	optionValues, exists := options.Options[string(SearchQuery)]
	if !exists {
		return "", ErrNoOption
	}
	if len(optionValues.GetValues()) > 0 {
		searchQuery := optionValues.GetValues()[0]
		// delete(options.Options, string(SearchQuery))
		return searchQuery, nil
	}
	return "", ErrWrongValueFormat
}

func GetPageNum(options *searchEnginePB.SearchOptions) (int64, error) {
	optionValues, exists := options.Options[string(PageNum)]
	if !exists {
		return 0, ErrNoOption
	}
	if len(optionValues.GetValues()) > 0 {
		pageNumStr := optionValues.GetValues()[0]
		pageNum, convErr := strconv.ParseInt(pageNumStr, 10, 64)
		if convErr != nil {
			// responseTemplates.SendErrorMessage(w, ErrWrongQueryParam, http.StatusBadRequest)
			return 0, convErr
		}
		// delete(options.Options, string(PageNum))
		return pageNum, nil
	}
	return 0, ErrWrongValueFormat
}

func GetResultsPerPage(options *searchEnginePB.SearchOptions) (int64, error) {
	optionValues, exists := options.Options[string(ResultsPerPage)]
	if !exists {
		return 0, ErrNoOption
	}
	if len(optionValues.GetValues()) > 0 {
		resultsPerPage := optionValues.GetValues()[0]
		pageNum, convErr := strconv.ParseInt(resultsPerPage, 10, 64)
		if convErr != nil {
			// responseTemplates.SendErrorMessage(w, ErrWrongQueryParam, http.StatusBadRequest)
			return 0, convErr
		}
		// delete(options.Options, string(ResultsPerPage))
		return pageNum, nil
	}
	return 0, ErrNoOption
}
