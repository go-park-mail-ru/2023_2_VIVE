package domain

import "HnH/pkg/nullTypes"

type Role string

const (
	Applicant Role = "applicant"
	Employer  Role = "employer"
)

func (r Role) IsRole() bool {
	return r == Applicant || r == Employer
}

type User struct {
	ID          int                  `json:"id,omitempty"`
	Email       string               `json:"email"`
	Password    string               `json:"password,omitempty"`
	FirstName   string               `json:"first_name,omitempty"`
	LastName    string               `json:"last_name,omitempty"`
	Birthday    nullTypes.NullString `json:"birthday,omitempty"`
	PhoneNumber nullTypes.NullString `json:"phone_number,omitempty"`
	Location    nullTypes.NullString `json:"location,omitempty"`
	Type        Role                 `json:"role,omitempty"`
}

type UserUpdate struct {
	Email       string               `json:"email,omitempty"`
	FirstName   string               `json:"first_name,omitempty"`
	LastName    string               `json:"last_name,omitempty"`
	Birthday    nullTypes.NullString `json:"birthday,omitempty"`
	PhoneNumber nullTypes.NullString `json:"phone_number,omitempty"`
	Location    nullTypes.NullString `json:"location,omitempty"`
	Password    string               `json:"password"`
	NewPassword string               `json:"new_password,omitempty"`
}
