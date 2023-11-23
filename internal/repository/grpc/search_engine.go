package grpc

import (
	"HnH/internal/domain"
	"HnH/pkg/contextUtils"
	pb "HnH/pkg/services/searchEngineService/searchEnginePB"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type ISearchEngineRepository interface {
	SearchVacancies(ctx context.Context, query string, pageNumber, resultsPerPage int64) ([]domain.ApiVacancy, error)
}

type grpcSearchEngineRepository struct {
	client pb.SearchEngineClient
}

func NewGrpcSearchEngineRepository(client pb.SearchEngineClient) ISearchEngineRepository {
	return &grpcSearchEngineRepository{
		client: client,
	}
}

func (repo *grpcSearchEngineRepository) castVacanciesResponse(response *pb.VacanciesSearchResponse) []domain.ApiVacancy {
	fmt.Printf("response: %v\n", response)
	return []domain.ApiVacancy{}
}

func (repo *grpcSearchEngineRepository) SearchVacancies(ctx context.Context, query string, pageNumber, resultsPerPage int64) ([]domain.ApiVacancy, error) {
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
	searchResponce, err := repo.client.SearchVacancies(ctx, &request)
	if err != nil {
		return nil, err
	}

	foundVacancies := repo.castVacanciesResponse(searchResponce)

	return foundVacancies, nil
}
