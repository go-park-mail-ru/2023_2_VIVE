package middleware

import (
	"HnH/pkg/responseTemplates"

	"encoding/json"
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

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func AccessLogMiddleware(logger *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := newResponseWriterWrapper(w)

		next.ServeHTTP(wrappedWriter, r)

		toLog := logger.WithFields(logrus.Fields{
			"status":         wrappedWriter.status,
			"method":         r.Method,
			"URL":            r.URL.Path,
			"endpoint":       r.RemoteAddr,
			"execution_time": time.Since(start).String(),
		})

		if wrappedWriter.body != "" {
			toLog.Error(wrappedWriter.body)
		} else {
			toLog.Info("HTTP Request")
		}
	})
}
