package domain

import "time"

type Status string

const (
	Searching    Status = "searching"
	NotSearching Status = "not searching"
)

type CV struct {
	ID             int       `json:"id"`
	ApplicantID    int       `json:"applicant_id"`
	ProfessionName string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	Status         Status    `json:"status,omitempty"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
	// UserID         int       `json:"user_id,omitempty"`
	// FirstName   string      `json:"first_name"`
	// LastName    string      `json:"last_name"`
	// DateOfBirth string      `json:"date_of_bitrh,omitempty"`
	// Skills      []string    `json:"skills,omitempty"`
	// Contact     string      `json:"contact,omitempty"`
	// PhoneNumber string      `json:"phone_number,omitempty"`
	// Location    string      `json:"location,omitempty"`
	// Grades      []Education `json:"grades,omitempty"`
	// Languages   []Language  `json:"languages,omitempty"`
	// MainText    string      `json:"main_text,omitempty"`
}
