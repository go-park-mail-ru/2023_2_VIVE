package middleware

import (
	"HnH/internal/appErrors"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/serverErrors"

	"context"
	"errors"
	"net/http"
)

func AuthMiddleware(sessionUsecase usecase.ISessionUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if errors.Is(err, http.ErrNoCookie) {
			responseTemplates.SendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusBadRequest)
			return
		} else if err != nil {
			responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
			return
		}

		ctxWithCookie := context.WithValue(r.Context(), contextUtils.SESSION_ID_KEY, cookie.Value)
		userID, authErr := sessionUsecase.CheckLogin(ctxWithCookie)
		if authErr != nil {
			errToSend, code := appErrors.GetErrAndCodeToSend(authErr)
			responseTemplates.SendErrorMessage(w, errToSend, code)
			return
		}

		ctxWithUID := context.WithValue(ctxWithCookie, contextUtils.USER_ID_KEY, userID)

		next.ServeHTTP(w, r.WithContext(ctxWithUID))
	})
}
