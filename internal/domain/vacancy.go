package domain

type Vacancy struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	CompanyName   string `json:"company_name"`
	Description   string `json:"description,omitempty"`
	Salary        int    `json:"salary,omitempty"`
	Employment    string `json:"employment,omitempty"`
	Experience    string `json:"experience,omitempty"`
	EducationType string `json:"education_type,omitempty"`
	Location      string `json:"location,omitempty"`
}
