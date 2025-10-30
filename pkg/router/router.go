package router

import (
	"github.com/elijahthis/ngatex/pkg/config"
	"github.com/elijahthis/ngatex/pkg/loadbalancer"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	chi.Router
}

func New() *Router {
	r := &Router{
		Router: chi.NewRouter(),
	}

	// default routes, e.g /health
	return r
}

func (r *Router) AddRoute(path string, service *config.Service, lb loadbalancer.Balancer) {
	// ... logic to create the proxy handler ...
	// r.Router.Handle(path, proxyHandler)
}
