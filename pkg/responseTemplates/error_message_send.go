package responseTemplates

import (
	"HnH/pkg/serverErrors"
	"fmt"
	"net/http"
)

//easyjson:json
type ErrorToSend struct {
	Message string `json:"message"`
}

func SendErrorMessage(w http.ResponseWriter, err error, statusCode int) error {
	errResp := ErrorToSend{Message: err.Error()}
	mess, marshalErr := errResp.MarshalJSON()
	if marshalErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":%s}`, serverErrors.INTERNAL_SERVER_ERROR.Error())))

		return marshalErr
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(mess)

	return nil
}
