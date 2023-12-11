package interceptors

import (
	"HnH/pkg/contextUtils"
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (reply any, err error) {
		defer func() {
			if err := recover(); err != nil {
				contextLogger := contextUtils.GetContextLogger(ctx)
				contextLogger.WithFields(logrus.Fields{
					"panic": err,
				}).
					Error("panic recovered")
			}
		}()

		reply, err = next(ctx, req)
		return reply, err
	}
}
