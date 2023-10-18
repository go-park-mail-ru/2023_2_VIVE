package repository

import (
	"HnH/internal/domain"
	"HnH/internal/repository/mock"
	"HnH/pkg/authUtils"
	"HnH/pkg/serverErrors"
)

type psqlUserRepository struct {
	userStorage *mock.Users
}

func NewPsqlUserRepository(users *mock.Users) *psqlUserRepository {
	return &psqlUserRepository{
		userStorage: users,
	}
}

func (p *psqlUserRepository) checkPassword(user *domain.User) error {
	actualUserIndex, ok := p.userStorage.EmailToUser.Load(user.Email)

	if !ok {
		return serverErrors.NO_DATA_FOUND
	}

	actualUser := p.userStorage.UsersList[actualUserIndex.(int)]

	hashedPass := authUtils.GetHash(user.Password)

	if hashedPass != actualUser.Password {
		return serverErrors.INCORRECT_CREDENTIALS
	}

	return nil
}

func (p *psqlUserRepository) checkRole(user *domain.User) error {
	actualUserIndex, ok := p.userStorage.EmailToUser.Load(user.Email)

	if !ok {
		return serverErrors.NO_DATA_FOUND
	}

	actualUser := p.userStorage.UsersList[actualUserIndex.(int)]

	if user.Type != actualUser.Type {
		return serverErrors.INCORRECT_ROLE
	}

	return nil
}

func (p *psqlUserRepository) CheckUser(user *domain.User) error {
	passwordStatus := p.checkPassword(user)
	if passwordStatus != nil {
		return passwordStatus
	}

	roleStatus := p.checkRole(user)
	if roleStatus != nil {
		return roleStatus
	}

	return nil
}

func (p *psqlUserRepository) AddUser(user *domain.User) error {
	_, exist := p.userStorage.EmailToUser.Load(user.Email)

	if exist {
		return serverErrors.ACCOUNT_ALREADY_EXISTS
	}

	hashedPass := authUtils.GetHash(user.Password)

	p.userStorage.Mu.Lock()

	defer p.userStorage.Mu.Unlock()

	p.userStorage.CurrentID++
	user.ID = p.userStorage.CurrentID
	user.Password = hashedPass

	p.userStorage.UsersList = append(mock.UserDB.UsersList, user)

	p.userStorage.EmailToUser.Store(user.Email, len(mock.UserDB.UsersList)-1)
	p.userStorage.IdToUser.Store(user.ID, len(mock.UserDB.UsersList)-1)

	return nil
}

func (p *psqlUserRepository) GetUserInfo(userID int) (*domain.User, error) {
	userIndex, exist := p.userStorage.IdToUser.Load(userID)
	if !exist {
		return nil, serverErrors.NO_DATA_FOUND
	}

	user := p.userStorage.UsersList[userIndex.(int)]

	user.Password = ""

	return user, nil
}

func (p *psqlUserRepository) GetUserIdByEmail(email string) (int, error) {
	userIndex, exist := p.userStorage.EmailToUser.Load(email)
	if !exist {
		return 0, serverErrors.INVALID_EMAIL
	}

	user := p.userStorage.UsersList[userIndex.(int)]

	return user.ID, nil
}
