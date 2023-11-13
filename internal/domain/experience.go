package domain

type ExperienceTime string

const (
	None          ExperienceTime = "none"
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

type ApiExperience struct {
	ID               int     `json:"id"`
	OrganizationName string  `json:"name"`
	JobPosition      string  `json:"job_position"`
	Description      string  `json:"description"`
	StartDate        string  `json:"start_date"`
	EndDate          *string `json:"end_date,omitempty"`
}

func (apiExp *ApiExperience) ToDb() DbExperience {
	dbExp := DbExperience{
		ID:               apiExp.ID,
		OrganizationName: apiExp.OrganizationName,
		Position:         apiExp.JobPosition,
		Description:      apiExp.Description,
		StartDate:        apiExp.StartDate,
		EndDate:          apiExp.EndDate,
	}

	return dbExp
}
