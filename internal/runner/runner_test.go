package runner

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/brewer/internal/config"
)

func sleepService() config.Service {
	cmd := []string{"sleep", "10"}
	if runtime.GOOS == "windows" {
		cmd = []string{"timeout", "/t", "10"}
	}
	return config.Service{
		Name:    "sleeper",
		Command: cmd,
	}
}

func echoService() config.Service {
	return config.Service{
		Name:    "echo-svc",
		Command: []string{"echo", "hello"},
	}
}

func TestStart_Success(t *testing.T) {
	r := New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc := sleepService()
	if err := r.Start(ctx, svc); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	status, ok := r.GetStatus(svc.Name)
	if !ok {
		t.Fatal("expected service to be tracked")
	}
	if status != StatusRunning {
		t.Errorf("expected StatusRunning, got %v", status)
	}
}

func TestStart_AlreadyRunning(t *testing.T) {
	r := New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc := sleepService()
	_ = r.Start(ctx, svc)

	if err := r.Start(ctx, svc); err == nil {
		t.Fatal("expected error for already-running service")
	}
}

func TestStart_NoCommand(t *testing.T) {
	r := New()
	svc := config.Service{Name: "empty"}

	if err := r.Start(context.Background(), svc); err == nil {
		t.Fatal("expected error for service with no command")
	}
}

func TestStop_Success(t *testing.T) {
	r := New()
	ctx := context.Background()

	svc := sleepService()
	if err := r.Start(ctx, svc); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	if err := r.Stop(svc.Name); err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	status, _ := r.GetStatus(svc.Name)
	if status != StatusStopped {
		t.Errorf("expected StatusStopped, got %v", status)
	}
}

func TestStop_NotTracked(t *testing.T) {
	r := New()
	if err := r.Stop("ghost"); err == nil {
		t.Fatal("expected error for untracked service")
	}
}

func TestGetStatus_Unknown(t *testing.T) {
	r := New()
	_, ok := r.GetStatus("unknown")
	if ok {
		t.Fatal("expected false for unknown service")
	}
}
