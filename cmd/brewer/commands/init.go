package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var defaultConfig = `# brewer.yaml — service dependency configuration
services:
  - name: postgres
    command: docker run --rm --name pg -e POSTGRES_PASSWORD=secret -p 5432:5432 postgres:15
    health:
      type: tcp
      target: localhost:5432
      interval: 2s
      timeout: 30s

  - name: redis
    command: docker run --rm --name redis -p 6379:6379 redis:7
    health:
      type: tcp
      target: localhost:6379
      interval: 2s
      timeout: 15s

  - name: api
    command: go run ./cmd/api
    depends_on:
      - postgres
      - redis
    health:
      type: http
      target: http://localhost:8080/healthz
      interval: 3s
      timeout: 60s
`

// Init returns a cobra command that scaffolds a default brewer.yaml.
func Init() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a default brewer.yaml in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			const filename = "brewer.yaml"

			if !force {
				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("%s already exists; use --force to overwrite", filename)
				}
			}

			if err := os.WriteFile(filename, []byte(defaultConfig), 0o644); err != nil {
				return fmt.Errorf("writing %s: %w", filename, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", filename)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing brewer.yaml")
	return cmd
}
