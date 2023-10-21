package domain

type CV struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id,omitempty"`
	CVName      string      `json:"name,omitempty"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	DateOfBirth string      `json:"date_of_bitrh,omitempty"`
	Skills      []string    `json:"skills,omitempty"`
	Contact     string      `json:"contact,omitempty"`
	PhoneNumber string      `json:"phone_number,omitempty"`
	Location    string      `json:"location,omitempty"`
	Status      string      `json:"status,omitempty"`
	Grades      []Education `json:"grades,omitempty"`
	Languages   []Language  `json:"languages,omitempty"`
	MainText    string      `json:"main_text,omitempty"`
}
