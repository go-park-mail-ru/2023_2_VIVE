package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var HitCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_amount_count",
		Help: "Accumulates incoming requests",
	},
	[]string{"path", "method"},
)

func HitCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HitCounter.WithLabelValues(r.URL.Path, r.Method).Inc()

		next.ServeHTTP(w, r)
	})
}
