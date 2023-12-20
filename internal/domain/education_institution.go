package domain

type DbEducationInstitution struct {
	ID             int    `json:"id"`
	CvID           int    `json:"cv_id"`
	Name           string `json:"name"`
	MajorField     string `json:"major_field"`
	GraduationYear string `json:"graduation_year"`
}

func (dbInst *DbEducationInstitution) ToAPI() *ApiEducationInstitution {
	return &ApiEducationInstitution{
		ID:             dbInst.ID,
		CvID:           dbInst.CvID,
		Name:           dbInst.Name,
		MajorField:     dbInst.MajorField,
		GraduationYear: dbInst.GraduationYear,
	}
}

//easyjson:json
type ApiEducationInstitution struct {
	ID             int    `json:"id"`
	CvID           int    `json:"cv_id"`
	Name           string `json:"name" pdf:"header"`
	MajorField     string `json:"major_field" pdf:"content,-,Специальность"`
	GraduationYear string `json:"graduation_year" pdf:"content,-,Год выпуска"`
}

func (apiEdInst *ApiEducationInstitution) ToDb() *DbEducationInstitution {
	return &DbEducationInstitution{
		ID:             apiEdInst.ID,
		CvID:           apiEdInst.CvID,
		Name:           apiEdInst.Name,
		MajorField:     apiEdInst.MajorField,
		GraduationYear: apiEdInst.GraduationYear,
	}
}
