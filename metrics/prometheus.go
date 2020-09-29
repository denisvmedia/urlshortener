package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const Namespace = "urlshortener"

var (
	RequestProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "requests_processed_total",
			Help:      "Number of total redirect requests processed by response code.",
		},
		[]string{"code"},
	)
)

func RegisterAll() {
	prometheus.MustRegister(RequestProcessed)
}
