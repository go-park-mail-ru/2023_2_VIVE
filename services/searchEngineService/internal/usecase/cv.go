package usecase

import (
	"HnH/internal/repository/psql"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"context"
)

func (u *SearchUsecase) SearchCVs(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := request.GetQuery()
	pageNumber := request.GetPageNumber()
	resultsPerPage := request.GetResultsPerPage()
	limit := resultsPerPage
	offset := (pageNumber - 1) * resultsPerPage
	// contextLogger := contextUtils.GetContextLogger(ctx)

	cvsIDs, count, err := u.searchRepo.SearchCVsIDs(ctx, query, limit, offset)
		if err == psql.ErrEntityNotFound {
			return &pb.SearchResponse{}, nil
		}
		if err != nil {
			return nil, err
		}
		res := pb.SearchResponse{
			Ids:   cvsIDs,
			Count: count,
		}
		return &res, nil

	// if strings.TrimSpace(query) == "" {
	// 	contextLogger.Debug("got empty search query")
	// 	cvsIDs, count, err := u.searchRepo.GetAllVacanciesIDs(ctx, limit, offset)
	// 	if err == psql.ErrEntityNotFound {
	// 		return &pb.SearchResponse{}, nil
	// 	}
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	filters := []*pb.Filter{}

	// 	cityFilterValues, err := u.searchRepo.FilterCitiesAll()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	filters = append(filters, &pb.Filter{
	// 		Name:   string(domain.CityFilter),
	// 		Type:   string(domain.CheckBoxSearch),
	// 		Values: cityFilterValues,
	// 	})

	// 	salaryFilterValues, err := u.searchRepo.FilterSalaryAll()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	filters = append(filters, &pb.Filter{
	// 		Name:   string(domain.SalaryFilter),
	// 		Type:   string(domain.Radio),
	// 		Values: salaryFilterValues,
	// 	})
		
	// 	experienceFilterValues, err := u.searchRepo.FilterExperienceAll()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	filters = append(filters, &pb.Filter{
	// 		Name:   string(domain.SalaryFilter),
	// 		Type:   string(domain.Radio),
	// 		Values: experienceFilterValues,
	// 	})

	// 	employmentFilterValues, err := u.searchRepo.FilterEmploymentAll()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	filters = append(filters, &pb.Filter{
	// 		Name:   string(domain.SalaryFilter),
	// 		Type:   string(domain.Radio),
	// 		Values: employmentFilterValues,
	// 	})

	// 	educationTypeFilterValues, err := u.searchRepo.FilterEducationTypeAll()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	filters = append(filters, &pb.Filter{
	// 		Name:   string(domain.SalaryFilter),
	// 		Type:   string(domain.Radio),
	// 		Values: educationTypeFilterValues,
	// 	})

	// 	res := pb.SearchResponse{
	// 		Ids:   cvsIDs,
	// 		Count: count,
	// 		Filters: filters,
	// 	}
	// 	return &res, nil
	// } else {
		
	// }

}
