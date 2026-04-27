package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nickcorin/brewer/internal/config"
	"github.com/nickcorin/brewer/internal/process"
	"github.com/nickcorin/brewer/internal/ui"
)

// Ps returns a cobra command that lists all services defined in the config
// along with their current process state (PID, status, uptime).
func Ps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: "List services and their current process state",
		Long: `Display a detailed process table for all services defined in brewer.yaml.

Shows each service's PID, current status, and how long it has been running.
Services that are not currently managed by brewer will appear as stopped.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			registry := process.NewRegistry()
			table := ui.NewStatusTable(os.Stdout)

			for _, svc := range cfg.Services {
				proc, ok := registry.Get(svc.Name)
				if !ok {
					// Service is known but not tracked — render as stopped.
					table.Render(svc.Name, process.Snapshot{
						Name:   svc.Name,
						Status: process.StatusStopped,
					})
					continue
				}

				table.Render(svc.Name, proc.Snapshot())
			}

			return nil
		},
	}

	return cmd
}
