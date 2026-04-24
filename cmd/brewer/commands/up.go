// Package commands contains CLI sub-command implementations.
package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/your-org/brewer/internal/config"
	"github.com/your-org/brewer/internal/lifecycle"
	"github.com/your-org/brewer/internal/runner"
	"github.com/your-org/brewer/internal/ui"
)

// Up loads the config, starts all services in order, and blocks until
// SIGINT or SIGTERM is received, then shuts everything down gracefully.
func Up(cfgPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	p := ui.NewPrinter()
	r := runner.New()
	m := lifecycle.New(cfg, r, p)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := m.StartAll(ctx); err != nil {
		return fmt.Errorf("startup failed: %w", err)
	}

	p.Info("all services running — press Ctrl+C to stop")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	p.Info("shutting down...")
	if err := m.StopAll(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	p.Info("all services stopped")
	return nil
}
