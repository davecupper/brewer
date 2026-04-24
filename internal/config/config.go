package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Service represents a single managed service definition.
type Service struct {
	Name      string            `yaml:"name"`
	Command   string            `yaml:"command"`
	Dir       string            `yaml:"dir"`
	Env       map[string]string `yaml:"env"`
	DependsOn []string          `yaml:"depends_on"`
	ReadyOn   string            `yaml:"ready_on"`
}

// Config is the top-level brewer configuration.
type Config struct {
	Version  string    `yaml:"version"`
	Services []Service `yaml:"services"`
}

// Load reads and parses a brewer YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ServiceByName returns the service with the given name, or an error if not found.
func (c *Config) ServiceByName(name string) (*Service, error) {
	for i := range c.Services {
		if c.Services[i].Name == name {
			return &c.Services[i], nil
		}
	}
	return nil, fmt.Errorf("service %q not found", name)
}

func validate(cfg *Config) error {
	seen := make(map[string]struct{}, len(cfg.Services))
	for _, svc := range cfg.Services {
		if svc.Name == "" {
			return fmt.Errorf("service is missing a name")
		}
		if svc.Command == "" {
			return fmt.Errorf("service %q is missing a command", svc.Name)
		}
		if _, dup := seen[svc.Name]; dup {
			return fmt.Errorf("duplicate service name: %q", svc.Name)
		}
		seen[svc.Name] = struct{}{}
	}
	return nil
}
