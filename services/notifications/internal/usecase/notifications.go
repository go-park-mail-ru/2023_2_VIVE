package usecase

import (
	"HnH/pkg/contextUtils"
	notificationsPB "HnH/services/notifications/api/proto"
	repositoryIM "HnH/services/notifications/internal/repository/inMemory"
	repositoryPSQL "HnH/services/notifications/internal/repository/psql"
	"HnH/services/notifications/pkg/serviceErrors"
	"context"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"message": message,
	}).
		Info("sending notification")
	addErr := u.notificationRepo.AddNotification(ctx, message)
	if addErr != nil {
		return addErr
	}

	userID := message.UserId
	conn, err := u.connRepo.GetConn(ctx, userID)
	if errors.Is(err, serviceErrors.ErrNoConn) || errors.Is(err, serviceErrors.ErrInvalidConnection) {
		return nil
	} else if err != nil {
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
	if err != nil && err != serviceErrors.ErrConnAlreadyExists {
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
