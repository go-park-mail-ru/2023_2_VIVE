package repository

import (
	"HnH/internal/repository/psql"
	"HnH/pkg/contextUtils"
	notificationsPB "HnH/services/notifications/api/proto"
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type INotificationRepository interface {
	AddNotification(ctx context.Context, message *notificationsPB.NotificationMessage) error
	GetUsersNotifications(ctx context.Context, userID int64) ([]*notificationsPB.NotificationMessage, error)
	DeleteUsersNotifications(ctx context.Context, userID int64) error
}

type PsqlNotificationRepository struct {
	db *sql.DB
}

func NewPsqlNotificationRepository(db *sql.DB) INotificationRepository {
	return &PsqlNotificationRepository{
		db: db,
	}
}

func (repo *PsqlNotificationRepository) AddNotification(ctx context.Context, message *notificationsPB.NotificationMessage) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"message": message,
	}).
		Info("adding new notification message")

	query := `INSERT
			INTO
				hnh_data.notification (user_id, message)
			VALUES ($1, $2)`

	result, err := repo.db.Exec(query, message.GetUserId(), message.GetMessage())
	if err == sql.ErrNoRows {
		return psql.ErrNotInserted
	}
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (repo *PsqlNotificationRepository) GetUsersNotifications(ctx context.Context, userID int64) ([]*notificationsPB.NotificationMessage, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("getting user's notifications")

	query := `SELECT
			n.message 
		FROM
			hnh_data.notification n
		WHERE
			n.user_id = $1`

	rows, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	notificationsToReturn := []*notificationsPB.NotificationMessage{}
	for rows.Next() {
		notification := notificationsPB.NotificationMessage{UserId: userID}
		err := rows.Scan(&notification.Message)
		if err != nil {
			return nil, err
		}
		notificationsToReturn = append(notificationsToReturn, &notification)
	}

	if len(notificationsToReturn) == 0 {
		return nil, psql.ErrEntityNotFound
	}

	return notificationsToReturn, nil
}

func (repo *PsqlNotificationRepository) DeleteUsersNotifications(ctx context.Context, userID int64) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("deleting user's notifications")

	query := `DELETE
			FROM
				hnh_data.notification n
			WHERE
				n.user_id = $1`

	result, err := repo.db.Exec(query, userID)
	if err == sql.ErrNoRows {
		return psql.ErrNoRowsDeleted
	}
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
