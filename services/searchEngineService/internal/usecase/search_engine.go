package usecase

import (
	"HnH/services/searchEngineService/internal/repository/psql"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"
)

type ISearchUsecase interface {
	SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error)
}

type SearchUsecase struct {
	searchRepo psql.ISearchRepository
}

func NewSearchUscase(searchRepo psql.ISearchRepository) ISearchUsecase {
	return &SearchUsecase{
		searchRepo: searchRepo,
	}
}

func (u *SearchUsecase) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := request.GetQuery()
	pageNumber := request.GetPageNumber()
	resultsPerPage := request.GetResultsPerPage()

	vacanciesIDs, err := u.searchRepo.SearchVacanciesIDs(ctx, query, pageNumber, resultsPerPage)
	if err != nil {
		return nil, err
	}

	res := pb.SearchResponse{
		Ids: vacanciesIDs,
	}

	return &res, nil



	// return searchEngineUtils.DbVacanciesToGrpc(vacanciesIDs), nil
}
