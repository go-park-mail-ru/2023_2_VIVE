package usecase

import (
	"HnH/internal/repository/psql"
	grpcPsql "HnH/services/searchEngineService/internal/repository/psql"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"
)

type ISearchUsecase interface {
	SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error)
}

type SearchUsecase struct {
	searchRepo grpcPsql.ISearchRepository
}

func NewSearchUscase(searchRepo grpcPsql.ISearchRepository) ISearchUsecase {
	return &SearchUsecase{
		searchRepo: searchRepo,
	}
}

func (u *SearchUsecase) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := request.GetQuery()
	pageNumber := request.GetPageNumber()
	resultsPerPage := request.GetResultsPerPage()

	vacanciesIDs, count, err := u.searchRepo.SearchVacanciesIDs(ctx, query, pageNumber, resultsPerPage)
	if err == psql.ErrEntityNotFound {
		return &pb.SearchResponse{}, nil
	}
	if err != nil {
		return nil, err
	}

	res := pb.SearchResponse{
		Ids:   vacanciesIDs,
		Count: count,
	}

	return &res, nil
}
