package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/pkg/contextUtils"
	grpcPsql "HnH/services/searchEngineService/internal/repository/psql"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"
	"strings"
)

type ISearchUsecase interface {
	SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error)
	SearchCVs(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error)
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

	limit := resultsPerPage
	offset := (pageNumber - 1) * resultsPerPage
	contextLogger := contextUtils.GetContextLogger(ctx)

	if strings.TrimSpace(query) == "" {
		contextLogger.Debug("got empty search query")
		cvsIDs, count, err := u.searchRepo.GetAllVacanciesIDs(ctx, limit, offset)
		if err == psql.ErrEntityNotFound {
			return &pb.SearchResponse{}, nil
		}
		if err != nil {
			return nil, err
		}

		filters := []*pb.Filter{}

		cityFilterValues, err := u.searchRepo.FilterCitiesAllVacancies(ctx)
		if err != nil {
			return nil, err
		}
		filters = append(filters, &pb.Filter{
			Name:   string(domain.CityFilter),
			Type:   string(domain.CheckBoxSearch),
			Values: cityFilterValues,
		})

		salaryFilterValues, err := u.searchRepo.FilterSalaryAllVacancies(ctx)
		if err != nil {
			return nil, err
		}
		filters = append(filters, &pb.Filter{
			Name:   string(domain.SalaryFilter),
			Type:   string(domain.Radio),
			Values: salaryFilterValues,
		})

		experienceFilterValues, err := u.searchRepo.FilterExperienceAllVacancies(ctx)
		if err != nil {
			return nil, err
		}
		filters = append(filters, &pb.Filter{
			Name:   string(domain.ExperienceFilter),
			Type:   string(domain.Radio),
			Values: experienceFilterValues,
		})

		employmentFilterValues, err := u.searchRepo.FilterEmploymentAllVacancies(ctx)
		if err != nil {
			return nil, err
		}
		filters = append(filters, &pb.Filter{
			Name:   string(domain.EmploymentFilter),
			Type:   string(domain.Radio),
			Values: employmentFilterValues,
		})

		educationTypeFilterValues, err := u.searchRepo.FilterEducationTypeAllVacancies(ctx)
		if err != nil {
			return nil, err
		}
		filters = append(filters, &pb.Filter{
			Name:   string(domain.EducationTypeFilter),
			Type:   string(domain.Radio),
			Values: educationTypeFilterValues,
		})

		res := pb.SearchResponse{
			Ids:     cvsIDs,
			Count:   count,
			Filters: filters,
		}
		return &res, nil
	} else {
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
}
