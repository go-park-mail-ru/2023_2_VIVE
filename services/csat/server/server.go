package server

import (
	"HnH/services/csat/config"
	pb "HnH/services/csat/csatPB"
	"HnH/services/csat/internal/delivery"
	"HnH/services/csat/internal/delivery/interceptors"
	"HnH/services/csat/pkg/logger"
	"database/sql"
	"fmt"
	"net"
	"os"

	_ "github.com/jackc/pgx/stdlib"
	"google.golang.org/grpc"
)

func initLogger() error {
	logFile, err := os.OpenFile(config.LOGS_DIR+config.CsatServiceConfig.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logger.InitLogger(logFile)
	return nil
}

func initListen() (net.Listener, error) {
	listenAddr := fmt.Sprintf(":%d", config.CsatServiceConfig.Port)
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

func initPostgres() (*sql.DB, error) {
	dsn := config.CsatPostgresConfig.GetConnectionString()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func initGrpcServer(opts []grpc.ServerOption) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(opts...)

	db, err := initPostgres()
	if err != nil {
		return nil, err
	}

	server, err := delivery.NewServer(db)
	if err != nil {
		return nil, err
	}
	pb.RegisterCsatServer(grpcServer, server)

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
	fmt.Printf("\tstarting %s server at %d port\n", config.CsatServiceConfig.ServiceName, config.CsatServiceConfig.Port)
	logger.Logger.Infof("starting %s server at %d port", config.CsatServiceConfig.ServiceName, config.CsatServiceConfig.Port)
	grpcServer.Serve(listner)
}
