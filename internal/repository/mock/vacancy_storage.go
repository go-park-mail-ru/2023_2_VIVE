package mock

import (
	"HnH/internal/domain"
	"sync"
)

type Vacancies struct {
	VacancyList []domain.Vacancy
	Mu          *sync.RWMutex
	CurrentID   int
}

// var VacancyDB = Vacancies{
// 	VacancyList: []domain.Vacancy{
// 		{
// 			ID:          1,
// 			VacancyName: "C++ developer",
// 			CompanyName: "VK",
// 			Description: "Middle C++ developer in Mail.ru team",
// 			Salary:      250000,
// 		},
// 		{
// 			ID:          2,
// 			VacancyName: "Go developer",
// 			CompanyName: "VK",
// 			Description: "Golang junior developer without any experience",
// 			Salary:      100000,
// 		},
// 		{
// 			ID:          3,
// 			VacancyName: "HR",
// 			CompanyName: "Yandex",
// 			Description: "Human resources specialist",
// 			Salary:      70000,
// 		},
// 		{
// 			ID:          4,
// 			VacancyName: "Frontend developer",
// 			CompanyName: "Google",
// 			Description: "Middle Frontend developer, JavaScript, HTML, Figma",
// 			Salary:      500000,
// 		},
// 		{
// 			ID:          5,
// 			VacancyName: "Project Manager",
// 			CompanyName: "VK",
// 			Description: "Experienced specialist in IT-management",
// 			Salary:      200000,
// 		}},

// 	CurrentID: 5,

// 	Mu: &sync.RWMutex{},
// }
