package responseTemplates

import (
	"encoding/json"
	"net/http"
)

type ErrorToSend struct {
	Message string `json:"message"`
}

func SendErrorMessage(w http.ResponseWriter, err error, statusCode int) error {
	errResp := ErrorToSend{Message: err.Error()}
	errJs, _ := json.Marshal(errResp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(errJs)
	if err != nil {
		return err
	}
	return nil
}
