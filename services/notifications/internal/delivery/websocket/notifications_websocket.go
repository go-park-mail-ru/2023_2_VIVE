package websocket

import (
	"HnH/services/notifications/internal/usecase"
	"fmt"
	"log"
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
	useCase *usecase.INotificationUseCase
}

func NewNotificationWebSocketHandler(useCase *usecase.INotificationUseCase) *NotificationWebSocketHandler {
	return &NotificationWebSocketHandler{
		useCase: useCase,
	}
}

func (h *NotificationWebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error: ", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}
		fmt.Printf("message: %s\n", message)
	}
	// TODO: handle incoming connection
}
