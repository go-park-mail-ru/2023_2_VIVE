package delivery

import (
	"HnH/app"
	"HnH/pkg/contextUtils"
	"HnH/pkg/sanitizer"
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

func (s *SearchEngineServer) ctxLoggerWithMethod(ctx context.Context, methodName string) context.Context {
	contextLogger := contextUtils.GetContextLogger(ctx)
	newContextLogger := contextLogger.WithFields(logrus.Fields{
		"method": methodName,
	})
	return context.WithValue(ctx, contextUtils.LOGGER_KEY, newContextLogger)
}

func (s *SearchEngineServer) sanitizeFilterValues(filterValues ...*pb.FilterValue) []*pb.FilterValue {
	result := make([]*pb.FilterValue, 0, len(filterValues))

	for _, value := range filterValues {
		value.Value = sanitizer.XSS.Sanitize(value.Value)
		result = append(result, value)
	}

	return result
}

func (s *SearchEngineServer) sanitizeFilters(filters ...*pb.Filter) []*pb.Filter {
	result := make([]*pb.Filter, 0, len(filters))

	for _, filter := range filters {
		filter.Name = sanitizer.XSS.Sanitize(filter.Name)
		filter.Type = sanitizer.XSS.Sanitize(filter.Type)
		filter.Values = s.sanitizeFilterValues(filter.Values...)
		result = append(result, filter)
	}

	return result
}

func (s *SearchEngineServer) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	ctx = s.ctxLoggerWithMethod(ctx, "SearchVacancies")

	vacanciesResponse, err := s.searchUscase.SearchVacancies(ctx, request)
	if err != nil {
		return nil, err
	}
	vacanciesResponse.Filters = s.sanitizeFilters(vacanciesResponse.Filters...)
	
	return vacanciesResponse, nil
}

func (s *SearchEngineServer) SearchCVs(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	ctx = s.ctxLoggerWithMethod(ctx, "SearchCVs")

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
