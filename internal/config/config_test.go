package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/brewer-cli/brewer/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "brewer.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  - name: postgres
    command: docker
    args: ["run", "--rm", "-p", "5432:5432", "postgres:16"]
  - name: api
    command: go
    args: ["run", "./cmd/api"]
    depends_on: [postgres]
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(cfg.Services))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/brewer.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_DuplicateServiceName(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  - name: db
    command: postgres
  - name: db
    command: mysql
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for duplicate service name")
	}
}

func TestServiceByName(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  - name: redis
    command: redis-server
`)
	cfg, _ := config.Load(path)
	svc, err := cfg.ServiceByName("redis")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.Command != "redis-server" {
		t.Errorf("expected command redis-server, got %s", svc.Command)
	}

	_, err = cfg.ServiceByName("unknown")
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
}
