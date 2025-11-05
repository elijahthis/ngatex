package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

type responseSpy struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func NewResponseSpy(w http.ResponseWriter) *responseSpy {
	return &responseSpy{
		ResponseWriter: w,
		body:           new(bytes.Buffer),
		statusCode:     http.StatusOK,
	}
}

func (rs *responseSpy) WriteHeader(statusCode int) {
	rs.statusCode = statusCode
	rs.ResponseWriter.WriteHeader(statusCode)
}

func (rs *responseSpy) Write(b []byte) (int, error) {
	rs.body.Write(b)
	return rs.ResponseWriter.Write(b)
}

// ------------------- ------------------- ------------------- ------------------- -------------------

type Cache struct {
	store *cache.Cache
	ttl   time.Duration
}

type cacheEntry struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func NewCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		store: cache.New(defaultTTL, 2*defaultTTL),
		ttl:   defaultTTL,
	}
}

func (c *Cache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		cacheKey := r.Method + r.URL.String()
		if entry, found := c.store.Get(cacheKey); found {
			cached := entry.(cacheEntry)

			for k, v := range cached.Header {
				w.Header()[k] = v
			}
			w.WriteHeader(cached.StatusCode)
			w.Write(cached.Body)
			return
		}

		spy := NewResponseSpy(w)
		next.ServeHTTP(spy, r)

		if spy.statusCode == http.StatusOK {
			entry := cacheEntry{
				StatusCode: spy.statusCode,
				Header:     spy.Header().Clone(),
				Body:       spy.body.Bytes(),
			}
			c.store.Set(cacheKey, entry, c.ttl)
		}

	})
}
