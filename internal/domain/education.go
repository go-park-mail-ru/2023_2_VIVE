package domain

type EducationLevel string

const (
	Secondary        EducationLevel = "secondary"         //среднее
	SecondarySpecial EducationLevel = "secondary_special" // средне профессиональное
	IncompleteHigher EducationLevel = "incomplete_higher" // неоконченное высшее
	Higher           EducationLevel = "higher"            // высшее
	Bachelor         EducationLevel = "bachelor"          // бакалавр
	Master           EducationLevel = "master"            // магистр
	PhDJunior        EducationLevel = "phd_junior"        // кандидат наук
	PhD              EducationLevel = "phd"               // доктор наук
)

type Education struct {
	Level          string `json:"level,omitempty"`
	Major          string `json:"major,omitempty"`
	Institution    string `json:"institution,omitempty"`
	GraduationYear string `json:"grad_year,omitempty"`
}
