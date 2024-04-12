package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func RegisterPrometheusMetrics() {
	prometheus.MustRegister(latency)
}

func RecordRequestLatency() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start).Seconds()
		latency.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(elapsed)
	}
}

var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "auth",
		Name:       "latency_seconds",
		Help:       "Latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)
