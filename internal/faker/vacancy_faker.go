package main

import (
	"HnH/app"
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/pkg/authUtils"
	"HnH/pkg/contextUtils"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"

	//"os"
	"reflect"
	"strconv"

	"time"

	"github.com/go-faker/faker/v4"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type FakeEmployer struct {
	FirstName    string
	LastName     string
	Email        string
	EmployerName string
	Description  string
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

type EducationLevelList struct {
	levels []domain.EducationLevel
	mu     *sync.RWMutex
}

var edu = EducationLevelList{
	levels: []domain.EducationLevel{domain.Nothing, domain.Secondary, domain.SecondarySpecial,
		domain.IncompleteHigher, domain.Higher, domain.Bachelor,
		domain.Master, domain.PhDJunior, domain.PhD},
	mu: &sync.RWMutex{},
}

type VacancyFaker struct {
	idToEmployer *sync.Map
	idToVacancy  *sync.Map

	userRepo    psql.IUserRepository
	vacancyRepo psql.IVacancyRepository
	skillRepo   psql.ISkillRepository

	DB *sql.DB

	logger *logrus.Entry
}

func NewVacancyFaker(conn *sql.DB, log *logrus.Entry) *VacancyFaker {
	pUserRepo := psql.NewPsqlUserRepository(conn)
	pVacancyRepo := psql.NewPsqlVacancyRepository(conn)
	pSkillRepo := psql.NewPsqlSkillRepository(conn)

	newFaker := &VacancyFaker{
		idToEmployer: &sync.Map{},
		idToVacancy:  &sync.Map{},
		userRepo:     pUserRepo,
		vacancyRepo:  pVacancyRepo,
		skillRepo:    pSkillRepo,
		DB:           conn,
		logger:       log,
	}

	return newFaker
}

func (vFaker *VacancyFaker) PushEmployers() error {
	ctx := context.Background()
	ctxLogger := context.WithValue(ctx, contextUtils.LOGGER_KEY, vFaker.logger)

	pusher := func(key any, value any) bool {
		employer, ok := value.(FakeEmployer)
		if !ok {
			fmt.Printf("CAST ERROR:\n%v\n\n\n", value)
			return false
		}

		user := &domain.ApiUser{
			Email:                   employer.Email,
			Password:                "Qwerty123",
			FirstName:               employer.FirstName,
			LastName:                employer.LastName,
			Type:                    domain.Employer,
			OrganizationName:        employer.EmployerName,
			OrganizationDescription: employer.Description,
		}

		addStatus := vFaker.userRepo.AddUser(ctxLogger, user, authUtils.GenerateHash)
		if addStatus != nil {
			fmt.Printf("user addition error: %v\n\n\n", addStatus)
			return true
		}

		return true
	}

	vFaker.idToEmployer.Range(pusher)

	return nil
}

func (vFaker *VacancyFaker) PushVacancies() error {
	ctx := context.Background()
	ctxLogger := context.WithValue(ctx, contextUtils.LOGGER_KEY, vFaker.logger)

	emailToEmpID, err := vFaker.getNewEmployersIDs()
	if err != nil {
		return err
	}

	pusher := func(key any, value any) bool {
		vacancy, ok := value.(FakeVacancy)
		if !ok {
			fmt.Printf("CAST ERROR:\n%v\n\n\n", value)
			return false
		}

		hhEmpID := vacancy.EmployerID
		employer, ok := vFaker.idToEmployer.Load(hhEmpID)
		if !ok {
			fmt.Printf("employer by hh id not found\n\n\n")
			return false
		}

		castedEmployer, ok := employer.(FakeEmployer)
		if !ok {
			fmt.Printf("CAST ERROR:\n%v\n\n\n", employer)
			return false
		}

		empEmail := castedEmployer.Email

		empID, ok := emailToEmpID.Load(empEmail)
		if !ok {
			fmt.Printf("employer id by email not found\n\n\n")
			return false
		}

		castedEmpID, ok := empID.(int)
		if !ok {
			fmt.Printf("CAST ERROR:\n%v\n\n\n", empID)
			return false
		}

		vacToAdd := &domain.DbVacancy{
			EmployerID:       castedEmpID,
			VacancyName:      vacancy.VacancyName,
			Description:      vacancy.Description,
			SalaryLowerBound: vacancy.SalaryLowerBound,
			SalaryUpperBound: vacancy.SalaryUpperBound,
			Employment:       vacancy.Employment,
			Experience:       vacancy.Experience,
			EducationType:    vacancy.EducationType,
			Location:         vacancy.Location,
		}

		vacancyID, addStatus := vFaker.vacancyRepo.AddVacancy(ctxLogger, castedEmpID, vacToAdd)
		if addStatus != nil {
			fmt.Printf("vacancy addition error: %v\n\n\n", addStatus)
			return true
		}

		addSkillsErr := vFaker.skillRepo.AddSkillsByVacID(ctxLogger, vacancyID, vacancy.Skills)
		if addSkillsErr != nil {
			fmt.Printf("skills addition error: %v\n\n\n", addSkillsErr)
			return true
		}

		return true
	}

	vFaker.idToVacancy.Range(pusher)

	return nil
}

func (vFaker *VacancyFaker) getNewEmployersIDs() (*sync.Map, error) {
	emails := []string{}
	mu := &sync.Mutex{}

	emailCollector := func(key any, value any) bool {
		employer, ok := value.(FakeEmployer)
		if !ok {
			fmt.Printf("CAST ERROR:\n%v\n\n\n", value)
			return false
		}

		defer mu.Unlock()
		mu.Lock()

		emails = append(emails, employer.Email)

		return true
	}

	vFaker.idToEmployer.Range(emailCollector)

	rows, err := vFaker.DB.Query(`SELECT u.email, emp.id
								FROM hnh_data.employer emp 
								JOIN hnh_data.user_profile u ON emp.user_id = u.id
								WHERE u.email = ANY($1)`, pq.Array(emails))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	emailToEmpID := &sync.Map{}

	for rows.Next() {
		var empID int
		var empEmail string

		err = rows.Scan(&empEmail, &empID)
		if err != nil {
			return nil, err
		}

		empEmail = strings.TrimSpace(empEmail)
		emailToEmpID.Store(empEmail, empID)
	}

	return emailToEmpID, nil
}

func (vFaker *VacancyFaker) GetData(vacsPerCategory int, categories ...string) error {
	jsonRes, err := GETWithQuery("https://api.hh.ru/vacancies/", 0, 100, "golang")
	if err != nil {
		return err
	}

	for _, value := range jsonRes["items"].([]interface{}) {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		var vacToAdd FakeVacancy

		vacID, err := getStringField(item, "id")
		if err != nil {
			continue
		}

		name, err := getStringField(item, "name")
		if err != nil {
			continue
		}
		vacToAdd.VacancyName = name

		salary, ok := item["salary"].(map[string]interface{})
		if salary == nil {
			vacToAdd.SalaryLowerBound = nil
			vacToAdd.SalaryUpperBound = nil
		} else {
			if !ok {
				continue
			}

			salaryFrom := salary["from"]
			salaryTo := salary["to"]

			if salaryFrom == nil {
				vacToAdd.SalaryLowerBound = nil
			} else {
				salaryFromFloat, ok := salaryFrom.(float64)
				if !ok {
					continue
				}

				salaryInt := int(salaryFromFloat)

				vacToAdd.SalaryLowerBound = &salaryInt
			}

			if salaryTo == nil {
				vacToAdd.SalaryUpperBound = nil
			} else {
				salaryToFloat, ok := salaryTo.(float64)
				if !ok {
					continue
				}

				salaryInt := int(salaryToFloat)

				vacToAdd.SalaryUpperBound = &salaryInt
			}
		}

		address, ok := item["address"].(map[string]interface{})
		if address == nil {
			area, areaOk := item["area"].(map[string]interface{})
			if area == nil {
				vacToAdd.Location = nil
			} else {
				if !areaOk {
					continue
				}

				areaName, areaNameOk := area["name"].(string)
				if !areaNameOk {
					continue
				}

				vacToAdd.Location = &areaName
			}
		} else {
			if !ok {
				continue
			}

			city, ok := address["city"].(string)
			if !ok {
				continue
			}
			vacToAdd.Location = &city
		}

		experienceID, err := getNestedField(item, "experience", "id")
		if err != nil {
			continue
		}
		vacToAdd.Experience = getDomainExp(experienceID)

		employmentID, err := getNestedField(item, "employment", "id")
		if err != nil {
			continue
		}
		vacToAdd.Employment = getDomainEmp(employmentID)

		var employerToAdd FakeEmployer

		employer, ok := item["employer"].(map[string]interface{})
		if !ok {
			continue
		}

		employerID, err := getStringField(employer, "id")
		if err != nil {
			continue
		}
		vacToAdd.EmployerID = employerID

		employerName, err := getStringField(employer, "name")
		if err != nil {
			continue
		}
		employerToAdd.EmployerName = employerName

		logoUrl, err := getNestedField(employer, "logo_urls", "240")
		if err != nil {
			continue
		}
		employerToAdd.LogoURL = logoUrl

		employerUrl, err := getStringField(employer, "url")
		if err != nil {
			continue
		}

		employerJsonRes, err := GET(employerUrl)
		if err != nil {
			continue
		}

		employerDesc, err := getStringField(employerJsonRes, "description")
		if err != nil {
			continue
		}
		employerToAdd.Description = employerDesc

		vacToAdd.EducationType = getRandomEducationLevel()

		vacancyUrl, err := getStringField(item, "url")
		if err != nil {
			continue
		}

		vacancyJsonRes, err := GET(vacancyUrl)
		if err != nil {
			continue
		}

		vacancyDesc, err := getStringField(vacancyJsonRes, "description")
		if err != nil {
			continue
		}
		vacToAdd.Description = vacancyDesc

		skills, err := getSkills(vacancyJsonRes)
		if err != nil {
			continue
		}
		vacToAdd.Skills = skills

		name, surname, err := getRussianName()
		if err != nil {
			continue
		}

		employerToAdd.FirstName = name
		employerToAdd.LastName = surname

		employerToAdd.Email = strings.ToLower(employerToAdd.FirstName) + "_" + strings.ToLower(employerToAdd.LastName) + "@mail.ru"

		vFaker.idToVacancy.Store(vacID, vacToAdd)
		vFaker.idToEmployer.Store(employerID, employerToAdd)
	}

	return nil
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

func GETWithQuery(endpoint string, pageNum, perPage int, skill string) (map[string]interface{}, error) {
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}

	request.URL.RawQuery = url.Values{
		"text":     {skill},
		"per_page": {strconv.Itoa(perPage)},
		"page":     {strconv.Itoa(pageNum)},
	}.Encode()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading error: %v", err)
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %v", err)
	}

	return jsonRes, nil
}

func GET(endpoint string) (map[string]interface{}, error) {
	response, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading error: %v", err)
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %v", err)
	}

	return jsonRes, nil
}

