package ui

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/jasonuc/brewer/internal/process"
)

func TestRender_Headers(t *testing.T) {
	buf := &bytes.Buffer{}
	tbl := NewStatusTableTo(buf)
	tbl.Render(nil)
	out := buf.String()
	for _, h := range []string{"NAME", "STATUS", "PID", "UPTIME"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output:\n%s", h, out)
		}
	}
}

func TestRender_RunningService(t *testing.T) {
	buf := &bytes.Buffer{}
	tbl := NewStatusTableTo(buf)
	snaps := []process.Snapshot{
		{
			Name:      "postgres",
			Status:    "running",
			PID:       5678,
			StartedAt: time.Now().Add(-2 * time.Minute),
		},
	}
	tbl.Render(snaps)
	out := buf.String()
	if !strings.Contains(out, "postgres") {
		t.Errorf("expected service name in output:\n%s", out)
	}
	if !strings.Contains(out, "5678") {
		t.Errorf("expected pid in output:\n%s", out)
	}
	if !strings.Contains(out, "running") {
		t.Errorf("expected status in output:\n%s", out)
	}
}

func TestRender_StoppedService(t *testing.T) {
	buf := &bytes.Buffer{}
	tbl := NewStatusTableTo(buf)
	snaps := []process.Snapshot{
		{
			Name:   "redis",
			Status: "stopped",
			PID:    0,
		},
	}
	tbl.Render(snaps)
	out := buf.String()
	if !strings.Contains(out, "redis") {
		t.Errorf("expected service name in output:\n%s", out)
	}
	// PID should display as dash when zero
	if !strings.Contains(out, "-") {
		t.Errorf("expected '-' for missing pid/uptime in output:\n%s", out)
	}
}

func TestRender_MultipleServices(t *testing.T) {
	buf := &bytes.Buffer{}
	tbl := NewStatusTableTo(buf)
	snaps := []process.Snapshot{
		{Name: "api", Status: "running", PID: 100, StartedAt: time.Now().Add(-10 * time.Second)},
		{Name: "worker", Status: "stopped", PID: 0},
	}
	tbl.Render(snaps)
	out := buf.String()
	if !strings.Contains(out, "api") || !strings.Contains(out, "worker") {
		t.Errorf("expected both services in output:\n%s", out)
	}
}
