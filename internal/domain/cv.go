package domain

import (
	"time"
)

const DATE_FORMAT = "2006-01-02"

type Status string
type Gender string
type EducationLevel string

const (
	Searching    Status = "searching"
	NotSearching Status = "not searching"
)

const (
	Male   Gender = "male"
	Female Gender = "female"
)

const (
	Nothing          EducationLevel = "nothing"
	Secondary        EducationLevel = "secondary"         // среднее
	SecondarySpecial EducationLevel = "secondary_special" // средне профессиональное
	IncompleteHigher EducationLevel = "incomplete_higher" // неоконченное высшее
	Higher           EducationLevel = "higher"            // высшее
	Bachelor         EducationLevel = "bachelor"          // бакалавр
	Master           EducationLevel = "master"            // магистр
	PhDJunior        EducationLevel = "phd_junior"        // кандидат наук
	PhD              EducationLevel = "phd"               // доктор наук
)

type DbCV struct {
	ID             int            `json:"id"`
	ApplicantID    int            `json:"applicant_id"`
	ProfessionName string         `json:"name"`
	FirstName      string         `json:"first_name"`            // имя
	LastName       string         `json:"last_name"`             // фамилия
	MiddleName     *string        `json:"middle_name,omitempty"` // отчество
	Gender         Gender         `json:"gender"`
	Birthday       *string        `json:"birthday,omitempty"`
	Location       *string        `json:"city,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Status         Status         `json:"status,omitempty"`
	EducationLevel EducationLevel `json:"education_level"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (cv *DbCV) ToAPI() *ApiCV {
	res := ApiCV{
		ID:             cv.ID,
		ProfessionName: cv.ProfessionName,
		FirstName:      cv.FirstName,
		LastName:       cv.LastName,
		MiddleName:     cv.MiddleName,
		Gender:         cv.Gender,
		Location:       cv.Location,
		Description:    cv.Description,
		Status:         cv.Status,
		EducationLevel: cv.EducationLevel,
		CreatedAt:      cv.CreatedAt,
		UpdatedAt:      cv.UpdatedAt,
		// Birthday:       cv.Birthday,
	}

	if cv.Birthday != nil {
		birthday, birthdayErr := time.Parse(time.RFC3339, *cv.Birthday)
		if birthdayErr == nil {
			birthdayStr := birthday.Format(time.DateOnly)
			res.Birthday = &birthdayStr
		}
	}

	return &res
}

type ApiCV struct {
	ID                    int                       `json:"id"`
	FirstName             string                    `json:"first_name"`            // имя
	LastName              string                    `json:"last_name"`             // фамилия
	MiddleName            *string                   `json:"middle_name,omitempty"` // отчество
	ProfessionName        string                    `json:"profession_name"`
	Gender                Gender                    `json:"gender"`
	Location              *string                   `json:"city,omitempty"`
	Birthday              *string                   `json:"birthday,omitempty"`
	EducationLevel        EducationLevel            `json:"education_level"`
	Status                Status                    `json:"status,omitempty"`
	EducationInstitutions []ApiEducationInstitution `json:"institutions"`
	Experience            []ApiExperience           `json:"companies"`
	Description           *string                   `json:"description,omitempty"`
	CreatedAt             time.Time                 `json:"created_at,omitempty"`
	UpdatedAt             time.Time                 `json:"updated_at,omitempty"`
}

func (cv *ApiCV) ToDb() *DbCV {
	return &DbCV{
		ID:             cv.ID,
		ProfessionName: cv.ProfessionName,
		FirstName:      cv.FirstName,
		LastName:       cv.LastName,
		MiddleName:     cv.MiddleName,
		Gender:         cv.Gender,
		Birthday:       cv.Birthday,
		Location:       cv.Location,
		Description:    cv.Description,
		Status:         cv.Status,
		EducationLevel: cv.EducationLevel,
		CreatedAt:      cv.CreatedAt,
		UpdatedAt:      cv.UpdatedAt,
	}
}

type ApiMetaCV struct {
	Count int64   `json:"count"`
	CVs   []ApiCV `json:"list"`
}
