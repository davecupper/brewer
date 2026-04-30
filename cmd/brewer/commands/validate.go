package commands

import (
	"fmt"

	"github.com/drew/brewer/internal/config"
	"github.com/drew/brewer/internal/graph"
	"github.com/spf13/cobra"
)

// Validate returns a cobra command that parses and validates brewer.yaml
// without starting any services.
func Validate() *cobra.Command {
	var cfgPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate brewer.yaml without starting services",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			order, err := graph.Build(cfg)
			if err != nil {
				return fmt.Errorf("dependency error: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Config is valid. Startup order:\n")
			for i, svc := range order {
				fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s\n", i+1, svc.Name)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&cfgPath, "config", "c", "brewer.yaml", "Path to config file")
	return cmd
}
