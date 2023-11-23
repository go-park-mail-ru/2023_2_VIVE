package searchEngineUtils

import (
	"HnH/internal/domain"
	pb "HnH/services/searchEngineService/searchEnginePB"
	"time"
)

func OneDbVacanciyToGrpc(vacancy *domain.DbVacancy) *pb.Vacancy {
	res := &pb.Vacancy{
		Id:          int64(vacancy.ID),
		EmployerId:  int64(vacancy.EmployerID),
		VacancyName: vacancy.VacancyName,
		Description: vacancy.Description,
		// SalaryLowerBound: int64(*vacancy.SalaryLowerBound),
		Employment:    string(vacancy.Employment),
		EducationType: string(vacancy.EducationType),
		CreatedAt:     vacancy.CreatedAt.Format(time.DateTime),
	}

	if vacancy.SalaryLowerBound != nil {
		res.SalaryLowerBound = int64(*vacancy.SalaryLowerBound)
	}
	if vacancy.SalaryUpperBound != nil {
		res.SalaryLowerBound = int64(*vacancy.SalaryUpperBound)
	}
	if vacancy.ExperienceLowerBound != nil {
		res.SalaryLowerBound = int64(*vacancy.ExperienceLowerBound)
	}
	if vacancy.ExperienceUpperBound != nil {
		res.SalaryLowerBound = int64(*vacancy.ExperienceUpperBound)
	}

	return res
}

func DbVacanciesToGrpc(vacancies []domain.DbVacancy) *pb.VacanciesSearchResponse {
	res := pb.VacanciesSearchResponse{}

	for _, vacancy := range vacancies {
		pbVacancy := OneDbVacanciyToGrpc(&vacancy)
		res.Vacancies = append(res.Vacancies, pbVacancy)
	}

	return &res
}
