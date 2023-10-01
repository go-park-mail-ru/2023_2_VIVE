package models

import "sync"

type Vacancy struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CompanyName string `json:"company_name"`
	Description string `json:"description,omitempty"`
	Salary      int    `json:"salary,omitempty"`
}

type vacancies struct {
	VacancyList []Vacancy
	Mu          *sync.Mutex
}

var VacancyDB = vacancies{
	VacancyList: []Vacancy{
		{
			ID:          1,
			Name:        "C++ developer",
			CompanyName: "VK",
			Description: "Middle C++ developer in Mail.ru team",
			Salary:      2500,
		},
		{
			ID:          2,
			Name:        "Go developer",
			CompanyName: "VK",
			Description: "Golang junior developer without any experience",
			Salary:      1000,
		},
		{
			ID:          3,
			Name:        "HR",
			CompanyName: "Yandex",
			Description: "Human resources specialist",
			Salary:      700,
		},
		{
			ID:          4,
			Name:        "Frontend developer",
			CompanyName: "Google",
			Description: "Middle Frontend developer, JavaScript, HTML, Figma",
			Salary:      5000,
		},
		{
			ID:          5,
			Name:        "Project Manager",
			CompanyName: "VK",
			Description: "Experienced specialist in IT-management",
			Salary:      2000,
		}},

	Mu: &sync.Mutex{},
}
