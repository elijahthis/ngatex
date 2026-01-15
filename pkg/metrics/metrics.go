package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GatewayRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_total",
			Help: "Total number of requests handled by the gateway",
		},
		[]string{"route_path", "method", "status_code"},
	)

	GatewayRequestLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_request_latency_seconds",
			Help:    "Histogram of request latencies",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route_path", "method"},
	)
)
