package websocket

import (
	"HnH/pkg/contextUtils"
	"HnH/services/notifications/internal/usecase"
	"HnH/services/notifications/pkg/serviceErrors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	USER_ID_KEY = "user_id"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NotificationWebSocketHandler struct {
	useCase usecase.INotificationUseCase
}

func NewNotificationWebSocketHandler(useCase usecase.INotificationUseCase) *NotificationWebSocketHandler {
	return &NotificationWebSocketHandler{
		useCase: useCase,
	}
}

func (h *NotificationWebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"err": err,
		}).
			Error("unable to connect")
		http.Error(w, serviceErrors.ErrOpenConn.Error(), http.StatusBadRequest)
		return
	}

	contextLogger.WithFields(logrus.Fields{
		"addr": conn.RemoteAddr(),
	}).
		Info("got new websocket connection")

	userIDStr := r.URL.Query().Get(USER_ID_KEY)
	if strings.TrimSpace(userIDStr) == "" {
		http.Error(w, serviceErrors.ErrInvalidUserID.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, serviceErrors.ErrInvalidUserID.Error(), http.StatusBadRequest)
		return
	}

	err = h.useCase.SaveConn(r.Context(), userID, conn)
	if err != nil {
		conn.Close()
		return
	}
}
