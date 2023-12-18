package responseTemplates

import (
	"net/http"
)

//easyjson:json
type ErrorToSend struct {
	Message string `json:"message"`
}

func SendErrorMessage(w http.ResponseWriter, err error, statusCode int) {
	errResp := ErrorToSend{Message: err.Error()}
	mess, _ := errResp.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(mess)
}
