package domain

import (
	"HnH/services/searchEngineService/searchEnginePB"
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
		ApplicantID:    cv.ApplicantID,
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

//easyjson:json
type ApiCV struct {
	ID                    int                        `json:"id"`
	ApplicantID           int                        `json:"applicant_id"`
	LastName              string                     `json:"last_name" pdf:"header"`             // фамилия
	FirstName             string                     `json:"first_name" pdf:"header"`            // имя
	MiddleName            *string                    `json:"middle_name,omitempty" pdf:"header"` // отчество
	Gender                Gender                     `json:"gender" pdf:"content,Основная информация,Пол,gender"`
	Location              *string                    `json:"city,omitempty" pdf:"content,Основная информация,Расположение"`
	Birthday              *string                    `json:"birthday,omitempty" pdf:"content,Основная информация,День рождения,dd.mm.yyyy"`
	ProfessionName        string                     `json:"profession_name" pdf:"content,Желаемая должность"`
	EducationLevel        EducationLevel             `json:"education_level"`
	Status                Status                     `json:"status,omitempty"`
	EducationInstitutions []*ApiEducationInstitution `json:"institutions" pdf:"content,Образование"`
	Experience            []*ApiExperience           `json:"companies" pdf:"content,Опыт работы"`
	Skills                []string                   `json:"skills" pdf:"content,Ключевые навыки"`
	Description           *string                    `json:"description,omitempty" pdf:"content,О себе"`
	AvatarURL             string                     `json:"avatar_url"`
	CreatedAt             time.Time                  `json:"created_at,omitempty"`
	UpdatedAt             time.Time                  `json:"updated_at,omitempty"`
}

//easyjson:json
type ApiCVSlice []ApiCV

func (cv *ApiCV) ToDb() *DbCV {
	return &DbCV{
		ID:             cv.ID,
		ApplicantID:    cv.ApplicantID,
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

//easyjson:json
type ApiCVCount struct {
	Count int64   `json:"count"`
	CVs   []ApiCV `json:"list"`
}

//easyjson:json
type ApiMetaCV struct {
	Filters []*searchEnginePB.Filter `json:"filters,omitempty"`
	CVs     ApiCVCount               `json:"cvs"`
}
