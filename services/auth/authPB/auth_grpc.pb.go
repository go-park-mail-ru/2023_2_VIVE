//
//protoc --go_out=. --go_opt=paths=source_relative \
//--go-grpc_out=. --go-grpc_opt=paths=source_relative \
//authPB/auth.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.11.2
// source: authPB/auth.proto

package authPB

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Auth_AddSession_FullMethodName         = "/auth.auth/AddSession"
	Auth_DeleteSession_FullMethodName      = "/auth.auth/DeleteSession"
	Auth_ValidateSession_FullMethodName    = "/auth.auth/ValidateSession"
	Auth_GetUserIdBySession_FullMethodName = "/auth.auth/GetUserIdBySession"
)

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	AddSession(ctx context.Context, in *AuthData, opts ...grpc.CallOption) (*Empty, error)
	DeleteSession(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Empty, error)
	ValidateSession(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Empty, error)
	GetUserIdBySession(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*UserID, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) AddSession(ctx context.Context, in *AuthData, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Auth_AddSession_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) DeleteSession(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Auth_DeleteSession_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) ValidateSession(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Auth_ValidateSession_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) GetUserIdBySession(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*UserID, error) {
	out := new(UserID)
	err := c.cc.Invoke(ctx, Auth_GetUserIdBySession_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	AddSession(context.Context, *AuthData) (*Empty, error)
	DeleteSession(context.Context, *SessionID) (*Empty, error)
	ValidateSession(context.Context, *SessionID) (*Empty, error)
	GetUserIdBySession(context.Context, *SessionID) (*UserID, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) AddSession(context.Context, *AuthData) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSession not implemented")
}
func (UnimplementedAuthServer) DeleteSession(context.Context, *SessionID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSession not implemented")
}
func (UnimplementedAuthServer) ValidateSession(context.Context, *SessionID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateSession not implemented")
}
func (UnimplementedAuthServer) GetUserIdBySession(context.Context, *SessionID) (*UserID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserIdBySession not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_AddSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).AddSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_AddSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).AddSession(ctx, req.(*AuthData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_DeleteSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).DeleteSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_DeleteSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).DeleteSession(ctx, req.(*SessionID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_ValidateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).ValidateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_ValidateSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).ValidateSession(ctx, req.(*SessionID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_GetUserIdBySession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).GetUserIdBySession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_GetUserIdBySession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).GetUserIdBySession(ctx, req.(*SessionID))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddSession",
			Handler:    _Auth_AddSession_Handler,
		},
		{
			MethodName: "DeleteSession",
			Handler:    _Auth_DeleteSession_Handler,
		},
		{
			MethodName: "ValidateSession",
			Handler:    _Auth_ValidateSession_Handler,
		},
		{
			MethodName: "GetUserIdBySession",
			Handler:    _Auth_GetUserIdBySession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authPB/auth.proto",
}
