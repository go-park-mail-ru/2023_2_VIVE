package websocket

import (
	"HnH/services/notifications/internal/model"
	"HnH/services/notifications/internal/usecase"
	"HnH/services/notifications/pkg/serviceErrors"
	"net/http"

	"github.com/gorilla/websocket"
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, serviceErrors.ErrOpenConn.Error(), http.StatusBadRequest)
		return
	}

	handshakeMsg := model.HandshakeMessage{}
	err = conn.ReadJSON(handshakeMsg)
	if err != nil {
		http.Error(w, serviceErrors.ErrHandshakeMsg.Error(), http.StatusBadRequest)
		return
	}

	err = h.useCase.SaveConn(r.Context(), handshakeMsg.UserID, conn)
	if err != nil {
		conn.Close()
		return
	}
	// defer conn.Close()

	// for {
	// 	_, message, err := conn.ReadMessage()
	// 	if err != nil {
	// 		fmt.Printf("error: %v\n", err)
	// 		break
	// 	}
	// 	fmt.Printf("message: %s\n", message)
	// }
	// TODO: handle incoming connection
}
