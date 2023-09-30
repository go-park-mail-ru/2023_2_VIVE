package requestHandlers

import (
	"encoding/json"
	"models/modelHandlers"
	"models/serverErrors"
	"net/http"
)

func GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies := modelHandlers.GetVacancies()

	js, err := json.Marshal(vacancies)
	if err != nil {
		errResp := serverErrors.ServerError{Message: err.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJs)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
