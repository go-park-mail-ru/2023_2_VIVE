package domain

type Vacancy struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CompanyName string `json:"company_name"`
	Description string `json:"description,omitempty"`
	Salary      int    `json:"salary,omitempty"`
}

type VacancyRepository interface {
	GetVacancies() ([]Vacancy, error)
}
