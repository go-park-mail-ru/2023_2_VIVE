package domain

import "net/http"

const SpecialChars = `~!?@#$%^&*_-+()[]{}></\|"'.,:;`

type Role string

const (
	Applicant Role = "applicant"
	Employer  Role = "employer"
)

func (r Role) IsRole() bool {
	return r == Applicant || r == Employer
}

type User struct {
	ID        int    `json:"id,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Type      Role   `json:"role,omitempty"`
}

type UserRepository interface {
	CheckPassword(user *User) error
	CheckRole(user *User) error
	ValidatePassword(password string) error
	CheckUser(user *User) error
	AddUser(user *User) error
	GetUserInfo(cookie *http.Cookie) (*User, error)
}
