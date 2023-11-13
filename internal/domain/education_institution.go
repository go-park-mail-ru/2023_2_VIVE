package domain

type DbEducationInstitution struct {
	ID             int    `json:"id"`
	CvID           int    `json:"cv_id"`
	Name           string `json:"name"`
	MajorField     string `json:"major_field"`
	GraduationYear string `json:"graduation_year"`
}

type ApiEducationInstitution struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	MajorField     string `json:"major_field"`
	GraduationYear string `json:"graduation_year"`
}

func (apiEdInst *ApiEducationInstitution) ToDb() DbEducationInstitution {
	return DbEducationInstitution{
		ID:             apiEdInst.ID,
		Name:           apiEdInst.Name,
		MajorField:     apiEdInst.MajorField,
		GraduationYear: apiEdInst.GraduationYear,
	}
}
