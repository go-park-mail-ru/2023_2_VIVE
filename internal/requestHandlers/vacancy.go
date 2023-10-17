package requestHandlers

import (
	"HnH/internal/modelHandlers"
	"encoding/json"
	"net/http"
)

func GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies := modelHandlers.GetVacancies()

	js, err := json.Marshal(vacancies)
	if err != nil {
		sendErrorMessage(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
