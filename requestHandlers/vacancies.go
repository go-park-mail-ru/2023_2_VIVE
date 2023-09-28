package requestHandlers

import (
	"encoding/json"
	"models/errors"
	"models/modelHandlers"
	"net/http"
)

func GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies := modelHandlers.GetVacancies()

	js, err := json.Marshal(vacancies)
	if err != nil {
		http.Error(w, errors.INTERNAL_SERVER_ERROR.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
