package graph_test

import (
	"testing"

	"github.com/nickcorin/brewer/internal/config"
	"github.com/nickcorin/brewer/internal/graph"
)

func makeConfig(services []config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestBuild_Valid(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "db", Command: "postgres"},
		{Name: "api", Command: "./api", DependsOn: []string{"db"}},
	})

	g, err := graph.Build(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if g == nil {
		t.Fatal("expected graph, got nil")
	}
}

func TestBuild_UnknownDependency(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "api", Command: "./api", DependsOn: []string{"db"}},
	})

	_, err := graph.Build(cfg)
	if err == nil {
		t.Fatal("expected error for unknown dependency, got nil")
	}
}

func TestBuild_CycleDetected(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "a", Command: "a", DependsOn: []string{"b"}},
		{Name: "b", Command: "b", DependsOn: []string{"a"}},
	})

	_, err := graph.Build(cfg)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestTopologicalOrder_DepsFirst(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "db", Command: "postgres"},
		{Name: "cache", Command: "redis"},
		{Name: "api", Command: "./api", DependsOn: []string{"db", "cache"}},
	})

	g, err := graph.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	order, err := g.TopologicalOrder()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order) != 3 {
		t.Fatalf("expected 3 services, got %d", len(order))
	}

	pos := make(map[string]int)
	for i, svc := range order {
		pos[svc.Name] = i
	}

	if pos["api"] <= pos["db"] {
		t.Errorf("expected db before api")
	}
	if pos["api"] <= pos["cache"] {
		t.Errorf("expected cache before api")
	}
}

func TestTopologicalOrder_NoDeps(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "worker", Command: "./worker"},
	})

	g, err := graph.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	order, _ := g.TopologicalOrder()
	if len(order) != 1 || order[0].Name != "worker" {
		t.Errorf("expected single service 'worker'")
	}
}
