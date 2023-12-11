package grpc

import (
	"HnH/pkg/contextUtils"
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
	useCase usecase.INotificationUseCase
}

func NewNotificationServer(useCase usecase.INotificationUseCase) *NotificationServer {
	return &NotificationServer{
		useCase: useCase,
	}
}

func (s *NotificationServer) NotifyUser(ctx context.Context, message *pb.NotificationMessage) (*empty.Empty, error) {
	ctx = contextUtils.UpdateCtxLoggerWithMethod(ctx, "NotifyUser")
	return &empty.Empty{}, s.useCase.SendNotification(ctx, message)
}

func (s *NotificationServer) GetUserNotifications(ctx context.Context, userID *pb.UserID) (*pb.UserNotifications, error) {
	ctx = contextUtils.UpdateCtxLoggerWithMethod(ctx, "GetUserNotifications")
	notifications, err := s.useCase.GetUsersNotifications(ctx, userID.UserId)
	return &pb.UserNotifications{Notifications: notifications}, err
}

func (s *NotificationServer) DeleteUserNotifications(ctx context.Context, userID *pb.UserID) (*empty.Empty, error) {
	ctx = contextUtils.UpdateCtxLoggerWithMethod(ctx, "DeleteUserNotifications")
	return &empty.Empty{}, s.useCase.DeleteUsersNotifications(ctx, userID.UserId)
}

func StartGRPCServer(
	useCase usecase.INotificationUseCase,
	lis net.Listener,
	opts ...grpc.ServerOption,
) {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterNotificationServiceServer(grpcServer, NewNotificationServer(useCase))

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
