package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nickcorin/brewer/internal/config"
	"github.com/nickcorin/brewer/internal/process"
)

// Logs returns a cobra command that tails the log output for a named service.
func Logs() *cobra.Command {
	var follow bool
	var lines int

	cmd := &cobra.Command{
		Use:   "logs <service>",
		Short: "Show log output for a service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load("brewer.yaml")
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			serviceName := args[0]
			svc := cfg.ServiceByName(serviceName)
			if svc == nil {
				return fmt.Errorf("unknown service: %q", serviceName)
			}

			reg := process.NewRegistry()
			proc := reg.Get(serviceName)
			if proc == nil {
				return fmt.Errorf("service %q is not running", serviceName)
			}

			snap := proc.Snapshot()
			if len(snap.Logs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no logs available for %q\n", serviceName)
				return nil
			}

			output := snap.Logs
			if !follow && lines > 0 && len(output) > lines {
				output = output[len(output)-lines:]
			}

			for _, line := range output {
				fmt.Fprintln(cmd.OutOrStdout(), line)
			}

			if follow {
				fmt.Fprintln(os.Stderr, "--follow not yet supported in this version")
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output (not yet implemented)")
	cmd.Flags().IntVarP(&lines, "lines", "n", 20, "Number of lines to show from the end")

	return cmd
}
