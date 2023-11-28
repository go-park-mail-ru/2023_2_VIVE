package domain

import (
	"HnH/services/searchEngineService/searchEnginePB"
	"time"
)

type EmploymentType string

const (
	NoneEmployment EmploymentType = "none"
	FullTime       EmploymentType = "full-time"
	PartTime       EmploymentType = "part-time"
	OneTime        EmploymentType = "one-time"
	Volunteering   EmploymentType = "volunteering"
	Internship     EmploymentType = "internship"
)

type DbVacancy struct {
	ID               int            `json:"id"`
	EmployerID       int            `json:"employer_id"`
	VacancyName      string         `json:"name"`
	Description      string         `json:"description,omitempty"`
	SalaryLowerBound *int           `json:"salary_lower_bound,omitempty"`
	SalaryUpperBound *int           `json:"salary_upper_bound,omitempty"`
	Employment       EmploymentType `json:"employment,omitempty"`
	Experience       ExperienceTime `json:"experience"`
	EducationType    EducationLevel `json:"education_type,omitempty"`
	Location         *string        `json:"location,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (vac *DbVacancy) ToAPI() *ApiVacancy {
	res := ApiVacancy{
		ID:               vac.ID,
		EmployerID:       vac.EmployerID,
		VacancyName:      vac.VacancyName,
		Description:      vac.Description,
		SalaryLowerBound: vac.SalaryLowerBound,
		SalaryUpperBound: vac.SalaryUpperBound,
		Experience:       vac.Experience,
		EducationType:    vac.EducationType,
		Employment:       vac.Employment,
		Location:         vac.Location,
		CreatedAt:        vac.CreatedAt,
		UpdatedAt:        vac.UpdatedAt,
	}

	// if vac.ExperienceUpperBound == nil {
	// 	if vac.ExperienceLowerBound == nil {
	// 		res.Experience = None
	// 	} else if *vac.ExperienceLowerBound == 6 {
	// 		res.Experience = SixMoreYears
	// 	}
	// } else if *vac.ExperienceLowerBound == 0 && *vac.ExperienceUpperBound == 0 {
	// 	res.Experience = NoExperience
	// } else if *vac.ExperienceLowerBound == 1 && *vac.ExperienceUpperBound == 3 {
	// 	res.Experience = OneThreeYears
	// } else if *vac.ExperienceLowerBound == 3 && *vac.ExperienceUpperBound == 6 {
	// 	res.Experience = ThreeSixYears
	// }

	return &res
}

type ApiVacancy struct {
	ID               int            `json:"id"`
	EmployerID       int            `json:"employer_id"`
	VacancyName      string         `json:"name"`
	SalaryLowerBound *int           `json:"salary_lower_bound,omitempty"`
	SalaryUpperBound *int           `json:"salary_upper_bound,omitempty"`
	Experience       ExperienceTime `json:"experience"`
	Employment       EmploymentType `json:"employment,omitempty"`
	EducationType    EducationLevel `json:"education_type,omitempty"`
	Location         *string        `json:"location,omitempty"`
	Description      string         `json:"description,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (vac *ApiVacancy) ToDb() *DbVacancy {
	res := DbVacancy{
		ID:               vac.ID,
		VacancyName:      vac.VacancyName,
		Description:      vac.Description,
		SalaryLowerBound: vac.SalaryLowerBound,
		SalaryUpperBound: vac.SalaryUpperBound,
		EducationType:    vac.EducationType,
		Experience:       vac.Experience,
		Employment:       vac.Employment,
		Location:         vac.Location,
		CreatedAt:        vac.CreatedAt,
		UpdatedAt:        vac.UpdatedAt,
	}

	// switch vac.Experience {
	// case NoExperience:
	// 	experienceLowerBound := 0
	// 	experienceUpperBound := 0
	// 	res.ExperienceLowerBound = &experienceLowerBound
	// 	res.ExperienceUpperBound = &experienceUpperBound
	// case OneThreeYears:
	// 	experienceLowerBound := 1
	// 	experienceUpperBound := 3
	// 	res.ExperienceLowerBound = &experienceLowerBound
	// 	res.ExperienceUpperBound = &experienceUpperBound

	// case ThreeSixYears:
	// 	experienceLowerBound := 3
	// 	experienceUpperBound := 6
	// 	res.ExperienceLowerBound = &experienceLowerBound
	// 	res.ExperienceUpperBound = &experienceUpperBound

	// case SixMoreYears:
	// 	experienceLowerBound := 6
	// 	res.ExperienceLowerBound = &experienceLowerBound

	// case None:
	// 	res.ExperienceLowerBound = nil
	// 	res.ExperienceUpperBound = nil
	// }

	return &res
}

type ApiVacancyCount struct {
	Count     int64        `json:"count"`
	Vacancies []ApiVacancy `json:"list"`
}

type ApiMetaVacancy struct {
	Filters   []*searchEnginePB.Filter `json:"filters,omitempty"`
	Vacancies ApiVacancyCount          `json:"vacancies"`
}
