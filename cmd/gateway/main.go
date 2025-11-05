package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/elijahthis/ngatex/pkg/config"
	"github.com/elijahthis/ngatex/pkg/health"
	"github.com/elijahthis/ngatex/pkg/loadbalancer"
	"github.com/elijahthis/ngatex/pkg/middleware"
	"github.com/elijahthis/ngatex/pkg/router"
)

func main() {
	port := flag.String("port", "8080", "Port to run gateway")
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	log.Info().Msgf("Loading config from %s", *configPath)

	configData, err := config.Load(*configPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Unable to load config from %s", *configPath)
	}

	r := router.New()
	r.Router.Use(middleware.RequestID)
	r.Router.Use(middleware.Logger)

	r.AddDefaultRoutes()

	routeMap := config.BuildRouteServiceMap(configData)

	mwFactory := make(map[string]func(http.Handler) http.Handler)

	if cfg := configData.Middlewares.RateLimit; cfg.RequestsPerSecond > 0 {
		limiter := middleware.NewIPRateLimiter(cfg.RequestsPerSecond, cfg.Burst, 5*time.Minute, 10*time.Minute)
		mwFactory["rate-limit"] = limiter.RateLimit
	}

	if cfg := configData.Middlewares.APIKeyAuth; len(cfg.Keys) > 0 {
		auth := middleware.NewAPIKeyAuth(cfg.Keys)
		mwFactory["api-key-auth"] = auth.Auth
	}

	if cfg := configData.Middlewares.JWTAuth; cfg.SecretKey != "" {
		auth := middleware.NewJWTAuth(cfg.SecretKey)
		mwFactory["jwt-auth"] = auth.Auth
	}
	if cfg := configData.Middlewares.Caching; cfg.TTL != "" {
		ttl, err := time.ParseDuration(cfg.TTL)
		if err != nil {
			log.Fatal().
				Err(err).
				Msgf("Invalid cache TTL")
		}
		c := middleware.NewCache(ttl)
		mwFactory["caching"] = c.Middleware
	}

	for _, route := range configData.Routes {
		route := route
		service := routeMap[route.Path]
		var lb loadbalancer.Balancer

		switch service.LoadBalancingPolicy {
		case "round-robin":
			lb, err = loadbalancer.NewRoundRobin(service.Upstreams)
		case "weighted-round-robin":
			lb, err = loadbalancer.NewWeightedRoundRobin(service.Upstreams)
		case "least-connections":
			lb, err = loadbalancer.NewLeastConnections(service.Upstreams)
		}

		if err != nil {
			log.Fatal().
				Err(err).
				Msgf("Failed to initialize load balancer for %s", route.Path)
		}

		health.StartActiveServiceChecks(lb.GetUpstreams(), 10*time.Second)

		proxyHandler := r.CreateProxyHandler(lb)
		finalHandler := http.StripPrefix(route.Path, proxyHandler)

		var mwStack []func(http.Handler) http.Handler
		for _, mwName := range route.MiddlewareNames {
			if mw, ok := mwFactory[mwName]; ok {
				mwStack = append(mwStack, mw)
			} else {
				log.Info().Msgf("Warning: middleware '%s' not found", mwName)
			}
		}

		r.Router.With(mwStack...).Handle(route.Path, finalHandler)
	}

	log.Info().Msgf("Starting PUBLIC gateway on :%s", *port)
	http.ListenAndServe(":"+*port, r.Router)

}
