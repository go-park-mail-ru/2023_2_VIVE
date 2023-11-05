package middleware

import (
	"HnH/pkg/responseTemplates"

	"mime"
	"net/http"
)

func JSONBodyValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		mt, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			responseTemplates.SendErrorMessage(w, MALFORMED_CONTENT_TYPE_HEADER, http.StatusBadRequest)
			return
		}

		if mt != "application/json" {
			responseTemplates.SendErrorMessage(w, INCORRECT_CONTENT_TYPE_JSON, http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}
