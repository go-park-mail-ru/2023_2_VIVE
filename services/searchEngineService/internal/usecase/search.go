package usecase

import (
	"HnH/internal/repository/psql"
	grpcPsql "HnH/services/searchEngineService/internal/repository/psql"
	"HnH/services/searchEngineService/pkg/searchOptions"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"
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

func (u *SearchUsecase) collectVacFilters(ctx context.Context, searchQuery string) ([]*pb.Filter, error) {
	filters := []*pb.Filter{}

	cityFilterValues, err := u.searchRepo.FilterCitiesVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(searchOptions.City),
		Type:   string(searchOptions.CheckBoxSearch),
		Values: cityFilterValues,
	})

	salaryFilterValues, err := u.searchRepo.FilterSalaryVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(searchOptions.Salary),
		Type:   string(searchOptions.DoubleRange),
		Values: salaryFilterValues,
	})

	experienceFilterValues, err := u.searchRepo.FilterExperienceVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(searchOptions.Experience),
		Type:   string(searchOptions.CheckBox),
		Values: experienceFilterValues,
	})

	employmentFilterValues, err := u.searchRepo.FilterEmploymentVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(searchOptions.Employment),
		Type:   string(searchOptions.CheckBox),
		Values: employmentFilterValues,
	})

	educationTypeFilterValues, err := u.searchRepo.FilterEducationTypeVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(searchOptions.EducationType),
		Type:   string(searchOptions.CheckBox),
		Values: educationTypeFilterValues,
	})

	return filters, nil
}

func (u *SearchUsecase) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	options := request.GetOptions()

	searchQuery, err := searchOptions.GetSearchQuery(options)
	if err != nil {
		return &pb.SearchResponse{}, nil
	}

	vacsIDs, count, err := u.searchRepo.SearchVacanciesIDs(ctx, options)
	if err == psql.ErrEntityNotFound {
		return &pb.SearchResponse{}, nil
	}
	if err != nil {
		return nil, err
	}

	filters, filtersErr := u.collectVacFilters(ctx, searchQuery)
	if filtersErr != nil {
		return nil, filtersErr
	}

	res := pb.SearchResponse{
		Ids:     vacsIDs,
		Count:   count,
		Filters: filters,
	}
	return &res, nil
}

func (u *SearchUsecase) collectCvFilters(ctx context.Context, searchQuery string) ([]*pb.Filter, error) {
	filters := []*pb.Filter{}

	cityFilterValues, err := u.searchRepo.FilterCitiesCVs(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(searchOptions.City),
		Type:   string(searchOptions.CheckBoxSearch),
		Values: cityFilterValues,
	})

	educationTypeFilterValues, err := u.searchRepo.FilterEducationTypeCVs(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(searchOptions.EducationType),
		Type:   string(searchOptions.CheckBox),
		Values: educationTypeFilterValues,
	})

	genderFilterValues, err := u.searchRepo.FilterGenderCVs(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(searchOptions.Gender),
		Type:   string(searchOptions.CheckBox),
		Values: genderFilterValues,
	})

	return filters, nil
}

func (u *SearchUsecase) SearchCVs(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	options := request.GetOptions()

	searchQuery, err := searchOptions.GetSearchQuery(options)
	if err != nil {
		return &pb.SearchResponse{}, nil
	}

	cvsIDs, count, err := u.searchRepo.SearchCVsIDs(ctx, options)
	if err == psql.ErrEntityNotFound {
		return &pb.SearchResponse{}, nil
	}
	if err != nil {
		return nil, err
	}

	filters, filtersErr := u.collectCvFilters(ctx, searchQuery)
	if filtersErr != nil {
		return nil, filtersErr
	}

	res := pb.SearchResponse{
		Ids:     cvsIDs,
		Count:   count,
		Filters: filters,
	}
	return &res, nil
}
