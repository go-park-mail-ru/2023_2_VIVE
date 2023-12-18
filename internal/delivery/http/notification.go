package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/middleware"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type NotificationHandler struct {
	notificationUsecase usecase.INotificationUsecase
}

func NewNotificationHandler(router *mux.Router, notificationUCase usecase.INotificationUsecase, sessionUCase usecase.ISessionUsecase) {
	handler := &NotificationHandler{
		notificationUsecase: notificationUCase,
	}

	router.Handle("/notifications/{userID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetUsersNotifications))).
		Methods("GET")

	router.Handle("/notifications/{userID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.DeleteUsersNotifications))).
		Methods("DELETE")
}

func (h *NotificationHandler) sanitizeNotifications(notifications *domain.UserNotifications) {
	// result := make([]*notificationsPB.NotificationMessage, 0, len(notifications))
	// result := &notificationsPB.UserNotifications{Notifications: make([]*notificationsPB.NotificationMessage, len(notifications.Notifications))}

	for _, notification := range notifications.Notifications {
		notification.Data = sanitizer.XSS.Sanitize(notification.Data)
		notification.Message = sanitizer.XSS.Sanitize(notification.Message)

		// result.Notifications = append(result.Notifications, notification)
	}

	// return result
}

func (h *NotificationHandler) GetUsersNotifications(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	vars := mux.Vars(r)

	userID, convErr := strconv.ParseInt(vars["userID"], 10, 64)
	if convErr != nil {
		sendErr := responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": convErr,
			}).
				Error("could not send error message")
		}
		return
	}

	userNotifications, err := h.notificationUsecase.GetUsersNotifications(r.Context(), userID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return
	}

	h.sanitizeNotifications(userNotifications)

	marshalErr := responseTemplates.MarshalAndSend(w, userNotifications)
	if marshalErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"error_msg": marshalErr,
			"data":      userNotifications,
		}).
			Error("could not send data")
	}
}

func (h *NotificationHandler) DeleteUsersNotifications(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	vars := mux.Vars(r)

	userID, convErr := strconv.ParseInt(vars["userID"], 10, 64)
	if convErr != nil {
		sendErr := responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": convErr,
			}).
				Error("could not send error message")
		}
		return
	}

	err := h.notificationUsecase.DeleteUsersNotifications(r.Context(), userID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
