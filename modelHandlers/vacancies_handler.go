package modelHandlers

import "HnH/models"

func GetVacancies() []models.Vacancy {
	models.VacancyDB.Mu.RLock()

	defer models.VacancyDB.Mu.RUnlock()

	listToReturn := models.VacancyDB.VacancyList

	return listToReturn
}