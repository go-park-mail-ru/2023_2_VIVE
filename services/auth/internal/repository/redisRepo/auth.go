package redisRepo

import (
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	pb "HnH/services/auth/authPB"
	"context"

	"github.com/gomodule/redigo/redis"
)

type IAuthRepository interface {
	AddSession(ctx context.Context, sessionID string, userID int64, expiryUnixSeconds int64) error
	DeleteSession(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, sessionID string) error
	GetUserIdBySession(ctx context.Context, sessionID string) (*pb.UserID, error)
}

type redisAuthRepository struct {
	sessionStorage *redis.Pool
}

func NewRedisAuthRepository(conn *redis.Pool) IAuthRepository {
	return &redisAuthRepository{
		sessionStorage: conn,
	}
}

func (p *redisAuthRepository) AddSession(ctx context.Context, sessionID string, userID int64, expiryUnixSeconds int64) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("adding new session in redis")

	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	result, err := redis.String(connection.Do("SET", sessionKey, userID, "EXAT", expiryUnixSeconds))
	if err != nil {
		contextLogger.WithField("err", err.Error()).Info("Error while adding session")
		return err
	} else if result != "OK" {
		return ERROR_WHILE_WRITING
	}

	return nil
}

func (p *redisAuthRepository) DeleteSession(ctx context.Context, sessionID string) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("deleting session from redis")

	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	_, err := redis.Int(connection.Do("DEL", sessionKey))
	if err != nil {
		contextLogger.WithField("err", err.Error()).Info("Error while deleting session")
		return err
	}

	return nil
}

func (p *redisAuthRepository) ValidateSession(ctx context.Context, sessionID string) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("checking session by its id in redis")

	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	result, err := redis.Int(connection.Do("EXISTS", sessionKey))
	if result == 0 {
		return serverErrors.NO_SESSION
	} else if err != nil {
		contextLogger.WithField("err", err.Error()).Info("Error while checking session")
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	return nil
}

func (p *redisAuthRepository) GetUserIdBySession(ctx context.Context, sessionID string) (*pb.UserID, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting user id by session id from redis")

	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	userID, err := redis.Int64(connection.Do("GET", sessionKey))
	if err != nil {
		contextLogger.WithField("err", err.Error()).Info("Error while checking session")
		return &pb.UserID{}, ENTITY_NOT_FOUND
	}

	return &pb.UserID{UserId: userID}, nil
}
