package domain

import (
	// "fmt"
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
	ID                     int             `json:"id"`
	Employer_id            int             `json:"employer_id"`
	VacancyName            string          `json:"name"`
	Description            string          `json:"description,omitempty"`
	Salary_lower_bound     *int            `json:"salary_lower_bound,omitempty"`
	Salary_upper_bound     *int            `json:"salary_upper_bound,omitempty"`
	Employment             *EmploymentType `json:"employment,omitempty"`
	Experience_lower_bound *int            `json:"experience_lower_bound,omitempty"`
	Experience_upper_bound *int            `json:"experience_upper_bound,omitempty"`
	EducationType          EducationLevel  `json:"education_type,omitempty"`
	Location               *string         `json:"location,omitempty"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
}

func (vac *DbVacancy) ToAPI() *ApiVacancy {
	res := ApiVacancy{
		ID:                 vac.ID,
		VacancyName:        vac.VacancyName,
		Description:        vac.Description,
		Salary_lower_bound: vac.Salary_lower_bound,
		Salary_upper_bound: vac.Salary_upper_bound,
		EducationType:      vac.EducationType,
		Employment:         vac.Employment,
		Location:           vac.Location,
		CreatedAt:          vac.CreatedAt,
		UpdatedAt:          vac.UpdatedAt,
	}

	if *vac.Experience_lower_bound == 1 && *vac.Experience_upper_bound == 3 {
		res.Experience = OneThreeYears
	} else if *vac.Experience_lower_bound == 3 && *vac.Experience_upper_bound == 6 {
		res.Experience = ThreeSixYears
	} else if *vac.Experience_lower_bound == 6 {
		res.Experience = SixMoreYears
	} else {
		res.Experience = None
	}

	return &res
}

type ApiVacancy struct {
	ID                 int             `json:"id"`
	VacancyName        string          `json:"name"`
	Salary_lower_bound *int            `json:"salary_lower_bound,omitempty"`
	Salary_upper_bound *int            `json:"salary_upper_bound,omitempty"`
	Experience         ExperienceTime  `json:"experience,omitempty"`
	Employment         *EmploymentType `json:"employment,omitempty"`
	EducationType      EducationLevel  `json:"education_type,omitempty"`
	Location           *string         `json:"location,omitempty"`
	Description        string          `json:"description,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

func (vac *ApiVacancy) ToDb() *DbVacancy {
	res := DbVacancy{
		ID:                 vac.ID,
		VacancyName:        vac.VacancyName,
		Description:        vac.Description,
		Salary_lower_bound: vac.Salary_lower_bound,
		Salary_upper_bound: vac.Salary_upper_bound,
		EducationType:      vac.EducationType,
		Employment:         vac.Employment,
		Location:           vac.Location,
		CreatedAt:          vac.CreatedAt,
		UpdatedAt:          vac.UpdatedAt,
	}

	switch vac.Experience {
	case OneThreeYears:
		experienceLowerBound := 1
		experienceUpperBound := 6
		res.Experience_lower_bound = &experienceLowerBound
		res.Experience_upper_bound = &experienceUpperBound

	case ThreeSixYears:
		experienceLowerBound := 3
		experienceUpperBound := 6
		res.Experience_lower_bound = &experienceLowerBound
		res.Experience_upper_bound = &experienceUpperBound

	case SixMoreYears:
		experienceLowerBound := 6
		res.Experience_lower_bound = &experienceLowerBound

	case None:
		res.Experience_lower_bound = nil
		res.Experience_upper_bound = nil
	}

	return &res
}

// type ApiVacancyUpdate struct {
// 	VacancyName        string         `json:"vacancy_name"`
// 	Salary_lower_bound int            `json:"salary_lower_bound,omitempty"`
// 	Salary_upper_bound int            `json:"salary_upper_bound,omitempty"`
// 	Experience         ExperienceTime `json:"experience,omitempty"`
// 	Employment         EmploymentType `json:"employment,omitempty"`
// 	EducationType      EducationLevel `json:"education_type,omitempty"`
// 	Location           string         `json:"location,omitempty"`
// 	Description        string         `json:"description,omitempty"`
// }

// func (vac *ApiVacancyUpdate) ToDb() *DbVacancy {
// 	// fmt.Printf("apiVac.VacancyName: %v\n", apiVac.VacancyName)
// 	// fmt.Printf("apiVac.Description: %v\n", apiVac.Description)
// 	// fmt.Printf("&apiVac.Salary_lower_bound: %v\n", &apiVac.Salary_lower_bound)
// 	// fmt.Printf("&apiVac.Salary_upper_bound: %v\n", &apiVac.Salary_upper_bound)
// 	// fmt.Printf("apiVac.EducationType: %v\n", apiVac.EducationType)
// 	// fmt.Printf("&apiVac.Employment: %v\n", &apiVac.Employment)
// 	// fmt.Printf("&apiVac.Location: %v\n", &apiVac.Location)

// 	res := DbVacancy{
// 		VacancyName:        vac.VacancyName,
// 		Description:        vac.Description,
// 		Salary_lower_bound: &vac.Salary_lower_bound,
// 		Salary_upper_bound: &vac.Salary_upper_bound,
// 		EducationType:      vac.EducationType,
// 		Employment:         &vac.Employment,
// 		Location:           &vac.Location,
// 	}

// 	// fmt.Printf("res.VacancyName: %v\n", res.VacancyName)
// 	// fmt.Printf("res.Description: %v\n", res.Description)
// 	// fmt.Printf("res.Salary_lower_bound: %v\n", *res.Salary_lower_bound)
// 	// fmt.Printf("res.Salary_upper_bound: %v\n", *res.Salary_upper_bound)
// 	// fmt.Printf("res.EducationType: %v\n", res.EducationType)
// 	// fmt.Printf("res.Employment: %v\n", *res.Employment)
// 	// fmt.Printf("res.Location: %v\n", *res.Location)

// 	switch vac.Experience {
// 	case OneThreeYears:
// 		experienceLowerBound := 1
// 		experienceUpperBound := 6
// 		res.Experience_lower_bound = &experienceLowerBound
// 		res.Experience_upper_bound = &experienceUpperBound

// 	case ThreeSixYears:
// 		experienceLowerBound := 3
// 		experienceUpperBound := 6
// 		res.Experience_lower_bound = &experienceLowerBound
// 		res.Experience_upper_bound = &experienceUpperBound

// 	case SixMoreYears:
// 		experienceLowerBound := 6
// 		res.Experience_lower_bound = &experienceLowerBound

// 	case None:
// 		res.Experience_lower_bound = nil
// 		res.Experience_upper_bound = nil
// 	}

// 	return &res
// }
