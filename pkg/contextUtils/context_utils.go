package contextUtils

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type ContextKey string

const (
	REQUEST_ID_KEY = ContextKey("request_id")
	LOGGER_KEY     = ContextKey("logger")
	SESSION_ID_KEY = ContextKey("session")
	USER_ID_KEY    = ContextKey("user_id")
)

func IsLoggedIn(ctx context.Context) (int, bool) {
	_, ok := ctx.Value(SESSION_ID_KEY).(string)
	if !ok {
		return 0, false
	}

	userID, ok := ctx.Value(USER_ID_KEY).(int)
	if !ok {
		return 0, false
	}

	return userID, true
}

func GetContextLogger(ctx context.Context) *logrus.Entry {
	return ctx.Value(LOGGER_KEY).(*logrus.Entry)
}

func GetRequestIDFromCtx(ctx context.Context) string {
	return ctx.Value(REQUEST_ID_KEY).(string)
}

func GetSessionIDFromCtx(ctx context.Context) string {
	return ctx.Value(SESSION_ID_KEY).(string)
}

func GetUserIDFromCtx(ctx context.Context) int {
	return ctx.Value(USER_ID_KEY).(int)
}

func UpdateCtxLoggerWithMethod(ctx context.Context, methodName string) context.Context {
	contextLogger := GetContextLogger(ctx)
	newContextLogger := contextLogger.WithFields(logrus.Fields{
		"method": methodName,
	})
	return context.WithValue(ctx, LOGGER_KEY, newContextLogger)
}

func PutRequestIDToMetaDataCtx(ctx context.Context) context.Context {
	md := metadata.Pairs(string(REQUEST_ID_KEY), GetRequestIDFromCtx(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
