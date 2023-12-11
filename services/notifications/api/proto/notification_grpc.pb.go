// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.11.2
// source: api/proto/notification.proto

package notificationsPB

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	// rpc AddNotification(NotificationMessage) returns ()
	NotifyUser(ctx context.Context, in *NotificationMessage, opts ...grpc.CallOption) (*empty.Empty, error)
	GetUserNotifications(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*UserNotifications, error)
	DeleteUserNotifications(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*empty.Empty, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) NotifyUser(ctx context.Context, in *NotificationMessage, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/notifications.NotificationService/NotifyUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) GetUserNotifications(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*UserNotifications, error) {
	out := new(UserNotifications)
	err := c.cc.Invoke(ctx, "/notifications.NotificationService/GetUserNotifications", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) DeleteUserNotifications(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/notifications.NotificationService/DeleteUserNotifications", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility
type NotificationServiceServer interface {
	// rpc AddNotification(NotificationMessage) returns ()
	NotifyUser(context.Context, *NotificationMessage) (*empty.Empty, error)
	GetUserNotifications(context.Context, *UserID) (*UserNotifications, error)
	DeleteUserNotifications(context.Context, *UserID) (*empty.Empty, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNotificationServiceServer struct {
}

func (UnimplementedNotificationServiceServer) NotifyUser(context.Context, *NotificationMessage) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyUser not implemented")
}
func (UnimplementedNotificationServiceServer) GetUserNotifications(context.Context, *UserID) (*UserNotifications, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserNotifications not implemented")
}
func (UnimplementedNotificationServiceServer) DeleteUserNotifications(context.Context, *UserID) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserNotifications not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_NotifyUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotificationMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).NotifyUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notifications.NotificationService/NotifyUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).NotifyUser(ctx, req.(*NotificationMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_GetUserNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).GetUserNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notifications.NotificationService/GetUserNotifications",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).GetUserNotifications(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_DeleteUserNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).DeleteUserNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notifications.NotificationService/DeleteUserNotifications",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).DeleteUserNotifications(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notifications.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NotifyUser",
			Handler:    _NotificationService_NotifyUser_Handler,
		},
		{
			MethodName: "GetUserNotifications",
			Handler:    _NotificationService_GetUserNotifications_Handler,
		},
		{
			MethodName: "DeleteUserNotifications",
			Handler:    _NotificationService_DeleteUserNotifications_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/notification.proto",
}