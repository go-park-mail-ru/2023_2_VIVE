package redisRepo

import (
	"HnH/pkg/serverErrors"

	"github.com/gomodule/redigo/redis"
)

type ISessionRepository interface {
	AddSession(sessionID string, userID int, expiryUnixSeconds int64) error
	DeleteSession(sessionID string) error
	ValidateSession(sessionID string) error
	GetUserIdBySession(sessionID string) (int, error)
}

type redisSessionRepository struct {
	sessionStorage *redis.Pool
}

func NewRedisSessionRepository(conn *redis.Pool) ISessionRepository {
	return &redisSessionRepository{
		sessionStorage: conn,
	}
}

func (p *redisSessionRepository) AddSession(sessionID string, userID int, expiryUnixSeconds int64) error {
	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	result, err := redis.String(connection.Do("SET", sessionKey, userID, "EXAT", expiryUnixSeconds))
	if err != nil {
		return err
	} else if result != "OK" {
		return ERROR_WHILE_WRITING
	}

	return nil
}

func (p *redisSessionRepository) DeleteSession(sessionID string) error {
	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	_, err := redis.Int(connection.Do("DEL", sessionKey))
	if err != nil {
		return err
	}

	return nil
}

func (p *redisSessionRepository) ValidateSession(sessionID string) error {
	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	result, err := redis.Int(connection.Do("EXISTS", sessionKey))
	if result == 0 {
		return serverErrors.NO_SESSION
	} else if err != nil {
		return err
	}

	return nil
}

func (p *redisSessionRepository) GetUserIdBySession(sessionID string) (int, error) {
	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	userID, err := redis.Int(connection.Do("GET", sessionKey))
	if err != nil {
		return 0, err
	}

	return userID, nil
}