func getNestedField(item map[string]interface{}, mapFieldName, targetFieldName string) (string, error) {
	mapField, ok := item[mapFieldName].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("cast error: %s", mapFieldName)
	}

	targetField, ok := mapField[targetFieldName].(string)
	if !ok {
		return "", fmt.Errorf("cast error: %s -> %s", mapFieldName, targetFieldName)
	}

	return targetField, nil
}

func getStringField(item map[string]interface{}, targetFieldName string) (string, error) {
	targetField, ok := item[targetFieldName].(string)
	if !ok {
		return "", fmt.Errorf("cast error: %s", targetFieldName)
	}

	return targetField, nil
}

func getRandomEducationLevel() domain.EducationLevel {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	defer edu.mu.RUnlock()
	edu.mu.RLock()

	eduLevel := r.Intn(len(edu.levels))

	return edu.levels[eduLevel]
}

func getSkills(vacancyJson map[string]interface{}) ([]string, error) {
	keySkills, ok := vacancyJson["key_skills"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("cast error: key skills")
	}

	result := make([]string, 0, len(keySkills))

	for _, keySkill := range keySkills {
		skillName, ok := keySkill.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cast error: key skill")
		}

		skill, ok := skillName["name"].(string)
		if !ok {
			return nil, fmt.Errorf("cast error: skill name")
		}

		result = append(result, skill)
	}

	return result, nil
}

