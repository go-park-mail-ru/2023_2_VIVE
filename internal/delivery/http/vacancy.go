package http

import (
	"HnH/internal/usecase"

	"encoding/json"
	"net/http"
)

func GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies, getErr := usecase.GetVacancies()
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
