// Package main is the entry point for the brewer CLI tool.
// It wires together configuration loading, dependency graph resolution,
// and the service runner to manage local development dependencies.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/brewer/internal/config"
	"github.com/yourorg/brewer/internal/graph"
	"github.com/yourorg/brewer/internal/runner"
	"github.com/yourorg/brewer/internal/ui"
)

const defaultConfigFile = "brewer.yaml"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		configFile = flag.String("config", defaultConfigFile, "path to brewer config file")
		showStatus = flag.Bool("status", false, "print service status table and exit")
	)
	flag.Parse()

	// Load and validate configuration.
	cfg, err := config.Load(*configFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Build the dependency graph and derive a safe startup order.
	g, err := graph.Build(cfg)
	if err != nil {
		return fmt.Errorf("building dependency graph: %w", err)
	}

	order, err := g.TopologicalOrder()
	if err != nil {
		return fmt.Errorf("resolving startup order: %w", err)
	}

	printer := ui.NewPrinter()
	r := runner.New(cfg, printer)

	// If the user only wants a status snapshot, print and exit.
	if *showStatus {
		table := ui.NewStatusTable()
		for _, name := range order {
			snap := r.Snapshot(name)
			table.Add(snap)
		}
		table.Render(os.Stdout)
		return nil
	}

	// Set up a context that is cancelled on SIGINT or SIGTERM so that
	// all services can be gracefully shut down.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "\nshutting down…")
		cancel()
	}()

	// Start services in dependency order.
	for _, name := range order {
		if err := r.Start(ctx, name); err != nil {
			return fmt.Errorf("starting service %q: %w", name, err)
		}
	}

	// Block until the context is cancelled (signal received).
	<-ctx.Done()

	// Stop all services in reverse order.
	for i := len(order) - 1; i >= 0; i-- {
		name := order[i]
		if err := r.Stop(name); err != nil {
			fmt.Fprintf(os.Stderr, "warning: stopping %q: %v\n", name, err)
		}
	}

	return nil
}
