package domain

import "time"

type Vacancy struct {
	ID                     int       `json:"id"`
	Employer_id            int       `json:"employer_id"`
	VacancyName            string    `json:"name"`
	Description            string    `json:"description,omitempty"`
	Salary_lower_bound     int       `json:"salary_lower_bound,omitempty"`
	Salary_upper_bound     int       `json:"salary_upper_bound,omitempty"`
	Employment             string    `json:"employment,omitempty"`
	Experience_lower_bound int       `json:"experience_lower_bound,omitempty"`
	Experience_upper_bound int       `json:"experience_upper_bound,omitempty"`
	EducationType          string    `json:"education_type,omitempty"`
	Location               string    `json:"location,omitempty"`
	Created_at             time.Time `json:"created_at"`
	Updated_at             time.Time `json:"updated_at"`
	// CompanyName            string    `json:"company_name"`
	// CompanyID              int       `json:"company_id,omitempty"`
	// Salary                 int       `json:"salary,omitempty"`
}
