package delivery

import (
	"HnH/app"
	"HnH/pkg/contextUtils"
	"HnH/services/searchEngineService/internal/repository/psql"
	"HnH/services/searchEngineService/internal/usecase"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"

	"github.com/sirupsen/logrus"
)

type SearchEngineServer struct {
	pb.UnimplementedSearchEngineServer
	searchUscase usecase.ISearchUsecase
}

func (s *SearchEngineServer) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	newContextLogger := contextLogger.WithFields(logrus.Fields{
		"method": "SearchVacancies",
	})
	ctx = context.WithValue(ctx, contextUtils.LOGGER_KEY, newContextLogger)

	vacanciesIDs, err := s.searchUscase.SearchVacancies(ctx, request)
	if err != nil {
		return nil, err
	}
	return vacanciesIDs, nil
}

func (s *SearchEngineServer) SearchCVs(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	newContextLogger := contextLogger.WithFields(logrus.Fields{
		"method": "SearchCVs",
	})
	ctx = context.WithValue(ctx, contextUtils.LOGGER_KEY, newContextLogger)

	cvsIDs, err := s.searchUscase.SearchCVs(ctx, request)
	if err != nil {
		return nil, err
	}
	return cvsIDs, nil
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
