package usecase

import (
	psqlmock "HnH/internal/repository/mock"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	notificationsPB "HnH/services/notifications/api/proto"
	"context"
	"fmt"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserNotificationsSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	notifUsecase := usecase.NewNotificationUsecase(notifRepo)

	defer mockCtrl.Finish()

	userID := 20

	notifMess := &notificationsPB.NotificationMessage{
		UserId:    int64(userID),
		VacancyId: int64(2),
		CvId:      int64(5),
		Message:   "Responsed",
		Data:      "Successfully responsed",
	}

	notifics := &notificationsPB.UserNotifications{
		Notifications: []*notificationsPB.NotificationMessage{notifMess},
	}

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	notifRepo.EXPECT().GetUserNotifications(ctxWithID, int64(userID)).Return(notifics, nil)

	_, err := notifUsecase.GetUsersNotifications(ctxWithID, int64(userID))
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestGetUserNotificationsFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	notifUsecase := usecase.NewNotificationUsecase(notifRepo)

	defer mockCtrl.Finish()

	userID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	notifRepo.EXPECT().GetUserNotifications(ctxWithID, int64(userID)).Return(nil, serverErrors.FORBIDDEN)

	_, err := notifUsecase.GetUsersNotifications(ctxWithID, int64(userID))
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.FORBIDDEN.Error())
}

func TestDeleteUsersNotificationsSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	notifUsecase := usecase.NewNotificationUsecase(notifRepo)

	defer mockCtrl.Finish()

	userID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	notifRepo.EXPECT().DeleteUserNotifications(ctxWithID, int64(userID)).Return(nil)

	err := notifUsecase.DeleteUsersNotifications(ctxWithID, int64(userID))
	if err != nil {
		fmt.Println("Error must be nil")
		t.Fail()
	}
}

func TestDeleteUsersNotificationsFail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	notifRepo := psqlmock.NewMockINotificationRepository(mockCtrl)

	notifUsecase := usecase.NewNotificationUsecase(notifRepo)

	defer mockCtrl.Finish()

	userID := 20

	ctxWithID := context.WithValue(ctx, contextUtils.USER_ID_KEY, userID)
	notifRepo.EXPECT().DeleteUserNotifications(ctxWithID, int64(userID)).Return(serverErrors.INTERNAL_SERVER_ERROR)

	err := notifUsecase.DeleteUsersNotifications(ctxWithID, int64(userID))
	if err == nil {
		fmt.Println("Error must be not nil")
		t.Fail()
	}

	assert.Error(t, err, serverErrors.INTERNAL_SERVER_ERROR.Error())
}
