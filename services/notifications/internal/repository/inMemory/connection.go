package repository

import (
	"HnH/pkg/contextUtils"
	"HnH/services/notifications/pkg/serviceErrors"
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type IConnectionRepository interface {
	SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error
	GetConn(ctx context.Context, userID int64) (*websocket.Conn, error)
	DeleteConn(ctx context.Context, userID int64)
}

type InMemoryConnectionRepository struct {
	storage sync.Map
}

func NewInMemoryConnectionRepository() IConnectionRepository {
	return &InMemoryConnectionRepository{
		storage: sync.Map{},
	}
}

func (repo *InMemoryConnectionRepository) SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("saving new incoming connection")

	if value, loaded := repo.storage.LoadAndDelete(userID); loaded {
		contextLogger.Info("rewriting existing connection")
		existingConn, ok := value.(*websocket.Conn)
		if !ok {
			return serviceErrors.ErrInvalidConnection
		}
		existingConn.Close()
	}
	repo.storage.Store(userID, connection)
	return nil
}

func (repo *InMemoryConnectionRepository) GetConn(ctx context.Context, userID int64) (*websocket.Conn, error) {
	value, exists := repo.storage.Load(userID)
	if !exists {
		return nil, serviceErrors.ErrNoConn
	}
	conn, ok := value.(*websocket.Conn)
	if !ok {
		return nil, serviceErrors.ErrInvalidConnection
	}
	return conn, nil
}

func (repo *InMemoryConnectionRepository) DeleteConn(ctx context.Context, userID int64) {
	if value, loaded := repo.storage.LoadAndDelete(userID); loaded {
		conn, ok := value.(*websocket.Conn)
		if !ok {
			return
		}
		conn.Close()
	}
}
