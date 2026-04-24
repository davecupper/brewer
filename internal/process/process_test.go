package process_test

import (
	"bytes"
	"runtime"
	"testing"
	"time"

	"github.com/yourorg/brewer/internal/process"
)

func TestStart_Success(t *testing.T) {
	var buf bytes.Buffer
	p := process.New("echo", []string{"echo", "hello"}, &buf, &buf)
	if err := p.Start([]string{"echo", "hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap := p.Snapshot()
	if snap.Status != process.StatusRunning {
		t.Errorf("expected running, got %s", snap.Status)
	}
	if snap.PID == 0 {
		t.Error("expected non-zero PID")
	}
	// Give the short-lived process time to finish.
	time.Sleep(50 * time.Millisecond)
}

func TestStart_NoCommand(t *testing.T) {
	p := process.New("empty", nil, nil, nil)
	if err := p.Start(nil); err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestStart_AlreadyRunning(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sleep not available on windows")
	}
	p := process.New("sleep", nil, nil, nil)
	if err := p.Start([]string{"sleep", "10"}); err != nil {
		t.Fatalf("first start: %v", err)
	}
	defer p.Stop() //nolint:errcheck

	if err := p.Start([]string{"sleep", "10"}); err == nil {
		t.Fatal("expected error on second start")
	}
}

func TestStop_Success(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sleep not available on windows")
	}
	p := process.New("sleep", nil, nil, nil)
	if err := p.Start([]string{"sleep", "30"}); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	snap := p.Snapshot()
	if snap.Status != process.StatusStopped {
		t.Errorf("expected stopped, got %s", snap.Status)
	}
}

func TestSnapshot_Uptime(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sleep not available on windows")
	}
	p := process.New("sleep", nil, nil, nil)
	if err := p.Start([]string{"sleep", "30"}); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer p.Stop() //nolint:errcheck

	time.Sleep(1100 * time.Millisecond)
	snap := p.Snapshot()
	if snap.Uptime() < time.Second {
		t.Errorf("expected uptime >= 1s, got %s", snap.Uptime())
	}
}

func TestSnapshot_IsRunning(t *testing.T) {
	p := process.New("stopped", nil, nil, nil)
	snap := p.Snapshot()
	if snap.IsRunning() {
		t.Error("newly created process should not be running")
	}
}
