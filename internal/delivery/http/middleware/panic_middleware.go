package middleware

import (
	"HnH/pkg/responseTemplates"
	"HnH/pkg/serverErrors"

	"net/http"

	"github.com/sirupsen/logrus"
)

func PanicRecoverMiddleware(logger *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				responseTemplates.SendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)

				logger.WithFields(logrus.Fields{
					"status":   http.StatusInternalServerError,
					"method":   r.Method,
					"URL":      r.URL.Path,
					"endpoint": r.RemoteAddr,
				}).Panic(err)

				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
