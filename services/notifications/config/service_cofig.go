package config

type NotificationsGRPCConfig struct {
	ServiceName string
	Host        string
	Port        int
	LogFile     string
}

var NotificationGRPCServiceConfig = NotificationsGRPCConfig{
	ServiceName: "Notifications",
	Host:        "localhost",
	Port:        8064,
	LogFile:     "notification_service.log",
}

type NotificationsWSConfig struct {
	Host    string
	Port    int
	LogFile string
}

var NotificationWSServiceConfig = NotificationsWSConfig{
	Host:    "84.23.53.171",
	Port:    8065,
}
