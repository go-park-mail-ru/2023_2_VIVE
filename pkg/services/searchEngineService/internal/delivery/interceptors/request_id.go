package interceptors

import (
	"HnH/pkg/contextUtils"
	"HnH/pkg/services/searchEngineService/pkg/logger"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestIDInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	next grpc.UnaryHandler,
) (any, error) {
	md, mdExists := metadata.FromIncomingContext(ctx)
	logger.Logger.WithFields(logrus.Fields{
		"meta_data": fmt.Sprintf("%#v", md),
	}).
		Debug("got request with meta data")
	requestID := "-"
	if mdExists && md != nil && len(md[string(contextUtils.REQUEST_ID_KEY)]) > 0 {
		requestID = md[string(contextUtils.REQUEST_ID_KEY)][0]
		logger.Logger.WithFields(logrus.Fields{
			"request_id": requestID,
		}).
			Info("got request with 'request_id'")
	}

	ctx = context.WithValue(ctx, contextUtils.REQUEST_ID_KEY, requestID)

	reply, err := next(ctx, req)
	return reply, err
}
