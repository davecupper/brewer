package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/nickcorin/brewer/internal/config"
	"github.com/nickcorin/brewer/internal/health"
)

// Health returns a cobra command that checks the health of a named service.
func Health() *cobra.Command {
	var timeout time.Duration

	cmd := &cobra.Command{
		Use:   "health <service>",
		Short: "Check the health of a service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load("brewer.yaml")
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			name := args[0]
			svc := cfg.ServiceByName(name)
			if svc == nil {
				return fmt.Errorf("unknown service: %q", name)
			}

			if svc.HealthCheck == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "service %q has no health check configured\n", name)
				return nil
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			checker := health.New(svc)
			if err := checker.Wait(ctx); err != nil {
				return fmt.Errorf("health check failed for %q: %w", name, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "service %q is healthy\n", name)
			return nil
		},
	}

	cmd.Flags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "maximum time to wait for health check")
	return cmd
}
