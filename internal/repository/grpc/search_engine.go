package grpc

import (
	"HnH/internal/domain"
	pb "HnH/pkg/services/searchEngineService/searchEnginePB"
	"context"
)

type ISearchEngineRepository interface {
	SearchVacancies(query string, pageNumber, resultsPerPage int32) ([]domain.ApiVacancy, error)
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
	return []domain.ApiVacancy{}
}

func (repo *grpcSearchEngineRepository) SearchVacancies(query string, pageNumber, resultsPerPage int32) ([]domain.ApiVacancy, error) {
	request := pb.SearchRequest{
		Query: query,
		PageNumber: pageNumber,
		ResultsPerPage: resultsPerPage,
	}

	searchResponce, err := repo.client.SearchVacancies(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	foundVacancies := repo.castVacanciesResponse(searchResponce)

	return foundVacancies, nil
}
