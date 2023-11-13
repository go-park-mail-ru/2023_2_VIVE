package domain

type Role string

const (
	Applicant Role = "applicant"
	Employer  Role = "employer"
)

func (r Role) IsRole() bool {
	return r == Applicant || r == Employer
}

type ApiUser struct {
	ID               int     `json:"id,omitempty"`
	EmployerID       *int    `json:"employer_id,omitempty"`
	ApplicantID      *int    `json:"applicant_id,omitempty"`
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

func (u *ApiUser) ToDb() *DbUser {
	return &DbUser{
		ID:          u.ID,
		Email:       u.Email,
		Password:    u.Password,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Birthday:    u.Birthday,
		PhoneNumber: u.PhoneNumber,
		Location:    u.Location,
		Type:        u.Type,
		// AvatarPath: u.,
	}
}

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

func (u *DbUser) ToAPI(empID, appID *int) *ApiUser {
	return &ApiUser{
		ID:          u.ID,
		EmployerID:  empID,
		ApplicantID: appID,
		Email:       u.Email,
		Password:    u.Password,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Birthday:    u.Birthday,
		PhoneNumber: u.PhoneNumber,
		Location:    u.Location,
		Type:        u.Type,
		// OrganizationName: u.OrganizationName,
	}
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
