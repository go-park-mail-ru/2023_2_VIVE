package websocket

import (
	"HnH/pkg/contextUtils"
	"HnH/services/notifications/internal/usecase"
	"HnH/services/notifications/pkg/serviceErrors"
	"net/http"

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

	// cookie, err := r.Cookie("session")
	// if err != nil {
	// 	fmt.Printf("no cookies provided\n")
	// 	http.Error(w, "no cookies", http.StatusForbidden)
	// 	return
	// }

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

	// fmt.Printf("cookie: %v\n", cookie)

	// userIDStr := r.URL.Query().Get(USER_ID_KEY)
	// if strings.TrimSpace(userIDStr) == "" {
	// 	http.Error(w, serviceErrors.ErrInvalidUserID.Error(), http.StatusBadRequest)
	// 	return
	// }

	// if err != nil {
	// 	http.Error(w, serviceErrors.ErrInvalidUserID.Error(), http.StatusBadRequest)
	// 	return
	// }

	userID, ok := (r.Context().Value(contextUtils.USER_ID_KEY)).(int64)
	if !ok {
		contextLogger.WithFields(logrus.Fields{
			"user_id": userID,
		}).
			Error("unable to convert user_id to int64")

		conn.Close()
	}

	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("new ws connection with user_id")

	err = h.useCase.SaveConn(r.Context(), userID, conn)
	if err != nil {
		conn.Close()
		return
	}
}
