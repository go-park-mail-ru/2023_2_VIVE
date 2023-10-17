package mock

import (
	"HnH/internal/domain"
	"sync"
)

const SpecialChars = `~!?@#$%^&*_-+()[]{}></\|"'.,:;`

type Role string

const (
	Applicant Role = "applicant"
	Employer  Role = "employer"
)

type Users struct {
	UsersList   []*domain.User
	IdToUser    sync.Map
	EmailToUser sync.Map
	CurrentID   int
	Mu          *sync.Mutex
}

var UserDB = Users{
	UsersList:   make([]*domain.User, 0),
	IdToUser:    sync.Map{},
	EmailToUser: sync.Map{},
	CurrentID:   0,
	Mu:          &sync.Mutex{},
}
