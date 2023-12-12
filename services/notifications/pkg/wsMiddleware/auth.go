package wsMiddleware

import (
	"HnH/pkg/contextUtils"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/serverErrors"
	"HnH/services/notifications/internal/usecase"
	"context"
	"errors"
	"net/http"
)

func AuthMiddleware(authUsecase usecase.IAuthUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if errors.Is(err, http.ErrNoCookie) {
			responseTemplates.SendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusBadRequest)
			return
		} else if err != nil {
			responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
			return
		}

		userID, userErr := authUsecase.ValidateAndGetUserID(r.Context(), cookie.Value)
		if userErr != nil {
			responseTemplates.SendErrorMessage(w, serverErrors.INCORRECT_CREDENTIALS, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextUtils.USER_ID_KEY, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
