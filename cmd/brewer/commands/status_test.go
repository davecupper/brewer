package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStatus_MissingConfig(t *testing.T) {
	err := Status("/nonexistent/brewer.yaml")
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestStatus_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brewer.yaml")

	content := `services:
  - name: api
    command: echo hello
  - name: worker
    command: echo world
    depends_on:
      - api
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	// Status should succeed even when no processes are running.
	if err := Status(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatus_InvalidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brewer.yaml")

	content := `services:
  - name: api
    command: echo hello
  - name: api
    command: echo duplicate
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	err := Status(path)
	if err == nil {
		t.Fatal("expected error for duplicate service name, got nil")
	}
}
