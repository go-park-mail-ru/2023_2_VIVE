package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var ErrorCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "errors_amount_count",
		Help: "Accumulates outgoing errors",
	},
	[]string{"path", "method", "status"},
)

func ErrorCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		status := rw.statusCode
		ErrorCounter.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(status)).Inc()
	})
}
