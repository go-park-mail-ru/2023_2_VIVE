package models

type Vacancy struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Salary      int    `json:"salary,omitempty"`
}
