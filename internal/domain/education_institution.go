package domain

type EducationInstitution struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	EducationLevel EducationLevel `json:"education_level"`
}
