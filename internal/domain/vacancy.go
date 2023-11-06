package domain

import (
	"HnH/pkg/nullTypes"
	"time"
)

type EmploymentType string

const (
	FullTime EmploymentType = "full-time"
	PartTime EmploymentType = "part-time"
)

type EducationType string

const (
	Higher    EducationType = "higher"
	Secondary EducationType = "secondary"
)

type Vacancy struct {
	ID                     int                  `json:"id"`
	Employer_id            int                  `json:"employer_id"`
	VacancyName            string               `json:"name"`
	Description            string               `json:"description,omitempty"`
	Salary_lower_bound     nullTypes.NullInt    `json:"salary_lower_bound,omitempty"`
	Salary_upper_bound     nullTypes.NullInt    `json:"salary_upper_bound,omitempty"`
	Employment             nullTypes.NullString `json:"employment,omitempty"`
	Experience_lower_bound nullTypes.NullInt    `json:"experience_lower_bound,omitempty"`
	Experience_upper_bound nullTypes.NullInt    `json:"experience_upper_bound,omitempty"`
	EducationType          EducationType        `json:"education_type,omitempty"`
	Location               nullTypes.NullString `json:"location,omitempty"`
	Created_at             time.Time            `json:"created_at"`
	Updated_at             time.Time            `json:"updated_at"`
	// CompanyID              int            `json:"company_id,omitempty"`
	// CompanyName            string    `json:"company_name"`
	// Salary                 int       `json:"salary,omitempty"`
}
