package domain

import "time"

type Status string
type Gender string

const (
	Searching    Status = "searching"
	NotSearching Status = "not searching"
)

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type DbCV struct {
	ID             int       `json:"id"`
	ApplicantID    int       `json:"applicant_id"`
	ProfessionName string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	Status         Status    `json:"status,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// FirstName      string    `json:"first_name"`
	// LastName       string    `json:"last_name"`
	// UserID         int       `json:"user_id,omitempty"`
	// DateOfBirth string      `json:"date_of_bitrh,omitempty"`
	// Skills      []string    `json:"skills,omitempty"`
	// Contact     string      `json:"contact,omitempty"`
	// PhoneNumber string      `json:"phone_number,omitempty"`
	// Location    string      `json:"location,omitempty"`
	// Grades      []Education `json:"grades,omitempty"`
	// Languages   []Language  `json:"languages,omitempty"`
	// MainText    string      `json:"main_text,omitempty"`
}

type ApiCV struct {
	FirstName                string         `json:"first_name"`            // имя
	LastName                 string         `json:"last_name"`             // фамилия
	MiddleName               string         `json:"middle_name,omitempty"` // отчество
	ProfessionName           string         `json:"profession_name"`
	Gender                   Gender         `json:"gender"`
	City                     string         `json:"city,omitempty"`
	Birthday                 string         `json:"birthday,omitempty"`
	EducationLevel           EducationLevel `json:"education_level"`
	EducationInstitutionName string         `json:"education_institution_name,omitempty"`
	Division                 string         `json:"division,omitempty"`
	MajorField               string         `json:"major_field,omitempty"`
	GraduationYear           string         `json:"graduation_year,omitempty"`
	OrganizationName         string         `json:"organization_name,omitempty"`
	JobPosition              string         `json:"job_position,omitempty"`
	StartDate                string         `json:"start_date,omitempty"`
	EndDate                  string         `json:"end_date,omitempty"`
	ExperienceDescription    string         `json:"experience_description,omitempty"`
	ApplicantDescription     string         `json:"applicant_description,omitempty"`
}
