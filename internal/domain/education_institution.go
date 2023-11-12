package domain

type DbEducationInstitution struct {
	ID             int    `json:"id"`
	CvID           int    `json:"cv_id"`
	Name           string `json:"name"`
	MajorField     string `json:"major_field"`
	GraduationYear string `json:"graduation_year"`
}

type ApiEducationInstitutionFromCV struct {
	Name           string `json:"name"`
	MajorField     string `json:"major_field"`
	GraduationYear string `json:"graduation_year"`
}

func (apiEdInst *ApiEducationInstitutionFromCV) ToDb() DbEducationInstitution {
	return DbEducationInstitution{
		Name:           apiEdInst.Name,
		MajorField:     apiEdInst.MajorField,
		GraduationYear: apiEdInst.GraduationYear,
	}
}
