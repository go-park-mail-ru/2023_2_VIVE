package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ContextKey string

const (
	REQUEST_ID_KEY = ContextKey("request_id")
)

func GetRequestIDCtx(ctx context.Context) (string, string) {
	return string(REQUEST_ID_KEY), ctx.Value(REQUEST_ID_KEY).(string)
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		ctx := context.WithValue(r.Context(), REQUEST_ID_KEY, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
