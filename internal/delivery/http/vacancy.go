package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type VacancyHandler struct {
	vacancyUsecase VacancyUsecase
}

func NewVacancyHandler(router *mux.Router, vacancyUCase VacancyUsecase) {
	handler := &VacancyHandler{
		vacancyUsecase: vacancyUCase,
	}

	router.HandleFunc("/vacancies", handler.GetVacancies).Methods("GET")
}

func (vacancyHandler *VacancyHandler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies, getErr := vacancyHandler.vacancyUsecase.GetVacancies()
	if getErr != nil {
		sendErrorMessage(w, getErr, http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(vacancies)
	if err != nil {
		sendErrorMessage(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
