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
