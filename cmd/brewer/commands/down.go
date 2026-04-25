package commands

import (
	"fmt"
	"os"

	"github.com/celicoo/brewer/internal/config"
	"github.com/celicoo/brewer/internal/lifecycle"
	"github.com/celicoo/brewer/internal/ui"
)

// Down stops all running services defined in the config file.
func Down(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	printer := ui.NewPrinter(os.Stdout)
	mgr := lifecycle.New(cfg, printer)

	if err := mgr.StopAll(); err != nil {
		return fmt.Errorf("stopping services: %w", err)
	}

	return nil
}
