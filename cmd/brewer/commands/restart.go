package commands

import (
	"fmt"

	"github.com/cnaize/brewer/internal/config"
	"github.com/cnaize/brewer/internal/graph"
	"github.com/cnaize/brewer/internal/lifecycle"
	"github.com/cnaize/brewer/internal/process"
	"github.com/cnaize/brewer/internal/ui"
	"github.com/spf13/cobra"
)

// Restart stops all running services and starts them again in dependency order.
func Restart(registry *process.Registry) *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart all services defined in brewer.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load("brewer.yaml")
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			order, err := graph.Build(cfg)
			if err != nil {
				return fmt.Errorf("resolve dependencies: %w", err)
			}

			printer := ui.NewPrinter()

			mgr := lifecycle.New(cfg, order, registry, printer)

			printer.Info("Stopping services...")
			if err := mgr.StopAll(cmd.Context()); err != nil {
				return fmt.Errorf("stop services: %w", err)
			}

			printer.Info("Starting services...")
			if err := mgr.StartAll(cmd.Context()); err != nil {
				return fmt.Errorf("start services: %w", err)
			}

			return nil
		},
	}
}
