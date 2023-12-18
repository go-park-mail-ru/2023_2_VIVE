package responseTemplates

import (
	"HnH/pkg/serverErrors"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

func MarshalAndSend(w http.ResponseWriter, data any) error {
	js, err := json.Marshal(data)
	if err != nil {
		sendErr := SendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return sendErr
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func SendPDF(w http.ResponseWriter, pdf *gofpdf.Fpdf, fileName string) error {
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fileName))

	err := pdf.Output(w)
	if err != nil {
		sendErr := SendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return sendErr
	}
	return nil
}
