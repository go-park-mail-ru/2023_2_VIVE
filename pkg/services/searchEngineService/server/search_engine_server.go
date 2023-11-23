package main

import (
	"HnH/configs"
	pb "HnH/pkg/services/searchEngineService/searchEnginePB"
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchEngineServer struct {
	pb.UnimplementedSearchEngineServer
}

func (s *SearchEngineServer) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.VacanciesSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchVacancies not implemented")
}
func (s *SearchEngineServer) SearchCVs(ctx context.Context, request *pb.SearchRequest) (*pb.VacanciesSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchCVs not implemented")
}

func newServer() *SearchEngineServer {
	// s := &SearchEngineServer{}
	return &SearchEngineServer{}
}

// func (s *SearchEngineServer) mustEmbedUnimplementedSearchEngineServer() {}

func main() {

	listenAddr := fmt.Sprintf(":%d", configs.HnHSearchEngineConfig.Port)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		os.Exit(1)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSearchEngineServer(grpcServer, newServer())
	fmt.Printf("starting search engine server at %d port\n", configs.HnHSearchEngineConfig.Port)
	grpcServer.Serve(lis)
}
