package usecase

import (
	notificationsPB "HnH/services/notifications/api/proto"
	repository "HnH/services/notifications/internal/repository/inMemory"
	"context"

	"github.com/gorilla/websocket"
)

type INotificationUseCase interface {
    SendNotification(ctx context.Context, notification *notificationsPB.NotificationMessage) error
	SaveConn(ctx context.Context, connection *websocket.Conn) error
}

type NotificationUseCase struct {
	repo *repository.INotificationRepository
}

func NewNotificationUseCase(repo *repository.INotificationRepository) INotificationUseCase {
	return &NotificationUseCase{
		repo: repo,
	}
}

func (u *NotificationUseCase) SendNotification(ctx context.Context, notification *notificationsPB.NotificationMessage) error {
	// TODO: send notification to frontend
	return nil
}

func (u *NotificationUseCase) SaveConn(ctx context.Context, connection *websocket.Conn) error {
	// TODO: save new connection
	return nil
}
