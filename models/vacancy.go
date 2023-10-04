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
	Mu          *sync.RWMutex
}

var VacancyDB = vacancies{
	VacancyList: []Vacancy{
		{
			ID:          1,
			Name:        "C++ developer",
			CompanyName: "VK",
			Description: "Middle C++ developer in Mail.ru team",
			Salary:      250000,
		},
		{
			ID:          2,
			Name:        "Go developer",
			CompanyName: "VK",
			Description: "Golang junior developer without any experience",
			Salary:      100000,
		},
		{
			ID:          3,
			Name:        "HR",
			CompanyName: "Yandex",
			Description: "Human resources specialist",
			Salary:      70000,
		},
		{
			ID:          4,
			Name:        "Frontend developer",
			CompanyName: "Google",
			Description: "Middle Frontend developer, JavaScript, HTML, Figma",
			Salary:      500000,
		},
		{
			ID:          5,
			Name:        "Project Manager",
			CompanyName: "VK",
			Description: "Experienced specialist in IT-management",
			Salary:      200000,
		}},

	Mu: &sync.RWMutex{},
}
