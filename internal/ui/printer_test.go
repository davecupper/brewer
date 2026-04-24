package ui

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func newBuf() (*Printer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	p := NewPrinterTo(buf)
	return p, buf
}

func TestServiceStarting(t *testing.T) {
	p, buf := newBuf()
	p.ServiceStarting("postgres")
	if !strings.Contains(buf.String(), "postgres") {
		t.Errorf("expected service name in output, got: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "starting") {
		t.Errorf("expected 'starting' in output, got: %q", buf.String())
	}
}

func TestServiceRunning(t *testing.T) {
	p, buf := newBuf()
	p.ServiceRunning("redis", 1234)
	out := buf.String()
	if !strings.Contains(out, "redis") {
		t.Errorf("expected service name, got: %q", out)
	}
	if !strings.Contains(out, "1234") {
		t.Errorf("expected pid in output, got: %q", out)
	}
}

func TestServiceFailed(t *testing.T) {
	p, buf := newBuf()
	p.ServiceFailed("api", errors.New("exit status 1"))
	out := buf.String()
	if !strings.Contains(out, "failed") {
		t.Errorf("expected 'failed' in output, got: %q", out)
	}
	if !strings.Contains(out, "exit status 1") {
		t.Errorf("expected error message in output, got: %q", out)
	}
}

func TestServiceStopped(t *testing.T) {
	p, buf := newBuf()
	p.ServiceStopped("worker")
	if !strings.Contains(buf.String(), "stopped") {
		t.Errorf("expected 'stopped' in output, got: %q", buf.String())
	}
}

func TestServiceSkipped(t *testing.T) {
	p, buf := newBuf()
	p.ServiceSkipped("mailer", "dependency unavailable")
	out := buf.String()
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected 'skipped' in output, got: %q", out)
	}
	if !strings.Contains(out, "dependency unavailable") {
		t.Errorf("expected reason in output, got: %q", out)
	}
}

func TestUptime(t *testing.T) {
	p, buf := newBuf()
	p.Uptime("db", 90*time.Second)
	out := buf.String()
	if !strings.Contains(out, "uptime") {
		t.Errorf("expected 'uptime' in output, got: %q", out)
	}
}

func TestNoColorOutput(t *testing.T) {
	p, buf := newBuf()
	p.ServiceRunning("svc", 42)
	if strings.Contains(buf.String(), "\033[") {
		t.Errorf("expected no ANSI codes in non-color output")
	}
}
