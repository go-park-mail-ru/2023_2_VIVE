package grpc

import (
	pb "HnH/services/notifications/api/proto"
	"HnH/services/notifications/internal/usecase"
	"context"
	"fmt"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
	useCase *usecase.INotificationUseCase
}

func NewNotificationServer(useCase *usecase.INotificationUseCase) *NotificationServer {
	return &NotificationServer{
		useCase: useCase,
	}
}

func (s *NotificationServer) NotifyUser(ctx context.Context, message *pb.NotificationMessage) (*empty.Empty, error) {
	// TODO: send message to users
	return nil, nil
}

func StartGRPCServer(
	// ctx context.Context,
	useCase *usecase.INotificationUseCase,
	lis net.Listener,
	opts ...grpc.ServerOption,
) {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterNotificationServiceServer(grpcServer, NewNotificationServer(useCase))

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
