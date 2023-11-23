package usecase

import (
	"HnH/internal/domain"
	"HnH/pkg/services/searchEngineService/internal/repository/psql"
	"HnH/pkg/services/searchEngineService/pkg/searchEngineUtils"
	"HnH/pkg/services/searchEngineService/searchEnginePB"
)

type ISearchVacanciesUsecase interface {
	SearchVacancies(request *searchEnginePB.SearchRequest) ([]domain.ApiVacancy, error)
}

type SearchVacanciesUsecase struct {
	vacancyRepo psql.IVacancyRepository
}

func NewSearchVacanciesUscase(vacancyRepo psql.IVacancyRepository) ISearchVacanciesUsecase {
	return &SearchVacanciesUsecase{
		vacancyRepo: vacancyRepo,
	}
}

func (u *SearchVacanciesUsecase) SearchVacancies(request *searchEnginePB.SearchRequest) ([]domain.ApiVacancy, error) {
	query := request.GetQuery()
	pageNumber := request.GetPageNumber()
	resultsPerPage := request.GetResultsPerPage()

	queryWords := searchEngineUtils.GetQueryWords(query)

	return u.vacancyRepo.SearchVacancies(queryWords, pageNumber, resultsPerPage)
}
