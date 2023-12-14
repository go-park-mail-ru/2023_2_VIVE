package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/usecase"
	"HnH/pkg/middleware"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"
	notificationsPB "HnH/services/notifications/api/proto"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (h *NotificationHandler) sanitizeNotifications(notifications *notificationsPB.UserNotifications) {
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
	vars := mux.Vars(r)

	userID, convErr := strconv.ParseInt(vars["userID"], 10, 64)
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	userNotifications, err := h.notificationUsecase.GetUsersNotifications(r.Context(), userID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}
	h.sanitizeNotifications(userNotifications)

	responseTemplates.MarshalAndSend(w, userNotifications)
}

func (h *NotificationHandler) DeleteUsersNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, convErr := strconv.ParseInt(vars["userID"], 10, 64)
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	err := h.notificationUsecase.DeleteUsersNotifications(r.Context(), userID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}
