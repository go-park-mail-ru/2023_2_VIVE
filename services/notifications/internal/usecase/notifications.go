package usecase

import (
	notificationsPB "HnH/services/notifications/api/proto"
	repositoryIM "HnH/services/notifications/internal/repository/inMemory"
	repositoryPSQL "HnH/services/notifications/internal/repository/psql"
	"context"

	"github.com/gorilla/websocket"
)

type INotificationUseCase interface {
	SendNotification(ctx context.Context, message *notificationsPB.NotificationMessage) error
	SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error
	GetUsersNotifications(ctx context.Context, userID int64) ([]*notificationsPB.NotificationMessage, error)
	DeleteUsersNotifications(ctx context.Context, userID int64) error
}

type NotificationUseCase struct {
	connRepo         repositoryIM.IConnectionRepository
	notificationRepo repositoryPSQL.INotificationRepository
}

func NewNotificationUseCase(connRepo repositoryIM.IConnectionRepository, notificationRepo repositoryPSQL.INotificationRepository) INotificationUseCase {
	return &NotificationUseCase{
		connRepo:         connRepo,
		notificationRepo: notificationRepo,
	}
}

func (u *NotificationUseCase) SendNotification(ctx context.Context, message *notificationsPB.NotificationMessage) error {
	addErr := u.notificationRepo.AddNotification(ctx, message)
	if addErr != nil {
		return addErr
	}

	userID := message.UserId
	conn, err := u.connRepo.GetConn(ctx, userID)
	if err != nil {
		return err
	}

	err = conn.WriteJSON(message)
	if err != nil {
		conn.Close()
		u.connRepo.DeleteConn(ctx, userID)
		return err
	}

	return nil
}

func (u *NotificationUseCase) SaveConn(ctx context.Context, userID int64, connection *websocket.Conn) error {
	err := u.connRepo.SaveConn(ctx, userID, connection)
	if err != nil {
		return err
	}
	return nil
}

func (u *NotificationUseCase) GetUsersNotifications(ctx context.Context, userID int64) ([]*notificationsPB.NotificationMessage, error) {
	return u.notificationRepo.GetUsersNotifications(ctx, userID)
}

func (u *NotificationUseCase) DeleteUsersNotifications(ctx context.Context, userID int64) error {
	return u.notificationRepo.DeleteUsersNotifications(ctx, userID)
}
