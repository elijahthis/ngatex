package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Types
type Config struct {
	Routes   []Route  `yaml:"routes"`
	Services Services `yaml:"services"`
}

type Route struct {
	Path    string `yaml:"path"`
	Service string `yaml:"service"`
}

type Services struct {
	ServiceA Service `yaml:"service-a"`
	ServiceB Service `yaml:"service-b"`
}

type Service struct {
	Upstreams           []string `yaml:"upstreams"`
	LoadBalancingPolicy string   `yaml:"load_balancing_policy"`
}

// Load config file
func Load(path string) (*Config, error) {
	data, fileErr := os.ReadFile(path)
	if fileErr != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, fileErr)
	}

	var yamlData Config

	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		// Return the error
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &yamlData, nil
}
