package repository

import (
	"HnH/internal/repository/mock"
	"HnH/pkg/serverErrors"
)

type ISessionRepository interface {
	AddSession(sessionID string, userID int) error
	DeleteSession(sessionID string) error
	ValidateSession(sessionID string) error
	GetUserIdBySession(sessionID string) (int, error)
}

type psqlSessionRepository struct {
	sessionStorage *mock.Sessions
}

func NewPsqlSessionRepository(sessions *mock.Sessions) ISessionRepository {
	return &psqlSessionRepository{
		sessionStorage: sessions,
	}
}

func (p *psqlSessionRepository) AddSession(sessionID string, userID int) error {
	p.sessionStorage.SessionsList.Store(sessionID, userID)
	return nil
}

func (p *psqlSessionRepository) DeleteSession(sessionID string) error {
	_, exist := p.sessionStorage.SessionsList.Load(sessionID)

	if !exist {
		return serverErrors.AUTH_REQUIRED
	}

	p.sessionStorage.SessionsList.Delete(sessionID)
	return nil
}

func (p *psqlSessionRepository) ValidateSession(sessionID string) error {
	_, ok := p.sessionStorage.SessionsList.Load(sessionID)

	if !ok {
		return serverErrors.INVALID_COOKIE
	}

	return nil
}

func (p *psqlSessionRepository) GetUserIdBySession(sessionID string) (int, error) {
	userID, ok := p.sessionStorage.SessionsList.Load(sessionID)
	if !ok {
		return 0, serverErrors.AUTH_REQUIRED
	}

	return userID.(int), nil
}
