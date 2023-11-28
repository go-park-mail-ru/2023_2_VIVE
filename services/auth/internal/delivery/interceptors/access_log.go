package interceptors

import (
	"HnH/pkg/contextUtils"
	"HnH/services/auth/config"
	"HnH/services/auth/pkg/logger"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func AccesLogInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	next grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	// requestID := contextUtils.GetRequestIDFromCtx(ctx)
	requestID := contextUtils.GetRequestIDCtx(ctx)

	contextLogger := logger.Logger.WithFields(logrus.Fields{
		"service":    config.AuthServiceConfig.ServiceName,
		"request_id": requestID,
	})

	ctx = context.WithValue(ctx, contextUtils.LOGGER_KEY, contextLogger)

	reply, err := next(ctx, req)

	toLog := contextLogger.WithFields(logrus.Fields{
		"execution_time": time.Since(start).String(),
	})

	if err != nil {
		toLog.Error(err)
	} else {
		toLog.Info("Request successful")
	}
	return reply, err
}
