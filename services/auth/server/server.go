package server

import (
	pb "HnH/services/auth/authPB"
	"HnH/services/auth/config"
	"HnH/services/auth/internal/delivery"
	"HnH/services/auth/internal/delivery/interceptors"
	"HnH/services/auth/pkg/logger"

	"fmt"
	"net"
	"os"

	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

func initRedis() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:   5,
		MaxActive: 5,

		Wait: true,

		IdleTimeout:     0,
		MaxConnLifetime: 0,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(config.AuthRedisConfig.GetConnectionURL())
			if err != nil {
				return nil, err
			}

			_, err = redis.String(conn.Do("PING"))
			if err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
	}

	return pool
}

func initListen() (net.Listener, error) {
	listenAddr := fmt.Sprintf(":%d", config.AuthServiceConfig.Port)
	listner, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	return listner, nil
}

func initGrpcServer(opts []grpc.ServerOption, conn *redis.Pool) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(opts...)

	server, err := delivery.NewServer(conn)
	if err != nil {
		return nil, err
	}

	pb.RegisterAuthServer(grpcServer, server)

	return grpcServer, nil
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
	logFile, err := os.OpenFile(config.AuthServiceConfig.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("failed to open .log file: %v", err)
		os.Exit(1)
	}
	defer logFile.Close()

	logger.InitLogger(logFile)

	listner, listnerErr := initListen()
	if listnerErr != nil {
		fmt.Printf("failed to listen: %v", listnerErr)
		os.Exit(1)
	}

	opts := initInterceptors()

	redisDB := initRedis()
	defer redisDB.Close()

	grpcServer, err := initGrpcServer(opts, redisDB)
	if err != nil {
		fmt.Printf("error while starting server: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("\tstarting %s server at %d port\n", config.AuthServiceConfig.ServiceName, config.AuthServiceConfig.Port)
	logger.Logger.Infof("starting %s server at %d port", config.AuthServiceConfig.ServiceName, config.AuthServiceConfig.Port)
	grpcServer.Serve(listner)
}
