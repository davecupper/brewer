// Package lifecycle orchestrates the startup and shutdown of services
// in dependency order using the graph and runner packages.
package lifecycle

import (
	"context"
	"fmt"

	"github.com/your-org/brewer/internal/config"
	"github.com/your-org/brewer/internal/graph"
	"github.com/your-org/brewer/internal/runner"
	"github.com/your-org/brewer/internal/ui"
)

// Manager coordinates ordered startup and shutdown of services.
type Manager struct {
	cfg     *config.Config
	runner  *runner.Runner
	printer *ui.Printer
}

// New creates a Manager from the given config.
func New(cfg *config.Config, r *runner.Runner, p *ui.Printer) *Manager {
	return &Manager{cfg: cfg, runner: r, printer: p}
}

// StartAll starts all services in topological (dependency-first) order.
func (m *Manager) StartAll(ctx context.Context) error {
	order, err := graph.Build(m.cfg)
	if err != nil {
		return fmt.Errorf("resolving dependency order: %w", err)
	}

	for _, name := range order {
		svc, ok := m.cfg.ServiceByName(name)
		if !ok {
			return fmt.Errorf("service %q not found in config", name)
		}

		m.printer.ServiceStarting(name)

		if err := m.runner.Start(ctx, svc); err != nil {
			m.printer.ServiceFailed(name, err)
			return fmt.Errorf("starting service %q: %w", name, err)
		}

		m.printer.ServiceRunning(name)
	}

	return nil
}

// StopAll stops all services in reverse topological order.
func (m *Manager) StopAll(ctx context.Context) error {
	order, err := graph.Build(m.cfg)
	if err != nil {
		return fmt.Errorf("resolving dependency order: %w", err)
	}

	var firstErr error
	for i := len(order) - 1; i >= 0; i-- {
		name := order[i]
		if err := m.runner.Stop(ctx, name); err != nil {
			m.printer.ServiceFailed(name, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		m.printer.ServiceStopped(name)
	}

	return firstErr
}
