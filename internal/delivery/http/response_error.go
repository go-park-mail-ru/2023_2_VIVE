package http

import (
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
