package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = "brewer.yaml"

// Service represents a single managed service dependency.
type Service struct {
	Name      string   `yaml:"name"`
	Command   string   `yaml:"command"`
	Args      []string `yaml:"args,omitempty"`
	DependsOn []string `yaml:"depends_on,omitempty"`
	Dir       string   `yaml:"dir,omitempty"`
	Env       []string `yaml:"env,omitempty"`
}

// Config is the top-level brewer configuration.
type Config struct {
	Version  string    `yaml:"version"`
	Services []Service `yaml:"services"`
}

// Load reads and parses a brewer config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// ServiceByName returns a service by its name, or an error if not found.
func (c *Config) ServiceByName(name string) (*Service, error) {
	for i := range c.Services {
		if c.Services[i].Name == name {
			return &c.Services[i], nil
		}
	}
	return nil, fmt.Errorf("service %q not found", name)
}

func (c *Config) validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("no services defined")
	}
	names := make(map[string]struct{}, len(c.Services))
	for _, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("service is missing a name")
		}
		if svc.Command == "" {
			return fmt.Errorf("service %q is missing a command", svc.Name)
		}
		if _, dup := names[svc.Name]; dup {
			return fmt.Errorf("duplicate service name %q", svc.Name)
		}
		names[svc.Name] = struct{}{}
	}
	return nil
}
