package domain

import (
	"time"
)

type EmploymentType string

const (
	FullTime     EmploymentType = "full-time"
	PartTime     EmploymentType = "part-time"
	OneTime      EmploymentType = "one-time"
	Volunteering EmploymentType = "volunteering"
	Internship   EmploymentType = "internship"
)

type DbVacancy struct {
	ID                     int            `json:"id"`
	Employer_id            int            `json:"employer_id"`
	VacancyName            string         `json:"name"`
	Description            string         `json:"description,omitempty"`
	Salary_lower_bound     *int           `json:"salary_lower_bound,omitempty"`
	Salary_upper_bound     *int           `json:"salary_upper_bound,omitempty"`
	Employment             *string        `json:"employment,omitempty"`
	Experience_lower_bound *int           `json:"experience_lower_bound,omitempty"`
	Experience_upper_bound *int           `json:"experience_upper_bound,omitempty"`
	EducationType          EducationLevel `json:"education_type,omitempty"`
	Location               *string        `json:"location,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

type ApiVacancy struct {
	VacancyName        string         `json:"vacancy_name"`
	Salary_lower_bound int            `json:"salary_lower_bound,omitempty"`
	Salary_upper_bound int            `json:"salary_upper_bound,omitempty"`
	Experience         ExperienceTime `json:"experience,omitempty"`
	Employment         EmploymentType `json:"employment,omitempty"`
	Location           string         `json:"location,omitempty"`
	Description        string         `json:"description,omitempty"`
}
