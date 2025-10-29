package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Types
type Config struct {
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Path     string `yaml:"path"`
	Upstream string `yaml:"upstream"`
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
