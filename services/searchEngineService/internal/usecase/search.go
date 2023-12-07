package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	grpcPsql "HnH/services/searchEngineService/internal/repository/psql"
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
		Name:   string(domain.CityFilter),
		Type:   string(domain.CheckBoxSearch),
		Values: cityFilterValues,
	})

	salaryFilterValues, err := u.searchRepo.FilterSalaryVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(domain.SalaryFilter),
		Type:   string(domain.DoubleRange),
		Values: salaryFilterValues,
	})

	experienceFilterValues, err := u.searchRepo.FilterExperienceVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(domain.ExperienceFilter),
		Type:   string(domain.Radio),
		Values: experienceFilterValues,
	})

	employmentFilterValues, err := u.searchRepo.FilterEmploymentVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(domain.EmploymentFilter),
		Type:   string(domain.Radio),
		Values: employmentFilterValues,
	})

	educationTypeFilterValues, err := u.searchRepo.FilterEducationTypeVacancies(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(domain.EducationTypeFilter),
		Type:   string(domain.Radio),
		Values: educationTypeFilterValues,
	})

	return filters, nil
}

func (u *SearchUsecase) SearchVacancies(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	searchQuery := request.GetQuery()
	pageNumber := request.GetPageNumber()
	resultsPerPage := request.GetResultsPerPage()

	limit := resultsPerPage
	offset := (pageNumber - 1) * resultsPerPage
	// contextLogger := contextUtils.GetContextLogger(ctx)

	vacsIDs, count, err := u.searchRepo.SearchVacanciesIDs(ctx, searchQuery, limit, offset)
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
		Name:   string(domain.CityFilter),
		Type:   string(domain.CheckBoxSearch),
		Values: cityFilterValues,
	})

	educationTypeFilterValues, err := u.searchRepo.FilterEducationTypeCVs(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(domain.EducationTypeFilter),
		Type:   string(domain.CheckBox),
		Values: educationTypeFilterValues,
	})

	genderFilterValues, err := u.searchRepo.FilterGenderCVs(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	filters = append(filters, &pb.Filter{
		Name:   string(domain.GenderFilter),
		Type:   string(domain.CheckBox),
		Values: genderFilterValues,
	})

	return filters, nil
}

func (u *SearchUsecase) SearchCVs(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	searchQuery := request.GetQuery()
	pageNumber := request.GetPageNumber()
	resultsPerPage := request.GetResultsPerPage()

	limit := resultsPerPage
	offset := (pageNumber - 1) * resultsPerPage
	// contextLogger := contextUtils.GetContextLogger(ctx)

	cvsIDs, count, err := u.searchRepo.SearchCVsIDs(ctx, searchQuery, limit, offset)
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
