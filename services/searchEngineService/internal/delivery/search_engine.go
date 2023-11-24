package delivery

import (
	"HnH/app"
	"HnH/services/searchEngineService/internal/repository/psql"
	"HnH/services/searchEngineService/internal/usecase"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchEngineServer struct {
	pb.UnimplementedSearchEngineServer
	searchUscase usecase.ISearchUsecase
}

func (s *SearchEngineServer) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	vacanciesIDs, err := s.searchUscase.SearchVacancies(ctx, request)
	if err != nil {
		return nil, err
	}
	return vacanciesIDs, nil
	// return nil, status.Errorf(codes.Unimplemented, "method SearchVacancies not implemented")
}
func (s *SearchEngineServer) SearchCVs(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchCVs not implemented")
}

func NewServer() (*SearchEngineServer, error) {
	db, err := app.GetPostgres()
	if err != nil {
		return nil, err
	}
	searchRepo := psql.NewPsqlSearchRepository(db)
	searchUsecase := usecase.NewSearchUscase(searchRepo)

	return &SearchEngineServer{
		searchUscase: searchUsecase,
	}, nil
}
