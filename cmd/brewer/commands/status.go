package commands

import (
	"fmt"
	"os"

	"github.com/celicoo/brewer/internal/config"
	"github.com/celicoo/brewer/internal/process"
	"github.com/celicoo/brewer/internal/ui"
)

// Status prints the current status of all configured services.
// It loads the configuration from configPath, queries the process registry
// for each service's current state, and renders a formatted status table
// to stdout.
func Status(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if len(cfg.Services) == 0 {
		fmt.Fprintln(os.Stdout, "No services configured.")
		return nil
	}

	registry := process.NewRegistry()

	table := ui.NewStatusTable(os.Stdout)

	for _, svc := range cfg.Services {
		snap := registry.Snapshot(svc.Name)
		table.Add(snap)
	}

	return table.Render()
}
