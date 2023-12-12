package grpc

import (
	"HnH/pkg/contextUtils"
	notificationsPB "HnH/services/notifications/api/proto"
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type INotificationRepository interface {
	SendMessage(ctx context.Context, message *notificationsPB.NotificationMessage) error
	GetUserNotifications(ctx context.Context, userID int64) (*notificationsPB.UserNotifications, error)
	DeleteUserNotifications(ctx context.Context, userID int64) error
}

type grpcNotificationRepository struct {
	client notificationsPB.NotificationServiceClient
}

func NewGrpcNotificationRepository(client notificationsPB.NotificationServiceClient) INotificationRepository {
	return &grpcNotificationRepository{
		client: client,
	}
}

func (repo *grpcNotificationRepository) SendMessage(ctx context.Context, message *notificationsPB.NotificationMessage) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"notification_message": message,
	}).
		Info("sending notification")

	ctx = contextUtils.PutRequestIDToMetaDataCtx(ctx)
	_, err := repo.client.NotifyUser(ctx, message)
	if err != nil {
		grpcStatus := status.Convert(err)
		errMessage := grpcStatus.Message()

		errToReturn := GetErrByMessage(errMessage)

		return errToReturn
	}

	return nil
}

func (repo *grpcNotificationRepository) GetUserNotifications(ctx context.Context, userID int64) (*notificationsPB.UserNotifications, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("getting user's notifications")

	ctx = contextUtils.PutRequestIDToMetaDataCtx(ctx)
	notifications, err := repo.client.GetUserNotifications(ctx, &notificationsPB.UserID{UserId: userID})
	if err != nil {
		grpcStatus := status.Convert(err)
		errMessage := grpcStatus.Message()

		errToReturn := GetErrByMessage(errMessage)

		return nil, errToReturn
	}

	return notifications, nil
}

func (repo *grpcNotificationRepository) DeleteUserNotifications(ctx context.Context, userID int64) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("getting user's notifications")

	ctx = contextUtils.PutRequestIDToMetaDataCtx(ctx)
	_, err := repo.client.DeleteUserNotifications(ctx, &notificationsPB.UserID{UserId: userID})
	if err != nil {
		grpcStatus := status.Convert(err)
		errMessage := grpcStatus.Message()

		errToReturn := GetErrByMessage(errMessage)

		return errToReturn
	}

	return nil
}
