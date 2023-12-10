package WSMiddleware

import (
	"HnH/pkg/contextUtils"
	"HnH/pkg/responseTemplates"
	"HnH/services/notifications/pkg/logger"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter
	status        int
	body          string
	isHeaderWrote bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.isHeaderWrote {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.isHeaderWrote = true
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.isHeaderWrote && rw.status >= 400 {
		message := responseTemplates.ErrorToSend{}

		err := json.Unmarshal(data, &message)
		if err == nil {
			rw.body = message.Message
		}
	}

	return rw.ResponseWriter.Write(data)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
    h, ok := rw.ResponseWriter.(http.Hijacker)
    if !ok {
        return nil, nil, errors.New("hijack not supported")
    }
    return h.Hijack()
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := newResponseWriterWrapper(w)

		requestID := contextUtils.GetRequestIDFromCtx(r.Context())

		contextLogger := logger.Logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"URL":        r.URL.Path,
			"request_id": requestID,
		})

		ctx := context.WithValue(r.Context(), contextUtils.LOGGER_KEY, contextLogger)

		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		toLog := contextLogger.WithFields(logrus.Fields{
			"status":         wrappedWriter.status,
			"execution_time": time.Since(start).String(),
		})

		if wrappedWriter.body != "" {
			toLog.Error(wrappedWriter.body)
		} else {
			toLog.Info("HTTP Request")
		}
	})
}
