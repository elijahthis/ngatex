package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/elijahthis/ngatex/pkg/config"
	"github.com/elijahthis/ngatex/pkg/health"
	"github.com/elijahthis/ngatex/pkg/loadbalancer"
	"github.com/elijahthis/ngatex/pkg/middleware"
	"github.com/elijahthis/ngatex/pkg/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Main Server Setup
	port := flag.String("port", "8080", "Port to run gateway")
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Info().Str("configPath", *configPath).Msgf("Loading config from ")

	configData, err := config.Load(*configPath)
	if err != nil {
		log.Fatal().
			Str("configPath", *configPath).
			Err(err).
			Msg("Unable to load config from ")
	}

	r := router.New()

	r.Router.Use(middleware.Metrics)
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
				Msg("Invalid cache TTL")
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
				Str("route.Path", route.Path).
				Msg("Failed to initialize load balancer for ")
		}

		health.StartActiveServiceChecks(lb.GetUpstreams(), 10*time.Second)

		proxyHandler := r.CreateProxyHandler(lb)
		finalHandler := http.StripPrefix(route.Path, proxyHandler)

		var mwStack []func(http.Handler) http.Handler
		for _, mwName := range route.MiddlewareNames {
			if mw, ok := mwFactory[mwName]; ok {
				mwStack = append(mwStack, mw)
			} else {
				log.Info().
					Str("mwName", mwName).
					Msgf("Warning: middleware '%s' not found", mwName)
			}
		}

		r.Router.With(mwStack...).Handle(route.Path, finalHandler)
	}

	// Admin Server Setup

	go func() {
		adminRouter := chi.NewRouter()

		adminRouter.Handle("/metrics", promhttp.Handler())

		adminRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		log.Info().Msg("Starting ADMIN gateway on :8081")
		if err := http.ListenAndServe(":8081", adminRouter); err != nil {
			log.Error().Err(err).Msg("ADMIN server failed")
		}

	}()

	log.Info().Msgf("Starting PUBLIC gateway on :%s", *port)
	if err := http.ListenAndServe(":"+*port, r.Router); err != nil {
		log.Error().Err(err).Msg("PUBLIC gateway failed")
	}

}
