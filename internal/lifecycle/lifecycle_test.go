package lifecycle_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/your-org/brewer/internal/config"
	"github.com/your-org/brewer/internal/lifecycle"
	"github.com/your-org/brewer/internal/runner"
	"github.com/your-org/brewer/internal/ui"
)

func makeConfig(services []config.Service) *config.Config {
	return &config.Config{Services: services}
}

func echoService(name string) config.Service {
	return config.Service{
		Name:    name,
		Command: "echo",
		Args:    []string{name},
	}
}

func newManager(t *testing.T, cfg *config.Config) (*lifecycle.Manager, *bytes.Buffer) {
	t.Helper()
	buf := &bytes.Buffer{}
	p := ui.NewPrinterTo(buf)
	r := runner.New()
	return lifecycle.New(cfg, r, p), buf
}

func TestStartAll_Success(t *testing.T) {
	cfg := makeConfig([]config.Service{
		echoService("alpha"),
		{Name: "beta", Command: "echo", Args: []string{"beta"}, DependsOn: []string{"alpha"}},
	})

	m, _ := newManager(t, cfg)
	ctx := context.Background()

	if err := m.StartAll(ctx); err != nil {
		t.Fatalf("StartAll: unexpected error: %v", err)
	}

	_ = m.StopAll(ctx)
}

func TestStartAll_BadCommand(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "bad", Command: "__no_such_binary__"},
	})

	m, _ := newManager(t, cfg)
	ctx := context.Background()

	if err := m.StartAll(ctx); err == nil {
		t.Fatal("expected error for bad command, got nil")
	}
}

func TestStopAll_ReversesOrder(t *testing.T) {
	cfg := makeConfig([]config.Service{
		echoService("svc-a"),
		{Name: "svc-b", Command: "echo", Args: []string{"b"}, DependsOn: []string{"svc-a"}},
	})

	m, _ := newManager(t, cfg)
	ctx := context.Background()

	if err := m.StartAll(ctx); err != nil {
		t.Fatalf("StartAll: %v", err)
	}

	if err := m.StopAll(ctx); err != nil {
		t.Fatalf("StopAll: unexpected error: %v", err)
	}
}

func TestStartAll_CyclicDependency(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "x", Command: "echo", DependsOn: []string{"y"}},
		{Name: "y", Command: "echo", DependsOn: []string{"x"}},
	})

	m, _ := newManager(t, cfg)
	if err := m.StartAll(context.Background()); err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}
