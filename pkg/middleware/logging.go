package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		spy := NewResponseSpy(w)

		next.ServeHTTP(spy, r)

		duration := time.Since(start)
		reqId := r.Context().Value(RequestIDKey)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", spy.statusCode).
			Dur("duration_ms", duration).
			Interface("request_id", reqId).
			Msg("incoming request")
	})
}
