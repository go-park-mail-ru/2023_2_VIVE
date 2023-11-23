package server

import (
	"HnH/pkg/services/searchEngineService/config"
	"HnH/pkg/services/searchEngineService/internal/delivery"
	"HnH/pkg/services/searchEngineService/internal/delivery/interceptors"
	"HnH/pkg/services/searchEngineService/pkg/logger"
	pb "HnH/pkg/services/searchEngineService/searchEnginePB"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
)

func initLogger() error {
	logFile, err := os.OpenFile(config.SearchEngineServiceConfig.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logger.InitLogger(logFile)
	return nil
}

func initListen() (net.Listener, error) {
	listenAddr := fmt.Sprintf(":%d", config.SearchEngineServiceConfig.Port)
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

func initGrpcServer(opts []grpc.ServerOption) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(opts...)
	server, err := delivery.NewServer()
	if err != nil {
		return nil, err
	}
	pb.RegisterSearchEngineServer(grpcServer, server)

	return grpcServer, nil
}

func Run() {
	
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

	grpcServer, err := initGrpcServer(opts)
	if err != nil {
		fmt.Printf("error while starting server: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("\tstarting search engine server at %d port\n", config.SearchEngineServiceConfig.Port)
	logger.Logger.Infof("starting search engine server at %d port", config.SearchEngineServiceConfig.Port)
	grpcServer.Serve(listner)
}
