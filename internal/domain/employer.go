package domain

type DbEmployer struct {
	UserID                  int    `json:"user_id"`
	OrganizationName        string `json:"organization_name"`
	OrganizationDescription string `json:"organization_description"`
}

type EmployerInfo struct {
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	CompanyName string       `json:"organization_name"`
	Vacancies   []ApiVacancy `json:"vacancies,omitempty"`
}
