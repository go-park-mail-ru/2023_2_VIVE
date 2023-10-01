package models

import "sync"

type Role string

const (
	Applicant Role = "applicant"
	Employer  Role = "employer"
)

type Users struct {
	UsersList []*User
	CurrentID int
	Mu        *sync.Mutex
}

var IdToUser = sync.Map{}
var EmailToUser = sync.Map{}

var UserDB = Users{
	UsersList: make([]*User, 0),
	CurrentID: 0,
	Mu:        &sync.Mutex{},
}

type User struct {
	ID        int    `json:"id,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Type      Role   `json:"role,omitempty"`
}
