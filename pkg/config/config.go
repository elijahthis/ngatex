package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Routes      []*Route            `yaml:"routes"`
	Services    map[string]*Service `yaml:"services"`
	Middlewares MiddlewareConfig    `yaml:"middlewares"`
}

type MiddlewareConfig struct {
	RateLimit  RateLimitConfig  `yaml:"rate-limit"`
	APIKeyAuth APIKeyAuthConfig `yaml:"api-key-auth"`
	JWTAuth    JWTAuthConfig    `yaml:"jwt-auth"`
	Caching    CachingConfig    `yaml:"caching"`
}

type RateLimitConfig struct {
	RequestsPerSecond float64 `yaml:"requests_per_second"`
	Burst             int     `yaml:"burst"`
}

type APIKeyAuthConfig struct {
	Keys []string `yaml:"keys"`
}

type JWTAuthConfig struct {
	SecretKey string `yaml:"secret_key"`
}

type CachingConfig struct {
	TTL string `yaml:"ttl"`
}

type Route struct {
	Path            string   `yaml:"path"`
	Service         string   `yaml:"service"`
	MiddlewareNames []string `yaml:"middleware_names"`
}

type Service struct {
	Upstreams           []string `yaml:"upstreams"`
	LoadBalancingPolicy string   `yaml:"load_balancing_policy"`
}

func Load(path string) (*Config, error) {
	data, fileErr := os.ReadFile(path)
	if fileErr != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, fileErr)
	}

	var yamlData Config

	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &yamlData, nil
}

func BuildRouteServiceMap(cfg *Config) map[string]*Service {
	routeMap := make(map[string]*Service)

	for _, route := range cfg.Routes {
		service, ok := cfg.Services[route.Service]
		if ok {
			routeMap[route.Path] = service
		}
	}

	return routeMap
}
