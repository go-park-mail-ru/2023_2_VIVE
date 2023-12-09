package server

import (
	interceptors "HnH/pkg/serviceInterceptors"
	"HnH/services/notifications/config"
	deliveryGrpc "HnH/services/notifications/internal/delivery/grpc"
	"HnH/services/notifications/internal/delivery/websocket"
	repository "HnH/services/notifications/internal/repository/inMemory"
	"HnH/services/notifications/internal/usecase"
	"HnH/services/notifications/pkg/logger"
	"fmt"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
)

func initLogger() error {
	logFile, err := os.OpenFile(config.NotificationGRPCServiceConfig.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logger.InitLogger(logFile)
	return nil
}

func initListen() (net.Listener, error) {
	listenAddr := fmt.Sprintf(":%d", config.NotificationGRPCServiceConfig.Port)
	listner, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	return listner, nil
}

func initInterceptors() []grpc.ServerOption {
	var opts []grpc.ServerOption
	opts = append(opts, grpc.ChainUnaryInterceptor(
		interceptors.RecoverInterceptor,
		interceptors.RequestIDInterceptor,
		interceptors.AccesLogInterceptor,
	))

	return opts
}

func Run() {
	repo := repository.NewInMemoryNotificationRepository()
	useCase := usecase.NewNotificationUseCase(&repo)

	loggerErr := initLogger()
	if loggerErr != nil {
		fmt.Printf("Error while initializing logger\n")
		os.Exit(1)
	}

	listner, listnerErr := initListen()
	if listnerErr != nil {
		fmt.Printf("failed to listen: %v", listnerErr)
		os.Exit(1)
	}

	opts := initInterceptors()
	go deliveryGrpc.StartGRPCServer(&useCase, listner, opts...)

	fmt.Printf(
		"\tstarting %s grpc server at %d port\n",
		config.NotificationGRPCServiceConfig.ServiceName,
		config.NotificationGRPCServiceConfig.Port,
	)
	logger.Logger.Infof(
		"starting %s grpc server at %d port",
		config.NotificationGRPCServiceConfig.ServiceName,
		config.NotificationGRPCServiceConfig.Port,
	)

	wsHandler := websocket.NewNotificationWebSocketHandler(&useCase)
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)
	fmt.Printf(
		"\tstarting %s websocket server at %d port\n",
		config.NotificationGRPCServiceConfig.ServiceName,
		config.NotificationWSServiceConfig.Port,
	)
	logger.Logger.Infof(
		"starting %s websocket server at %d port",
		config.NotificationGRPCServiceConfig.ServiceName,
		config.NotificationWSServiceConfig.Port,
	)
	http.ListenAndServe(fmt.Sprintf("%s:%d", config.NotificationWSServiceConfig.Host, config.NotificationWSServiceConfig.Port), nil)
}
