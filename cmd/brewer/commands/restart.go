package commands

import (
	"fmt"

	"github.com/brewer/internal/config"
	"github.com/brewer/internal/lifecycle"
	"github.com/brewer/internal/process"
	"github.com/brewer/internal/ui"
)

// Restart stops and restarts a named service defined in the config file.
func Restart(configPath, serviceName string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	svc := cfg.ServiceByName(serviceName)
	if svc == nil {
		return fmt.Errorf("unknown service: %q", serviceName)
	}

	printer := ui.NewPrinter()
	reg := process.NewRegistry()
	mgr := lifecycle.New(cfg, reg, printer)

	printer.ServiceStarting(serviceName)

	// Stop the service if it is currently running.
	if proc, ok := reg.Get(serviceName); ok && proc.IsRunning() {
		printer.ServiceStopped(serviceName)
		if stopErr := proc.Stop(); stopErr != nil {
			return fmt.Errorf("stop service %q: %w", serviceName, stopErr)
		}
	}

	// Start only the requested service (not the full dependency graph).
	if startErr := mgr.StartOne(svc); startErr != nil {
		printer.ServiceFailed(serviceName, startErr)
		return fmt.Errorf("start service %q: %w", serviceName, startErr)
	}

	printer.ServiceRunning(serviceName)
	return nil
}
