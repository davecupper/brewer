package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth_MissingConfig(t *testing.T) {
	cmd := Health()
	cmd.SetArgs([]string{"api"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}

func TestHealth_InvalidConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "brewer.yaml")
	err := os.WriteFile(cfgPath, []byte("not: valid: yaml:"), 0644)
	require.NoError(t, err)

	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(old) })

	cmd := Health()
	cmd.SetArgs([]string{"api"})

	err = cmd.Execute()
	require.Error(t, err)
}

func TestHealth_UnknownService(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "brewer.yaml")
	content := `services:
  - name: api
    command: echo api
`
	require.NoError(t, os.WriteFile(cfgPath, []byte(content), 0644))

	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(old) })

	cmd := Health()
	cmd.SetArgs([]string{"unknown"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown service")
}

func TestHealth_NoHealthCheck(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "brewer.yaml")
	content := `services:
  - name: api
    command: echo api
`
	require.NoError(t, os.WriteFile(cfgPath, []byte(content), 0644))

	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(old) })

	var buf bytes.Buffer
	cmd := Health()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"api"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "no health check configured")
}
