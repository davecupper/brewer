package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInit_CreatesFile(t *testing.T) {
	tmp := t.TempDir()
	prev, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer os.Chdir(prev)

	cmd := Init()
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Created brewer.yaml") {
		t.Errorf("expected creation message, got: %q", out.String())
	}

	data, err := os.ReadFile(filepath.Join(tmp, "brewer.yaml"))
	if err != nil {
		t.Fatalf("brewer.yaml not created: %v", err)
	}
	if !strings.Contains(string(data), "services:") {
		t.Error("expected 'services:' key in generated config")
	}
}

func TestInit_FailsIfFileExists(t *testing.T) {
	tmp := t.TempDir()
	prev, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer os.Chdir(prev)

	_ = os.WriteFile("brewer.yaml", []byte("existing"), 0o644)

	cmd := Init()
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when file already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestInit_ForceOverwrites(t *testing.T) {
	tmp := t.TempDir()
	prev, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer os.Chdir(prev)

	_ = os.WriteFile("brewer.yaml", []byte("old content"), 0o644)

	cmd := Init()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--force"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error with --force: %v", err)
	}

	data, _ := os.ReadFile("brewer.yaml")
	if string(data) == "old content" {
		t.Error("expected file to be overwritten")
	}
}
