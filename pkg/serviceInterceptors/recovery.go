package interceptors

import (
	"HnH/services/searchEngineService/pkg/logger"
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RecoverInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	next grpc.UnaryHandler,
) (reply any, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"panic": err,
			}).
				Error("panic recovered")
		}
	}()

	reply, err = next(ctx, req)
	return reply, err
}
