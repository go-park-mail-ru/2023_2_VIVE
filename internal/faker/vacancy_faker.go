package main

import (
	"HnH/internal/domain"

	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type FakeEmployer struct {
	EmployerName string
	LogoURL      string
}

type FakeVacancy struct {
	EmployerID       string
	VacancyName      string
	SalaryLowerBound *int
	SalaryUpperBound *int
	Experience       domain.ExperienceTime
	Employment       domain.EmploymentType
	EducationType    domain.EducationLevel
	Location         *string
	Description      string
	Skills           []string
}

func getDomainExp(HHexpID string) domain.ExperienceTime {
	switch HHexpID {
	case "noExperience":
		return domain.None
	case "between1And3":
		return domain.OneThreeYears
	case "between3And6":
		return domain.ThreeSixYears
	case "moreThan6":
		return domain.SixMoreYears
	default:
		return domain.None
	}
}

func getDomainEmp(HHempID string) domain.EmploymentType {
	switch HHempID {
	case "full":
		return domain.FullTime
	case "part":
		return domain.PartTime
	case "project":
		return domain.OneTime
	case "volunteer":
		return domain.Volunteering
	case "probation":
		return domain.Internship
	default:
		return domain.NoneEmployment
	}
}

/*type ApiVacancy struct {
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
	Skills           []string       `json:"skills,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}*/

func main() {
	request, err := http.NewRequest("GET", "https://api.hh.ru/vacancies/", nil)
	if err != nil {
		fmt.Printf("Client error: %v\n", err)
		return
	}

	request.URL.RawQuery = url.Values{
		"text":     {"golang"},
		"per_page": {"1"},
		//"currency": {"RUR"},
	}.Encode()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("API error: %v\n", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Reading error: %v\n", err)
		return
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		fmt.Printf("Parsing error: %v\n", err)
		return
	}

	/*type ApiVacancy struct {
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
		Skills           []string       `json:"skills,omitempty"`
		CreatedAt        time.Time      `json:"created_at"`
		UpdatedAt        time.Time      `json:"updated_at"`
	}*/

	/*const (
		NoneEmployment EmploymentType = "none"
		FullTime       EmploymentType = "full-time"
		PartTime       EmploymentType = "part-time"
		OneTime        EmploymentType = "one-time"
		Volunteering   EmploymentType = "volunteering"
		Internship     EmploymentType = "internship"
	)*/

	/*const (
		None          ExperienceTime = "none"
		NoExperience  ExperienceTime = "no_experience"
		OneThreeYears ExperienceTime = "one_three_years"
		ThreeSixYears ExperienceTime = "three_six_years"
		SixMoreYears  ExperienceTime = "six_more_years"
	)*/

	/*const (
		Nothing          EducationLevel = "nothing"
		Secondary        EducationLevel = "secondary"         // среднее
		SecondarySpecial EducationLevel = "secondary_special" // средне профессиональное
		IncompleteHigher EducationLevel = "incomplete_higher" // неоконченное высшее
		Higher           EducationLevel = "higher"            // высшее
		Bachelor         EducationLevel = "bachelor"          // бакалавр
		Master           EducationLevel = "master"            // магистр
		PhDJunior        EducationLevel = "phd_junior"        // кандидат наук
		PhD              EducationLevel = "phd"               // доктор наук
	)*/

	idToVac := make(map[string]FakeVacancy, 1000)

	idToEmployer := make(map[string]FakeEmployer)

	edu := []domain.EducationLevel{domain.Nothing, domain.Secondary, domain.SecondarySpecial,
		domain.IncompleteHigher, domain.Higher, domain.Bachelor,
		domain.Master, domain.PhDJunior, domain.PhD}

	for _, value := range jsonRes["items"].([]interface{}) {
		item, ok := value.(map[string]interface{})
		if !ok {
			fmt.Println("Cast error")
			return
		}

		var vacToAdd FakeVacancy

		vacID, ok := item["id"].(string)
		if !ok {
			fmt.Println("Cast error: vacancy id")
			return
		}

		name, ok := item["name"].(string)
		if !ok {
			fmt.Println("Cast error: name")
			return
		}
		vacToAdd.VacancyName = name

		salary, ok := item["salary"].(map[string]interface{})
		if !ok {
			fmt.Println("Cast error: salary")
			return
		}

		salaryFrom := salary["from"]
		salaryTo := salary["to"]

		if salaryFrom == nil {
			vacToAdd.SalaryLowerBound = nil
		} else {
			salaryFromFloat, ok := salaryFrom.(float64)
			if !ok {
				fmt.Println("Cast error: salary from")
				return
			}

			salaryInt := int(salaryFromFloat)

			vacToAdd.SalaryLowerBound = &salaryInt
		}

		if salaryTo == nil {
			vacToAdd.SalaryUpperBound = nil
		} else {
			salaryToFloat, ok := salaryTo.(float64)
			if !ok {
				fmt.Println("Cast error: salary to")
				return
			}

			salaryInt := int(salaryToFloat)

			vacToAdd.SalaryUpperBound = &salaryInt
		}

		address, ok := item["address"].(map[string]interface{})
		if !ok {
			fmt.Println("Cast error: address")
			return
		}

		city, ok := address["city"].(string)
		if !ok {
			fmt.Println("Cast error: city")
			return
		}
		vacToAdd.Location = &city

		experience, ok := item["experience"].(map[string]interface{})
		if !ok {
			fmt.Println("Cast error: experience")
			return
		}

		expID, ok := experience["id"].(string)
		if !ok {
			fmt.Println("Cast error: experience ID")
			return
		}
		vacToAdd.Experience = getDomainExp(expID)

		employment, ok := item["employment"].(map[string]interface{})
		if !ok {
			fmt.Println("Cast error: employment")
			return
		}

		empID, ok := employment["id"].(string)
		if !ok {
			fmt.Println("Cast error: employment ID")
			return
		}
		vacToAdd.Employment = getDomainEmp(empID)

		employer, ok := item["employer"].(map[string]interface{})
		if !ok {
			fmt.Println("Cast error: employer")
			return
		}

		var employerToAdd FakeEmployer

		employerID, ok := employer["id"].(string)
		if !ok {
			fmt.Println("Cast error: employer ID")
			return
		}
		vacToAdd.EmployerID = employerID

		employerName, ok := employer["name"].(string)
		if !ok {
			fmt.Println("Cast error: employer name")
			return
		}
		employerToAdd.EmployerName = employerName

		logoUrls, ok := employer["logo_urls"].(map[string]interface{})
		if !ok {
			fmt.Println("Cast error: employer logo urls")
			return
		}

		logoUrl, ok := logoUrls["240"].(string)
		if !ok {
			fmt.Println("Cast error: employer logo url")
			return
		}
		employerToAdd.LogoURL = logoUrl

		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		eduLevel := r.Intn(len(edu))

		vacToAdd.EducationType = edu[eduLevel]

		verboseUrl, ok := item["url"].(string)
		if !ok {
			fmt.Println("Cast error: url")
			return
		}

		response, err := http.Get(verboseUrl)
		if err != nil {
			fmt.Printf("Verbose API error: %v\n", err)
			return
		}

		verbBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Reading error: %v\n", err)
			return
		}

		var verbJsonRes map[string]interface{}
		err = json.Unmarshal(verbBody, &verbJsonRes)
		if err != nil {
			fmt.Printf("Parsing error: %v\n", err)
			return
		}

		description, ok := verbJsonRes["description"].(string)
		if !ok {
			fmt.Println("Cast error: description")
			return
		}
		vacToAdd.Description = description

		keySkills, ok := verbJsonRes["key_skills"].([]interface{})
		if !ok {
			fmt.Println("Cast error: key skills")
			return
		}

		for _, keySkill := range keySkills {
			skillName, ok := keySkill.(map[string]interface{})
			if !ok {
				fmt.Println("Cast error: key skill")
				return
			}

			skill, ok := skillName["name"].(string)
			if !ok {
				fmt.Println("Cast error: skill name")
				return
			}

			vacToAdd.Skills = append(vacToAdd.Skills, skill)
		}

		idToVac[vacID] = vacToAdd
		idToEmployer[employerID] = employerToAdd
	}
}
