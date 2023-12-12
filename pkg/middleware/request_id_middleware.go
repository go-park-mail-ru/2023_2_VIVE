package middleware

import (
	"HnH/pkg/contextUtils"
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		ctx := context.WithValue(r.Context(), contextUtils.REQUEST_ID_KEY, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
