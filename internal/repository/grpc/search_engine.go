package grpc

import (
	"HnH/pkg/contextUtils"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type ISearchEngineRepository interface {
	SearchVacancyIDs(ctx context.Context, query string, pageNumber, resultsPerPage int64) ([]int64, int64, error)
	SearchCVsIDs(ctx context.Context, query string, pageNumber, resultsPerPage int64) ([]int64, int64, error)
}

type grpcSearchEngineRepository struct {
	client pb.SearchEngineClient
}

func NewGrpcSearchEngineRepository(client pb.SearchEngineClient) ISearchEngineRepository {
	return &grpcSearchEngineRepository{
		client: client,
	}
}

func (repo *grpcSearchEngineRepository) SearchVacancyIDs(ctx context.Context, query string, pageNumber, resultsPerPage int64) ([]int64, int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	request := pb.SearchRequest{
		Query:          query,
		PageNumber:     pageNumber,
		ResultsPerPage: resultsPerPage,
	}

	contextLogger.Info("sending request to search engine server via grpc")
	contextLogger.WithFields(logrus.Fields{
		"query":            query,
		"page_num":         pageNumber,
		"results_per_page": resultsPerPage,
	}).
		Debug("sending request data")

	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDFromCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)
	searchResponse, err := repo.client.SearchVacancies(ctx, &request)
	if err != nil {
		return nil, 0, err
	}

	// foundVacancyIDs := repo.castVacanciesResponse(searchResponse)
	foundVacancyIDs, count := searchResponse.Ids, searchResponse.Count

	return foundVacancyIDs, count, nil
}

func (repo *grpcSearchEngineRepository) SearchCVsIDs(ctx context.Context, query string, pageNumber, resultsPerPage int64) ([]int64, int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	request := pb.SearchRequest{
		Query:          query,
		PageNumber:     pageNumber,
		ResultsPerPage: resultsPerPage,
	}

	contextLogger.Info("sending request to search engine server via grpc")
	contextLogger.WithFields(logrus.Fields{
		"query":            query,
		"page_num":         pageNumber,
		"results_per_page": resultsPerPage,
	}).
		Debug("sending request data")

	md := metadata.Pairs(string(contextUtils.REQUEST_ID_KEY), contextUtils.GetRequestIDFromCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)
	searchResponse, err := repo.client.SearchCVs(ctx, &request)
	if err != nil {
		return nil, 0, err
	}

	// foundCVIDs := repo.castVacanciesResponse(searchResponse)
	foundCVIDs, count := searchResponse.Ids, searchResponse.Count

	return foundCVIDs, count, nil
}
