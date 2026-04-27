package commands_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nickbullock/brewer/cmd/brewer/commands"
)

func TestPs_MissingConfig(t *testing.T) {
	var buf bytes.Buffer
	err := commands.Ps("/nonexistent/brewer.yaml", &buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "load config")
}

func TestPs_InvalidConfig(t *testing.T) {
	cfg := writeTempConfig(t, `
services:
  - name: ""
    command: echo hello
`)
	var buf bytes.Buffer
	err := commands.Ps(cfg, &buf)
	assert.Error(t, err)
}

func TestPs_ValidConfig_NoRunningServices(t *testing.T) {
	cfg := writeTempConfig(t, `
services:
  - name: alpha
    command: echo hello
  - name: beta
    command: echo world
    depends_on:
      - alpha
`)
	var buf bytes.Buffer
	err := commands.Ps(cfg, &buf)
	assert.NoError(t, err)

	out := buf.String()
	// Table headers should be present
	assert.Contains(t, out, "NAME")
	assert.Contains(t, out, "STATUS")
	// Both services should appear with stopped/unknown state
	assert.Contains(t, out, "alpha")
	assert.Contains(t, out, "beta")
}
