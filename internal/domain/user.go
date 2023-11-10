package domain

type Role string

const (
	Applicant Role = "applicant"
	Employer  Role = "employer"
)

func (r Role) IsRole() bool {
	return r == Applicant || r == Employer
}

type ApiUserReg struct {
	// ID          int     `json:"id,omitempty"`
	Email            string  `json:"email"`
	Password         string  `json:"password,omitempty"`
	FirstName        string  `json:"first_name,omitempty"`
	LastName         string  `json:"last_name,omitempty"`
	Birthday         *string `json:"birthday,omitempty"`
	PhoneNumber      *string `json:"phone_number,omitempty"`
	Location         *string `json:"location,omitempty"`
	Type             Role    `json:"role,omitempty"`
	OrganizationName string  `json:"organization_name,omitempty"`
	// AvatarPath  *string `json:"avatar,omitempty"`
}

// func (user *ApiUserReg) ToDb() *DbUser {
// 	res :=
// }

type DbUser struct {
	ID          int     `json:"id,omitempty"`
	Email       string  `json:"email"`
	Password    string  `json:"password,omitempty"`
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	Birthday    *string `json:"birthday,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Location    *string `json:"location,omitempty"`
	Type        Role    `json:"role,omitempty"`
	AvatarPath  *string `json:"avatar,omitempty"`
}

type UserUpdate struct {
	Email       string  `json:"email,omitempty"`
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	Birthday    *string `json:"birthday,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Location    *string `json:"location,omitempty"`
	Password    string  `json:"password"`
	NewPassword string  `json:"new_password,omitempty"`
}
