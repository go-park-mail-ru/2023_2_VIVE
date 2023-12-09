package repository

import (
	"context"

	"github.com/gorilla/websocket"
)

type ConnectionStorage map[int64]*websocket.Conn

type INotificationRepository interface {
	SaveConn(ctx context.Context, connection *websocket.Conn) error
	GetConn(ctx context.Context, userID int64) (*websocket.Conn, error)
}

type InMemoryNotificationRepository struct {
	storage ConnectionStorage
}

func NewInMemoryNotificationRepository() INotificationRepository {
	return &InMemoryNotificationRepository{
		storage: make(ConnectionStorage),
	}
}

func (repo *InMemoryNotificationRepository) SaveConn(ctx context.Context, connection *websocket.Conn) error {
	// TODO: saves connection to memory
	return nil
}

func (repo *InMemoryNotificationRepository) GetConn(ctx context.Context, userID int64) (*websocket.Conn, error) {
	// TODO: get connection from memory
	return nil, nil
}
