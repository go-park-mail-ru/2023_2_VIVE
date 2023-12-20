package domain

//easyjson:json
type ApiResponse struct {
	Id               int    `json:"id"`
	VacancyName      string `json:"vacancy_name"`
	VacancyID        int    `json:"vacancy_id"`
	OrganizationName string `json:"organization_name"`
	EmployerID       int    `json:"employer_id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

//easyjson:json
type ApiResponseSlice []ApiResponse