func getRussianName() (firstName, lastName string, err error) {
	person := faker.GetPerson()

	now := time.Now().Unix()
	if now%2 == 0 {
		name, err := person.RussianFirstNameMale(reflect.Value{})
		if err != nil {
			return "", "", err
		}

		surname, err := person.RussianLastNameMale(reflect.Value{})
		if err != nil {
			return "", "", err
		}

		firstName = name.(string)
		lastName = surname.(string)
	} else {
		name, err := person.RussianFirstNameFemale(reflect.Value{})
		if err != nil {
			return "", "", err
		}

		surname, err := person.RussianLastNameFemale(reflect.Value{})
		if err != nil {
			return "", "", err
		}

		firstName = name.(string)
		lastName = surname.(string)
	}

	return firstName, lastName, nil
}

func main() {
	db, err := app.GetPostgres()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	logEntry := logrus.NewEntry(logger)

	vacFaker := NewVacancyFaker(db, logEntry)

	err = vacFaker.GetData(1, "писька")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = vacFaker.PushEmployers()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = vacFaker.PushVacancies()
	if err != nil {
		fmt.Println(err)
		return
	}

	vacFaker.idToVacancy.Range(func(key, value any) bool { fmt.Println(key, value); return true })
	fmt.Printf("\n\n\n")
	vacFaker.idToEmployer.Range(func(key, value any) bool { fmt.Println(key, value); return true })

	/*jsonRes, err := GETWithQuery("https://api.hh.ru/vacancies/", 0, 1, "golang")
	if err != nil {
		fmt.Println(err)
		return
	}

	idToVac := make(map[string]FakeVacancy, 1000)
	idToEmployer := make(map[string]FakeEmployer)

	for _, value := range jsonRes["items"].([]interface{}) {
		item, ok := value.(map[string]interface{})
		if !ok {
			fmt.Println("Cast error")
			return
		}

		var vacToAdd FakeVacancy

		vacID, err := getStringField(item, "id")
		if err != nil {
			fmt.Println(err)
			return
		}

		name, err := getStringField(item, "name")
		if err != nil {
			fmt.Println(err)
			return
		}
		vacToAdd.VacancyName = name

		salary, ok := item["salary"].(map[string]interface{})
		if salary == nil {
			vacToAdd.SalaryLowerBound = nil
			vacToAdd.SalaryUpperBound = nil
		} else {
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
		}

		address, ok := item["address"].(map[string]interface{})
		if address == nil {
			area, areaOk := item["area"].(map[string]interface{})
			if area == nil {
				vacToAdd.Location = nil
			} else {
				if !areaOk {
					fmt.Println("Cast error: area")
					return
				}

				areaName, areaNameOk := area["name"].(string)
				if !areaNameOk {
					fmt.Println("Cast error: area name")
					return
				}

				vacToAdd.Location = &areaName
			}
		} else {
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
		}

		experienceID, err := getNestedField(item, "experience", "id")
		if err != nil {
			fmt.Println(err)
			return
		}
		vacToAdd.Experience = getDomainExp(experienceID)

		employmentID, err := getNestedField(item, "employment", "id")
		if err != nil {
			fmt.Println(err)
			return
		}
		vacToAdd.Employment = getDomainEmp(employmentID)

		var employerToAdd FakeEmployer

		employer, ok := item["employer"].(map[string]interface{})
		if !ok {
			fmt.Println("Cast error: employer")
			return
		}

		employerID, err := getStringField(employer, "id")
		if err != nil {
			fmt.Println(err)
			return
		}
		vacToAdd.EmployerID = employerID

		employerName, err := getStringField(employer, "name")
		if err != nil {
			fmt.Println(err)
			return
		}
		employerToAdd.EmployerName = employerName

		logoUrl, err := getNestedField(employer, "logo_urls", "240")
		if err != nil {
			fmt.Println(err)
			return
		}
		employerToAdd.LogoURL = logoUrl

		employerUrl, err := getStringField(employer, "url")
		if err != nil {
			fmt.Println(err)
			return
		}

		employerJsonRes, err := GET(employerUrl)
		if err != nil {
			fmt.Println(err)
			return
		}

		employerDesc, err := getStringField(employerJsonRes, "description")
		if err != nil {
			fmt.Println(err)
			return
		}
		employerToAdd.Description = employerDesc

		vacToAdd.EducationType = getRandomEducationLevel()

		vacancyUrl, err := getStringField(item, "url")
		if err != nil {
			fmt.Println(err)
			return
		}

		vacancyJsonRes, err := GET(vacancyUrl)
		if err != nil {
			fmt.Println(err)
			return
		}

		vacancyDesc, err := getStringField(vacancyJsonRes, "description")
		if err != nil {
			fmt.Println(err)
			return
		}
		vacToAdd.Description = vacancyDesc

		skills, err := getSkills(vacancyJsonRes)
		if err != nil {
			fmt.Println(err)
			return
		}
		vacToAdd.Skills = skills

		idToVac[vacID] = vacToAdd
		idToEmployer[employerID] = employerToAdd

		fmt.Println(idToVac)
		fmt.Println(idToEmployer)

		params, err := url.Parse(employerToAdd.LogoURL)
		segments := strings.Split(params.Path, "/")
		fileName := segments[len(segments)-1]

		resp, err := http.Get(employerToAdd.LogoURL)
		if err != nil {
			fmt.Println(err)
			return
		}

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		name, surname, err := getRussianName()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(name, surname)
	}*/
}
