package requestHandlers

import (
	"encoding/json"
	"models/modelHandlers"
	"net/http"
)

func GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies := modelHandlers.GetVacancies()

	js, err := json.Marshal(vacancies)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
