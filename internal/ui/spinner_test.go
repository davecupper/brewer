package ui

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestSpinner_StartStop(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerTo("loading", &buf)

	s.Start()
	time.Sleep(200 * time.Millisecond)
	s.Stop()

	out := buf.String()
	if !strings.Contains(out, "loading") {
		t.Errorf("expected output to contain 'loading', got: %q", out)
	}
}

func TestSpinner_StopWithMessage(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerTo("working", &buf)

	s.Start()
	time.Sleep(100 * time.Millisecond)
	s.StopWithMessage("✔ done")

	out := buf.String()
	if !strings.Contains(out, "✔ done") {
		t.Errorf("expected output to contain '✔ done', got: %q", out)
	}
}

func TestSpinner_StopIdempotent(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerTo("task", &buf)

	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop()

	// Calling Stop again should not panic or block.
	done := make(chan struct{})
	go func() {
		s.Stop()
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(time.Second):
		t.Fatal("second Stop() call blocked")
	}
}

func TestSpinner_StartIdempotent(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerTo("task", &buf)

	s.Start()
	s.Start() // second start should be a no-op
	time.Sleep(50 * time.Millisecond)
	s.Stop()
}
