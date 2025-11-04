package middleware

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	cache cache.Cache
	r     rate.Limit
	b     int
}

func NewIPRateLimiter(r float64, b int, ttl time.Duration, cleanupInterval time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		cache: *cache.New(ttl, cleanupInterval),
		r:     rate.Limit(r),
		b:     b,
	}
}

func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	if limiter, found := i.cache.Get(ip); found {
		return limiter.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(i.r, i.b)
	i.cache.Set(ip, limiter, cache.DefaultExpiration)

	return limiter
}

func (i *IPRateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("could not get ip from %s: %v", r.RemoteAddr, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := i.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
