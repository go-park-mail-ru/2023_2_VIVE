package usecase

import (
	"HnH/internal/repository/grpc"
	notificationsPB "HnH/services/notifications/api/proto"
	"context"
)

type INotificationUsecase interface {
	GetUsersNotifications(ctx context.Context, user_id int64) (*notificationsPB.UserNotifications, error)
	DeleteUsersNotifications(ctx context.Context, user_id int64) error
}

type NotificationUsecase struct {
	notificationsRepo grpc.INotificationRepository
}

func NewNotificationUsecase(notificationRepo grpc.INotificationRepository) INotificationUsecase {
	return &NotificationUsecase{
		notificationsRepo: notificationRepo,
	}
}

func (u *NotificationUsecase) GetUsersNotifications(ctx context.Context, user_id int64) (*notificationsPB.UserNotifications, error) {
	return u.notificationsRepo.GetUserNotifications(ctx, user_id)
}

func (u *NotificationUsecase) DeleteUsersNotifications(ctx context.Context, user_id int64) error {
	return u.notificationsRepo.DeleteUserNotifications(ctx, user_id)
}
