package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/brewer/cmd/brewer/commands"
)

func TestRestart_MissingConfig(t *testing.T) {
	err := commands.Restart("nonexistent.yaml", "web")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "load config")
}

func TestRestart_InvalidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brewer.yaml")
	err := os.WriteFile(path, []byte("services:\n  - name: \"\"\n"), 0644)
	assert.NoError(t, err)

	err = commands.Restart(path, "web")
	assert.Error(t, err)
}

func TestRestart_UnknownService(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brewer.yaml")
	content := `services:
  - name: web
    command: echo hello
`
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)

	err = commands.Restart(path, "unknown-service")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown service")
}

func TestRestart_ValidService(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brewer.yaml")
	content := `services:
  - name: web
    command: echo hello
`
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)

	// Service is not running, restart should still succeed gracefully
	err = commands.Restart(path, "web")
	assert.NoError(t, err)
}
