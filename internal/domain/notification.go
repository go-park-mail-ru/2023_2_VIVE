package domain

//easyjson:json
type NotificationMessage struct {
	UserId    int64  `json:"user_id,omitempty"`
	Message   string `json:"message,omitempty"`
	Data      string `json:"data,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

//easyjson:json
type NotificationMessageSlice []NotificationMessage

//easyjson:json
type UserNotifications struct {
	Notifications NotificationMessageSlice `json:"notifications,omitempty"`
}
