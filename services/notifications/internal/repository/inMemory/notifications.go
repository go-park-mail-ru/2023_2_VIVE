package repository

import (
	"HnH/services/notifications/pkg/serviceErrors"
	"context"

	"github.com/gorilla/websocket"
)

type ConnectionStorage map[int64]*websocket.Conn

type INotificationRepository interface {
	SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error
	GetConn(ctx context.Context, userID int64) (*websocket.Conn, error)
	DeleteConn(ctx context.Context, userID int64) 
}

type InMemoryNotificationRepository struct {
	storage ConnectionStorage
}

func NewInMemoryNotificationRepository() INotificationRepository {
	return &InMemoryNotificationRepository{
		storage: make(ConnectionStorage),
	}
}

func (repo *InMemoryNotificationRepository) SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error {
	// TODO: saves connection to memory
	_, exists := repo.storage[userID]
	if exists {
		return serviceErrors.ErrConnAlreadyExists
	}
	repo.storage[userID] = connection
	return nil
}

func (repo *InMemoryNotificationRepository) GetConn(ctx context.Context, userID int64) (*websocket.Conn, error) {
	conn, exists := repo.storage[userID]
	if !exists {
		return nil, serviceErrors.ErrNoConn
	}
	// TODO: get connection from memory
	return conn, nil
}

func (repo *InMemoryNotificationRepository) DeleteConn(ctx context.Context, userID int64)  {
	conn, exists := repo.storage[userID]
	if !exists {
		return
	}
	conn.Close()
	delete(repo.storage, userID)
}
