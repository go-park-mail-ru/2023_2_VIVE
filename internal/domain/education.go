package domain

type Education struct {
	Level          string `json:"level,omitempty"`
	Major          string `json:"major,omitempty"`
	Institution    string `json:"institution,omitempty"`
	GraduationYear string `json:"grad_year,omitempty"`
}
