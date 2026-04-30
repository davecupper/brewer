package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/brewer/cmd/brewer/commands"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "brewer.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestUp_MissingConfig(t *testing.T) {
	err := commands.Up("/nonexistent/brewer.yaml")
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestUp_InvalidConfig(t *testing.T) {
	p := writeTempConfig(t, "services:\n  - name: \"\"\n    command: echo\n")
	err := commands.Up(p)
	if err == nil {
		t.Fatal("expected validation error for empty service name, got nil")
	}
}

func TestUp_CyclicDeps(t *testing.T) {
	p := writeTempConfig(t, `services:
  - name: a
    command: echo
    depends_on: [b]
  - name: b
    command: echo
    depends_on: [a]
`)
	err := commands.Up(p)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestUp_UnknownDependency(t *testing.T) {
	p := writeTempConfig(t, `services:
  - name: a
    command: echo
    depends_on: [nonexistent]
`)
	err := commands.Up(p)
	if err == nil {
		t.Fatal("expected error for unknown dependency, got nil")
	}
}
