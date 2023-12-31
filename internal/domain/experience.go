package domain

type ExperienceTime string

const (
	None          ExperienceTime = "none"
	NoExperience  ExperienceTime = "no_experience"
	OneThreeYears ExperienceTime = "one_three_years"
	ThreeSixYears ExperienceTime = "three_six_years"
	SixMoreYears  ExperienceTime = "six_more_years"
)

type DbExperience struct {
	ID               int     `json:"id"`
	CvID             int     `json:"cv_id"`
	OrganizationName string  `json:"organization_name"`
	Position         string  `json:"position"`
	Description      string  `json:"description"`
	StartDate        string  `json:"start_date"`
	EndDate          *string `json:"end_date,omitempty"`
}

func (apiExp *DbExperience) ToAPI() *ApiExperience {
	return &ApiExperience{
		ID:               apiExp.ID,
		CvID:             apiExp.CvID,
		OrganizationName: apiExp.OrganizationName,
		JobPosition:      apiExp.Position,
		Description:      apiExp.Description,
		StartDate:        apiExp.StartDate,
		EndDate:          apiExp.EndDate,
	}
}

//easyjson:json
type ApiExperience struct {
	ID               int     `json:"id"`
	CvID             int     `json:"cv_id"`
	OrganizationName string  `json:"name" pdf:"header"`
	JobPosition      string  `json:"job_position" pdf:"content,-,Должность"`
	Description      string  `json:"description" pdf:"content,-,Описание"`
	StartDate        string  `json:"start_date" pdf:"content,-,Начало работы,mm.yyyy"`
	EndDate          *string `json:"end_date,omitempty" pdf:"content,-,Конец работы,mm.yyyy"`
}

func (apiExp *ApiExperience) ToDb() *DbExperience {
	return &DbExperience{
		ID:               apiExp.ID,
		CvID:             apiExp.CvID,
		OrganizationName: apiExp.OrganizationName,
		Position:         apiExp.JobPosition,
		Description:      apiExp.Description,
		StartDate:        apiExp.StartDate,
		EndDate:          apiExp.EndDate,
	}
}
