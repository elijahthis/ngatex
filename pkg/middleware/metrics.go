package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/elijahthis/ngatex/pkg/metrics"
	"github.com/go-chi/chi/v5"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		spy := NewResponseSpy(w)
		next.ServeHTTP(spy, r)

		duration := time.Since(start)

		routePath := chi.RouteContext(r.Context()).RoutePattern()
		if routePath == "" {
			routePath = "unknown"
		}

		metrics.GatewayRequestLatency.WithLabelValues(
			routePath,
			r.Method,
		).Observe(duration.Seconds())

		metrics.GatewayRequestsTotal.WithLabelValues(
			routePath,
			r.Method,
			strconv.Itoa(spy.statusCode),
		).Inc()
	})
}
