package models

import (
	"crypto/sha256"
	"encoding/hex"
	"models/statuses"
	"sync"
)

type Role string

const (
	Applicant Role = "applicant"
	Employer  Role = "employer"
)

type User struct {
	ID        int    `json:"id,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Type      Role   `json:"role,omitempty"`
}

var users = &sync.Map{}

func CheckPassword(user User) statuses.Status {
	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPass := hasher.Sum(nil)

	actualPass, ok := users.Load(user.Email)

	if !ok {
		return statuses.NOT_FOUND
	}

	if hex.EncodeToString(hashedPass) != hex.EncodeToString(actualPass.([]byte)) {
		return statuses.UNAUTHORIZED
	}

	return 0
}

func AddUser(user User) statuses.Status {
	_, exist := users.Load(user.Email)

	if exist {
		return statuses.CONFLICT
	}

	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPass := hasher.Sum(nil)

	users.Store(user.Email, hashedPass)

	return 0
}
