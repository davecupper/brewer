// Package config loads and validates brewer.yaml.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// HealthCheck describes how to probe a service for readiness.
type HealthCheck struct {
	// Type is either "http" or "tcp".
	Type string `yaml:"type"`
	// Target is a URL (http) or host:port (tcp).
	Target string `yaml:"target"`
}

// Service represents a single managed process.
type Service struct {
	Name        string      `yaml:"name"`
	Command     string      `yaml:"command"`
	Args        []string    `yaml:"args"`
	Dir         string      `yaml:"dir"`
	Env         []string    `yaml:"env"`
	DependsOn   []string    `yaml:"depends_on"`
	HealthCheck *HealthCheck `yaml:"health_check"`
}

// Config is the top-level brewer.yaml structure.
type Config struct {
	Services []Service `yaml:"services"`
}

// ServiceByName returns the service with the given name, or an error.
func (c *Config) ServiceByName(name string) (Service, error) {
	for _, s := range c.Services {
		if s.Name == name {
			return s, nil
		}
	}
	return Service{}, fmt.Errorf("service %q not found", name)
}

// Load reads and validates the config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	seen := make(map[string]bool)
	for _, svc := range cfg.Services {
		if seen[svc.Name] {
			return fmt.Errorf("duplicate service name: %q", svc.Name)
		}
		seen[svc.Name] = true
	}
	return nil
}
