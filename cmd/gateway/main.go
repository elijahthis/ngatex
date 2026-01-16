package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/elijahthis/ngatex/pkg/config"
	"github.com/elijahthis/ngatex/pkg/gateway"
	"github.com/elijahthis/ngatex/pkg/middleware"
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

	mainEngine := &gateway.Engine{
		Config:    configData,
		MwFactory: mwFactory,
	}

	if err := mainEngine.ReloadConfig(configData); err != nil {
		log.Error().Err(err).Msg("Error while loading config")
	}

	// Admin Server Setup

	go func() {
		adminRouter := chi.NewRouter()

		adminRouter.Handle("/metrics", promhttp.Handler())

		adminRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		adminRouter.Post("/reload", func(w http.ResponseWriter, r *http.Request) {
			newCfg, err := config.Load(*configPath)
			if err != nil {
				log.Error().Err(err).Msg("Unable to reload config")
				http.Error(w, "Unable to reload config", http.StatusBadRequest)
				return
			}

			if err := mainEngine.ReloadConfig(newCfg); err != nil {
				log.Error().Err(err).Msg("Error while reloading config")
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Reloaded Successfully"))
		})

		log.Info().Msg("Starting ADMIN gateway on :8081")
		if err := http.ListenAndServe(":8081", adminRouter); err != nil {
			log.Error().Err(err).Msg("ADMIN server failed")
		}

	}()

	// SIGHUP Signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	go func() {
		for range sigChan {
			log.Info().Msg("Received SIGHUP signal")

			newCfg, err := config.Load(*configPath)
			if err != nil {
				log.Error().Err(err).Msg("Unable to reload config")
			}

			if err := mainEngine.ReloadConfig(newCfg); err != nil {
				log.Error().Err(err).Msg("Error while reloading config")
			}
		}
	}()

	log.Info().Msgf("Starting PUBLIC gateway on :%s", *port)
	if err := http.ListenAndServe(":"+*port, mainEngine); err != nil {
		log.Error().Err(err).Msg("PUBLIC gateway failed")
	}

}
