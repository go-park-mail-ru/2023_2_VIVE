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
	sessionStorage redis.Conn
}

func NewRedisSessionRepository(conn redis.Conn) ISessionRepository {
	return &redisSessionRepository{
		sessionStorage: conn,
	}
}

func (p *redisSessionRepository) AddSession(sessionID string, userID int, expiryUnixSeconds int64) error {
	sessionKey := "sessions:" + sessionID

	result, err := redis.String(p.sessionStorage.Do("SET", sessionKey, userID, "EXAT", expiryUnixSeconds))
	if err != nil {
		return err
	} else if result != "OK" {
		return ERROR_WHILE_WRITING
	}

	return nil
}

func (p *redisSessionRepository) DeleteSession(sessionID string) error {
	sessionKey := "sessions:" + sessionID

	_, err := redis.Int(p.sessionStorage.Do("DEL", sessionKey))
	if err != nil {
		return err
	}

	return nil
}

func (p *redisSessionRepository) ValidateSession(sessionID string) error {
	sessionKey := "sessions:" + sessionID

	result, err := redis.Int(p.sessionStorage.Do("EXISTS", sessionKey))
	if result == 0 {
		return serverErrors.NO_SESSION
	} else if err != nil {
		return err
	}

	return nil
}

func (p *redisSessionRepository) GetUserIdBySession(sessionID string) (int, error) {
	sessionKey := "sessions:" + sessionID

	userID, err := redis.Int(p.sessionStorage.Do("GET", sessionKey))
	if err != nil {
		return 0, err
	}

	return userID, nil
}
