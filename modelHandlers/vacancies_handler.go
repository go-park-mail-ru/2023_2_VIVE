package modelHandlers

import "models/models"

func GetVacancies() []models.Vacancy {
	defer models.VacancyDB.Mu.Unlock()

	models.VacancyDB.Mu.Lock()
	listToReturn := models.VacancyDB.VacancyList
	models.VacancyDB.Mu.Unlock()

	return listToReturn
}
