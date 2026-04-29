package health

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickward/brewer/internal/config"
)

// Checker waits for a service health check to pass.
type Checker struct {
	svc config.Service
}

// New returns a new Checker for the given service.
func New(svc config.Service) *Checker {
	return &Checker{svc: svc}
}

// Wait blocks until the health check passes, the timeout is exceeded,
// or the context is cancelled. Returns nil if no health check is defined.
func (c *Checker) Wait(ctx context.Context) error {
	hc := c.svc.HealthCheck
	if hc.Type == "" {
		return nil
	}

	timeout := 30 * time.Second
	if hc.Timeout != "" {
		d, err := time.ParseDuration(hc.Timeout)
		if err == nil {
			timeout = d
		}
	}

	interval := 500 * time.Millisecond
	if hc.Interval != "" {
		d, err := time.ParseDuration(hc.Interval)
		if err == nil {
			interval = d
		}
	}

	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("health check timed out after %s", timeout)
		}

		var err error
		switch hc.Type {
		case "http":
			err = checkHTTP(ctx, hc.Target)
		case "tcp":
			err = checkTCP(ctx, hc.Target)
		case "exec":
			err = checkExec(ctx, hc.Target)
		default:
			return fmt.Errorf("unknown health check type: %s", hc.Type)
		}

		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}

// check is an alias used in tests.
func (c *Checker) check(ctx context.Context) error {
	return c.Wait(ctx)
}
