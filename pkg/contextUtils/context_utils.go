package contextUtils

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ContextKey string

const (
	REQUEST_ID_KEY = ContextKey("request_id")
	LOGGER_KEY     = ContextKey("logger")
)

func GetContextLogger(ctx context.Context) *logrus.Entry {
	return ctx.Value(LOGGER_KEY).(*logrus.Entry)
}

func GetRequestIDCtx(ctx context.Context) string {
	return ctx.Value(REQUEST_ID_KEY).(string)
}
