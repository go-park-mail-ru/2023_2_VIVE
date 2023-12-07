package domain

type FilterName string

const (
	CityFilter          FilterName = "city"
	SalaryFilter        FilterName = "salary"
	EmploymentFilter    FilterName = "employment"
	ExperienceFilter    FilterName = "experience"
	EducationTypeFilter FilterName = "education_type"
)

type FilterType string

const (
	CheckBox       FilterType = "checkbox"
	Radio          FilterType = "radio"
	CheckBoxSearch FilterType = "checkbox_search"
	DoubleRange    FilterType = "double_range"
)
