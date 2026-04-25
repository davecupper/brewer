// Package health provides readiness checking for services.
// It supports TCP and HTTP probe strategies with configurable
// retry intervals and timeouts.
package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nickcorin/brewer/internal/config"
)

// Checker polls a service until it becomes ready or the context is cancelled.
type Checker struct {
	service config.Service
	interval time.Duration
}

// New returns a Checker for the given service.
func New(svc config.Service) *Checker {
	return &Checker{
		service:  svc,
		interval: 500 * time.Millisecond,
	}
}

// Wait blocks until the service is healthy or ctx expires.
// It returns nil on success and an error if the context is cancelled
// before the service becomes ready.
func (c *Checker) Wait(ctx context.Context) error {
	probe := c.service.HealthCheck
	if probe == nil {
		// No health check configured — assume immediately ready.
		return nil
	}

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("health check timed out for service %q: %w", c.service.Name, ctx.Err())
		case <-ticker.C:
			if err := check(ctx, probe); err == nil {
				return nil
			}
		}
	}
}

func check(ctx context.Context, probe *config.HealthCheck) error {
	switch probe.Type {
	case "http":
		return checkHTTP(ctx, probe.Target)
	case "tcp":
		return checkTCP(ctx, probe.Target)
	default:
		return fmt.Errorf("unknown health check type: %q", probe.Type)
	}
}
