package models

import "sync"

var CurrentID = 0

type Role string

const (
	Applicant Role = "applicant"
	Employer  Role = "employer"
)

var Users = sync.Map{}

type User struct {
	ID        int    `json:"id,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Type      Role   `json:"role,omitempty"`
}
