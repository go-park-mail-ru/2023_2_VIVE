package server

import (
	"HnH/app"
	"HnH/pkg/middleware"
	interceptors "HnH/pkg/serviceInterceptors"
	"HnH/services/auth/authPB"
	authConfig "HnH/services/auth/config"
	"HnH/services/notifications/config"
	deliveryGrpc "HnH/services/notifications/internal/delivery/grpc"
	"HnH/services/notifications/internal/delivery/websocket"
	repositoryGRPC "HnH/services/notifications/internal/repository/grpc"
	repositoryIM "HnH/services/notifications/internal/repository/inMemory"
	repositoryPSQL "HnH/services/notifications/internal/repository/psql"
	"HnH/services/notifications/internal/usecase"
	"HnH/services/notifications/pkg/logger"
	"HnH/services/notifications/pkg/wsMiddleware"
	"fmt"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		interceptors.RequestIDInterceptor(logger.Logger),
		interceptors.AccesLogInterceptor(logger.Logger, config.NotificationGRPCServiceConfig.ServiceName),
		interceptors.RecoverInterceptor(),
	))

	return opts
}

func initAuthClient(config authConfig.AuthConfig) (authPB.AuthClient, error) {
	connAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(connAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := authPB.NewAuthClient(conn)
	return client, nil
}

func Run() {
	authClient, err := initAuthClient(authConfig.AuthServiceConfig)
	if err != nil {
		fmt.Printf("Error while initializing auth client\n")
		os.Exit(1)
	}

	db, err := app.GetPostgres()
	if err != nil {
		fmt.Printf("Error while initializing psql db\n")
		os.Exit(1)
	}
	notificationRepo := repositoryPSQL.NewPsqlNotificationRepository(db)

	connRepo := repositoryIM.NewInMemoryConnectionRepository()
	authRepo := repositoryGRPC.NewGrpcAuthRepository(authClient)

	notificationUseCase := usecase.NewNotificationUseCase(connRepo, notificationRepo)
	authUseCase := usecase.NewAuthUsecase(authRepo)

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
	go deliveryGrpc.StartGRPCServer(notificationUseCase, listner, opts...)

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

	wsHandler := websocket.NewNotificationWebSocketHandler(notificationUseCase)

	wsHandlerWithMiddleware := http.HandlerFunc(wsHandler.HandleWebSocket)

	wsHandlerWithUserID := wsMiddleware.AuthMiddleware(authUseCase, wsHandlerWithMiddleware)
	wsHandlerWithlogger := wsMiddleware.AccessLogMiddleware(wsHandlerWithUserID)
	wsHandlerWithRequestID := middleware.RequestID(wsHandlerWithlogger)
	firstHandler := middleware.PanicRecoverMiddleware(logger.Logger, wsHandlerWithRequestID)

	http.Handle("/ws", firstHandler)

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
