package commands

import (
	"fmt"
	"os"

	"github.com/celicoo/brewer/internal/config"
	"github.com/celicoo/brewer/internal/process"
	"github.com/celicoo/brewer/internal/ui"
)

// Status prints the current status of all configured services.
func Status(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	registry := process.NewRegistry()

	table := ui.NewStatusTable(os.Stdout)

	for _, svc := range cfg.Services {
		snap := registry.Snapshot(svc.Name)
		table.Add(snap)
	}

	return table.Render()
}
