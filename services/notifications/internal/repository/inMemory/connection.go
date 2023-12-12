package repository

import (
	"HnH/pkg/contextUtils"
	"HnH/services/notifications/pkg/serviceErrors"
	"context"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ConnectionStorage map[int64]*websocket.Conn

type IConnectionRepository interface {
	SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error
	GetConn(ctx context.Context, userID int64) (*websocket.Conn, error)
	DeleteConn(ctx context.Context, userID int64)
}

type InMemoryConnectionRepository struct {
	storage ConnectionStorage
}

func NewInMemoryConnectionRepository() IConnectionRepository {
	return &InMemoryConnectionRepository{
		storage: make(ConnectionStorage),
	}
}

func (repo *InMemoryConnectionRepository) SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("saving new incoming connection")
	_, exists := repo.storage[userID]
	if exists {
		return serviceErrors.ErrConnAlreadyExists
	}
	repo.storage[userID] = connection
	return nil
}

func (repo *InMemoryConnectionRepository) GetConn(ctx context.Context, userID int64) (*websocket.Conn, error) {
	conn, exists := repo.storage[userID]
	if !exists {
		return nil, serviceErrors.ErrNoConn
	}
	// TODO: get connection from memory
	return conn, nil
}

func (repo *InMemoryConnectionRepository) DeleteConn(ctx context.Context, userID int64) {
	conn, exists := repo.storage[userID]
	if !exists {
		return
	}
	conn.Close()
	delete(repo.storage, userID)
}
