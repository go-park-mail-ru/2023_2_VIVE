package repository

import (
	"HnH/internal/domain"
	"HnH/internal/repository/mock"
)

type IVacancyRepository interface {
	GetAllVacancies() ([]domain.Vacancy, error)
	GetVacanciesByIds(idList []int) ([]domain.Vacancy, error)
	GetVacancy(vacancyID int) (*domain.Vacancy, error)
	GetOrgId(vacancyID int) (int, error)
	AddVacancy(vacancy *domain.Vacancy) (int, error)
	UpdateVacancy(vacancy *domain.Vacancy) error
	DeleteVacancy(vacancyID int) error
}

type psqlVacancyRepository struct {
	vacancyStorage *mock.Vacancies
}

func NewPsqlVacancyRepository(vacancies *mock.Vacancies) IVacancyRepository {
	return &psqlVacancyRepository{
		vacancyStorage: vacancies,
	}
}

func (p *psqlVacancyRepository) GetAllVacancies() ([]domain.Vacancy, error) {
	p.vacancyStorage.Mu.RLock()

	defer p.vacancyStorage.Mu.RUnlock()

	listToReturn := p.vacancyStorage.VacancyList

	return listToReturn, nil
}

func (p *psqlVacancyRepository) GetVacanciesByIds(idList []int) ([]domain.Vacancy, error) {
	p.vacancyStorage.Mu.RLock()

	defer p.vacancyStorage.Mu.RUnlock()

	listToReturn := make([]domain.Vacancy, 0, len(idList))
	for _, idToFind := range idList {
		for _, vac := range p.vacancyStorage.VacancyList {
			if vac.ID == idToFind {
				listToReturn = append(listToReturn, vac)
				break
			}
		}
	}

	return listToReturn, nil
}

func (p *psqlVacancyRepository) GetVacancy(vacancyID int) (*domain.Vacancy, error) {
	p.vacancyStorage.Mu.RLock()

	defer p.vacancyStorage.Mu.RUnlock()

	indexToReturn := -1
	for index, elem := range p.vacancyStorage.VacancyList {
		if elem.ID == vacancyID {
			indexToReturn = index
			break
		}
	}

	if indexToReturn == -1 {
		return nil, ENTITY_NOT_FOUND
	}

	return &p.vacancyStorage.VacancyList[indexToReturn], nil
}

func (p *psqlVacancyRepository) GetOrgId(vacancyID int) (int, error) {
	p.vacancyStorage.Mu.RLock()

	defer p.vacancyStorage.Mu.RUnlock()

	foundIndex := -1
	for index, elem := range p.vacancyStorage.VacancyList {
		if elem.ID == vacancyID {
			foundIndex = index
			break
		}
	}

	if foundIndex == -1 {
		return 0, ENTITY_NOT_FOUND
	}

	return p.vacancyStorage.VacancyList[foundIndex].CompanyID, nil
}

func (p *psqlVacancyRepository) AddVacancy(vacancy *domain.Vacancy) (int, error) {
	p.vacancyStorage.Mu.Lock()

	defer p.vacancyStorage.Mu.Unlock()

	p.vacancyStorage.CurrentID++
	vacancy.ID = p.vacancyStorage.CurrentID
	p.vacancyStorage.VacancyList = append(p.vacancyStorage.VacancyList, *vacancy)

	return vacancy.ID, nil
}

func (p *psqlVacancyRepository) UpdateVacancy(vacancy *domain.Vacancy) error {
	p.vacancyStorage.Mu.Lock()

	defer p.vacancyStorage.Mu.Unlock()

	indexToUpdate := -1
	for index, elem := range p.vacancyStorage.VacancyList {
		if elem.ID == vacancy.ID {
			indexToUpdate = index
			break
		}
	}

	if indexToUpdate == -1 {
		return ENTITY_NOT_FOUND
	}

	if vacancy.Name != "" {
		p.vacancyStorage.VacancyList[indexToUpdate].Name = vacancy.Name
	}
	if vacancy.CompanyName != "" {
		p.vacancyStorage.VacancyList[indexToUpdate].CompanyName = vacancy.CompanyName
	}
	if vacancy.Description != "" {
		p.vacancyStorage.VacancyList[indexToUpdate].Description = vacancy.Description
	}
	if vacancy.EducationType != "" {
		p.vacancyStorage.VacancyList[indexToUpdate].EducationType = vacancy.EducationType
	}
	if vacancy.Employment != "" {
		p.vacancyStorage.VacancyList[indexToUpdate].Employment = vacancy.Employment
	}
	if vacancy.Experience != "" {
		p.vacancyStorage.VacancyList[indexToUpdate].Experience = vacancy.Experience
	}
	if vacancy.Salary != 0 {
		p.vacancyStorage.VacancyList[indexToUpdate].Salary = vacancy.Salary
	}
	if vacancy.Location != "" {
		p.vacancyStorage.VacancyList[indexToUpdate].Location = vacancy.Location
	}

	return nil
}

func (p *psqlVacancyRepository) DeleteVacancy(vacancyID int) error {
	p.vacancyStorage.Mu.Lock()

	defer p.vacancyStorage.Mu.Unlock()

	indexToDelete := -1
	for index, elem := range p.vacancyStorage.VacancyList {
		if elem.ID == vacancyID {
			indexToDelete = index
			break
		}
	}

	if indexToDelete == -1 {
		return ENTITY_NOT_FOUND
	}

	p.vacancyStorage.VacancyList = append(p.vacancyStorage.VacancyList[:indexToDelete], p.vacancyStorage.VacancyList[indexToDelete+1:]...)

	return nil
}
