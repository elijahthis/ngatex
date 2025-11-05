package router

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/elijahthis/ngatex/pkg/loadbalancer"
	"github.com/elijahthis/ngatex/pkg/transport"
	"github.com/go-chi/chi/v5"
)

type contextKey string

const upstreamKey contextKey = "upstream"
const proxyErrorKey contextKey = "proxy_error"

type Router struct {
	Router    chi.Router
	Transport *http.Transport
}

func New() *Router {
	r := &Router{
		Router: chi.NewRouter(),

		Transport: transport.NewGatewayTransport(transport.TransportConfig{
			MaxIdleConns:        5000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     30 * time.Second,
			DialTimeout:         100 * time.Millisecond,
			TLSHandshakeTimeout: 5 * time.Second,
		}),
	}

	return r
}

// func withUpstream(ctx context.Context, u *loadbalancer.Upstream) context.Context {
// 	return context.WithValue(ctx, upstreamKey, u)
// }
// func getUpstream(ctx context.Context) (*loadbalancer.Upstream, bool) {
// 	u, ok := ctx.Value(upstreamKey).(*loadbalancer.Upstream)
// 	return u, ok
// }

func (r *Router) CreateProxyHandler(lb loadbalancer.Balancer) http.Handler {

	proxy := &httputil.ReverseProxy{
		Transport: r.Transport,

		Director: func(req *http.Request) {
			upstream, err := lb.Next()
			if err != nil {
				ctx := context.WithValue(req.Context(), proxyErrorKey, err)
				*req = *req.WithContext(ctx)

				// Force a transport error by setting an invalid host
				req.URL.Host = ""

				return
			}

			upstream.IncActiveConn()

			ctx := context.WithValue(req.Context(), upstreamKey, upstream)
			*req = *req.WithContext(ctx)

			req.URL.Scheme = upstream.URL.Scheme
			req.URL.Host = upstream.URL.Host
			req.Host = upstream.URL.Host
		},

		ModifyResponse: func(res *http.Response) error {
			if u := res.Request.Context().Value(upstreamKey); u != nil {
				u.(*loadbalancer.Upstream).DecActiveConn()
			}
			return nil
		},

		ErrorHandler: func(w http.ResponseWriter, req *http.Request, err error) {
			log.Printf("proxy error for %s: %v", req.URL, err)

			if val := req.Context().Value(proxyErrorKey); val != nil {
				http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
				return
			}

			if u := req.Context().Value(upstreamKey); u != nil {
				upstream := u.(*loadbalancer.Upstream)
				upstream.DecActiveConn()
				upstream.SetAlive(false)
			}

			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("Bad Gateway: upstream unavailable"))

		},
	}

	return proxy
}

func (r *Router) AddDefaultRoutes() {
	r.Router.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Endpoint is active"))
	})
}
