package responseTemplates

import (
	"HnH/pkg/serverErrors"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

func MarshalAndSend(w http.ResponseWriter, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		SendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func SendPDF(w http.ResponseWriter, pdf *gofpdf.Fpdf, fileName string) {
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fileName))

	err := pdf.Output(w)
	if err != nil {
		SendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
}
