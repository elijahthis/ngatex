package gateway

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/elijahthis/ngatex/pkg/config"
	"github.com/elijahthis/ngatex/pkg/health"
	"github.com/elijahthis/ngatex/pkg/loadbalancer"
	"github.com/elijahthis/ngatex/pkg/middleware"
	"github.com/elijahthis/ngatex/pkg/router"
	"github.com/rs/zerolog/log"
)

type Engine struct {
	CurrentRouter atomic.Pointer[router.Router]
	Config        *config.Config
	MwFactory     map[string]func(http.Handler) http.Handler
	PrevCancel    context.CancelFunc
}

func (e *Engine) ReloadConfig(cfg *config.Config) error {
	if e.PrevCancel != nil {
		e.PrevCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.PrevCancel = cancel

	newRouter := router.New()

	newRouter.Router.Use(middleware.Metrics)
	newRouter.Router.Use(middleware.RequestID)
	newRouter.Router.Use(middleware.Logger)

	newRouter.AddDefaultRoutes()

	routeMap := config.BuildRouteServiceMap(cfg)

	for _, route := range cfg.Routes {
		route := route
		service := routeMap[route.Path]
		var lb loadbalancer.Balancer
		var err error

		switch service.LoadBalancingPolicy {
		case "round-robin":
			lb, err = loadbalancer.NewRoundRobin(service.Upstreams)
		case "weighted-round-robin":
			lb, err = loadbalancer.NewWeightedRoundRobin(service.Upstreams)
		case "least-connections":
			lb, err = loadbalancer.NewLeastConnections(service.Upstreams)
		}

		if err != nil {
			log.Error().
				Err(err).
				Str("route.Path", route.Path).
				Msg("Failed to initialize load balancer for ")
			return err
		}

		health.StartActiveServiceChecks(lb.GetUpstreams(), 10*time.Second, ctx)

		proxyHandler := newRouter.CreateProxyHandler(lb)
		finalHandler := http.StripPrefix(route.Path, proxyHandler)

		var mwStack []func(http.Handler) http.Handler
		for _, mwName := range route.MiddlewareNames {
			if mw, ok := e.MwFactory[mwName]; ok {
				mwStack = append(mwStack, mw)
			} else {
				log.
					Warn().
					Str("mwName", mwName).
					Msgf("Warning: middleware '%s' not found", mwName)
			}
		}

		newRouter.Router.With(mwStack...).Handle(route.Path, finalHandler)
	}

	e.CurrentRouter.Store(newRouter)
	e.Config = cfg
	log.Info().Msg("Reloaded Successfully")

	return nil
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routerPtr := e.CurrentRouter.Load()
	if routerPtr != nil {
		routerPtr.Router.ServeHTTP(w, r)
	} else {
		log.Error().Msg("Router not found in Gateway engine struct")
		http.Error(w, "Gateway initializing...", http.StatusServiceUnavailable)
		return
	}
}
