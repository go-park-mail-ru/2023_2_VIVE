package usecase

import (
	"HnH/services/searchEngineService/internal/repository/psql"
	"HnH/services/searchEngineService/pkg/searchEngineUtils"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"
)

type ISearchUsecase interface {
	SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.VacanciesSearchResponse, error)
}

type SearchUsecase struct {
	searchRepo psql.ISearchRepository
}

func NewSearchUscase(searchRepo psql.ISearchRepository) ISearchUsecase {
	return &SearchUsecase{
		searchRepo: searchRepo,
	}
}

func (u *SearchUsecase) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.VacanciesSearchResponse, error) {
	query := request.GetQuery()
	pageNumber := request.GetPageNumber()
	resultsPerPage := request.GetResultsPerPage()

	queryWords := searchEngineUtils.GetQueryWords(query)

	vacancies, err := u.searchRepo.SearchVacancies(ctx, queryWords, pageNumber, resultsPerPage)
	if err != nil {
		return nil, err
	}

	return searchEngineUtils.DbVacanciesToGrpc(vacancies), nil
}
