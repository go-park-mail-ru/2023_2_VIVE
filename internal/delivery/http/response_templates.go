package http

import (
	"HnH/pkg/serverErrors"

	"encoding/json"
	"net/http"
)

type errorToSend struct {
	Message string `json:"message"`
}

func sendErrorMessage(w http.ResponseWriter, err error, statusCode int) {
	errResp := errorToSend{Message: err.Error()}
	errJs, _ := json.Marshal(errResp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(errJs)
}

func marshalAndSend(w http.ResponseWriter, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		sendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
