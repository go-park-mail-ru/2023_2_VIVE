package middleware

import (
	"HnH/internal/usecase"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/serverErrors"

	"errors"
	"net/http"
)

func AuthMiddleware(sessionUsecase usecase.ISessionUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if errors.Is(err, http.ErrNoCookie) {
			responseTemplates.SendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
			return
		} else if err != nil {
			responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
			return
		}

		authErr := sessionUsecase.CheckLogin(cookie.Value)
		if authErr != nil {
			responseTemplates.SendErrorMessage(w, authErr, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
