package usecase

import (
	notificationsPB "HnH/services/notifications/api/proto"
	repository "HnH/services/notifications/internal/repository/inMemory"
	"context"

	"github.com/gorilla/websocket"
)

type INotificationUseCase interface {
	SendNotification(ctx context.Context, message *notificationsPB.NotificationMessage) error
	SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error
}

type NotificationUseCase struct {
	repo repository.INotificationRepository
}

func NewNotificationUseCase(repo repository.INotificationRepository) INotificationUseCase {
	return &NotificationUseCase{
		repo: repo,
	}
}

func (u *NotificationUseCase) SendNotification(ctx context.Context, message *notificationsPB.NotificationMessage) error {
	userID := message.UserId
	conn, err := u.repo.GetConn(ctx, userID)
	if err != nil {
		return err
	}

	err = conn.WriteJSON(message)
	if err != nil {
		conn.Close()
		u.repo.DeleteConn(ctx, userID)
		return err
	}

	return nil
}

func (u *NotificationUseCase) SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error {
	err := u.repo.SaveConn(ctx, userID, connection)
	if err != nil {
		return err
	}
	return nil
}
