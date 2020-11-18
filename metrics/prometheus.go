package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "urlshortener"

var (
	// RequestProcessed defines a Prometheus counter for a total of redirect requests processed (by response code)
	RequestProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "requests_processed_total",
			Help:      "Number of total redirect requests processed by response code.",
		},
		[]string{"code"},
	)
)

// RegisterAll registers all the app's Prometheus metrics
func RegisterAll() {
	prometheus.MustRegister(RequestProcessed)
}
