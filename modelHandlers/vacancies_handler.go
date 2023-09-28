package modelHandlers

import "models/models"

func GetVacancies() []models.Vacancy {
	models.CreateVacancies()

	defer models.Vac.Mu.Unlock()

	models.Vac.Mu.Lock()
	listToReturn := models.Vac.VacancyList
	models.Vac.Mu.Unlock()

	return listToReturn
}
