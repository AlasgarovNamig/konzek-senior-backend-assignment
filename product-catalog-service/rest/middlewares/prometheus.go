package middlewares

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

func RequestLatencyTrackingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start).Seconds()
		latency.WithLabelValues(r.Method, r.URL.Path).Observe(elapsed)
	})
}

var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "product",
		Name:       "latency_seconds",
		Help:       "Request latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)

func init() {
	prometheus.MustRegister(latency)
}
