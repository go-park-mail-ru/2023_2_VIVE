package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	notificationsPB "HnH/services/notifications/api/proto"
	"context"
)

type INotificationUsecase interface {
	GetUsersNotifications(ctx context.Context, user_id int64) (*domain.UserNotifications, error)
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

func (u *NotificationUsecase) convertNotifications(notif *notificationsPB.UserNotifications) *domain.UserNotifications {
	toReturn := &domain.UserNotifications{}

	list := notif.Notifications
	for _, note := range list {
		toAppend := domain.NotificationMessage{}

		toAppend.UserId = note.UserId
		toAppend.Data = note.Data
		toAppend.Message = note.Message
		toAppend.CreatedAt = note.CreatedAt

		toReturn.Notifications = append(toReturn.Notifications, toAppend)
	}

	return toReturn
}

func (u *NotificationUsecase) GetUsersNotifications(ctx context.Context, user_id int64) (*domain.UserNotifications, error) {
	notifications, err := u.notificationsRepo.GetUserNotifications(ctx, user_id)
	if err != nil {
		return nil, err
	}

	toReturn := u.convertNotifications(notifications)

	return toReturn, nil
}

func (u *NotificationUsecase) DeleteUsersNotifications(ctx context.Context, user_id int64) error {
	return u.notificationsRepo.DeleteUserNotifications(ctx, user_id)
}
