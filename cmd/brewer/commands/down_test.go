package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDown_MissingConfig(t *testing.T) {
	err := Down("/nonexistent/brewer.yaml")
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestDown_ValidConfig_NoRunningServices(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brewer.yaml")

	content := `services:
  - name: api
    command: echo hello
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	// Down should succeed gracefully when no services are running.
	if err := Down(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDown_InvalidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brewer.yaml")

	content := `services:
  - name: svc
    command: echo one
  - name: svc
    command: echo two
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	err := Down(path)
	if err == nil {
		t.Fatal("expected error for invalid config, got nil")
	}
}
