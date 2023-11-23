// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.11.2
// source: searchEnginePB/search_engine.proto

package searchEnginePB

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

// SearchEngineClient is the client API for SearchEngine service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchEngineClient interface {
	SearchVacancies(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*VacanciesSearchResponse, error)
}

type searchEngineClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchEngineClient(cc grpc.ClientConnInterface) SearchEngineClient {
	return &searchEngineClient{cc}
}

func (c *searchEngineClient) SearchVacancies(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*VacanciesSearchResponse, error) {
	out := new(VacanciesSearchResponse)
	err := c.cc.Invoke(ctx, "/searchEngine.searchEngine/SearchVacancies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchEngineServer is the server API for SearchEngine service.
// All implementations must embed UnimplementedSearchEngineServer
// for forward compatibility
type SearchEngineServer interface {
	SearchVacancies(context.Context, *SearchRequest) (*VacanciesSearchResponse, error)
	mustEmbedUnimplementedSearchEngineServer()
}

// UnimplementedSearchEngineServer must be embedded to have forward compatible implementations.
type UnimplementedSearchEngineServer struct {
}

func (UnimplementedSearchEngineServer) SearchVacancies(context.Context, *SearchRequest) (*VacanciesSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchVacancies not implemented")
}
func (UnimplementedSearchEngineServer) mustEmbedUnimplementedSearchEngineServer() {}

// UnsafeSearchEngineServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchEngineServer will
// result in compilation errors.
type UnsafeSearchEngineServer interface {
	mustEmbedUnimplementedSearchEngineServer()
}

func RegisterSearchEngineServer(s grpc.ServiceRegistrar, srv SearchEngineServer) {
	s.RegisterService(&SearchEngine_ServiceDesc, srv)
}

func _SearchEngine_SearchVacancies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchEngineServer).SearchVacancies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/searchEngine.searchEngine/SearchVacancies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchEngineServer).SearchVacancies(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SearchEngine_ServiceDesc is the grpc.ServiceDesc for SearchEngine service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchEngine_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "searchEngine.searchEngine",
	HandlerType: (*SearchEngineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchVacancies",
			Handler:    _SearchEngine_SearchVacancies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "searchEnginePB/search_engine.proto",
}
