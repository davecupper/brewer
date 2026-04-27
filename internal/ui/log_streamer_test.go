package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestStreamLines_AllLines(t *testing.T) {
	var buf bytes.Buffer
	ls := NewLogStreamerTo(&buf, "redis", 0)
	lines := []string{"starting", "ready", "accepting connections"}
	ls.StreamLines(lines)

	out := buf.String()
	for _, line := range lines {
		if !strings.Contains(out, line) {
			t.Errorf("expected output to contain %q, got:\n%s", line, out)
		}
	}
	if !strings.Contains(out, "[redis]") {
		t.Errorf("expected output to contain service name [redis], got:\n%s", out)
	}
}

func TestStreamLines_TailLimitsOutput(t *testing.T) {
	var buf bytes.Buffer
	ls := NewLogStreamerTo(&buf, "postgres", 2)
	lines := []string{"line1", "line2", "line3", "line4"}
	ls.StreamLines(lines)

	out := buf.String()
	if strings.Contains(out, "line1") || strings.Contains(out, "line2") {
		t.Errorf("expected tail=2 to omit first lines, got:\n%s", out)
	}
	if !strings.Contains(out, "line3") || !strings.Contains(out, "line4") {
		t.Errorf("expected tail=2 to include last 2 lines, got:\n%s", out)
	}
}

func TestStream_ReadsFromReader(t *testing.T) {
	var buf bytes.Buffer
	ls := NewLogStreamerTo(&buf, "api", 0)

	input := "hello world\nfoo bar\n"
	err := ls.Stream(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected output to contain 'hello world', got:\n%s", out)
	}
	if !strings.Contains(out, "foo bar") {
		t.Errorf("expected output to contain 'foo bar', got:\n%s", out)
	}
	if !strings.Contains(out, "[api]") {
		t.Errorf("expected output to contain '[api]', got:\n%s", out)
	}
}

func TestStreamLines_EmptyLines(t *testing.T) {
	var buf bytes.Buffer
	ls := NewLogStreamerTo(&buf, "svc", 5)
	ls.StreamLines([]string{})
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty lines, got: %s", buf.String())
	}
}

func TestStreamLines_TailLargerThanInput(t *testing.T) {
	var buf bytes.Buffer
	ls := NewLogStreamerTo(&buf, "worker", 10)
	lines := []string{"only", "three", "lines"}
	ls.StreamLines(lines)

	out := buf.String()
	for _, line := range lines {
		if !strings.Contains(out, line) {
			t.Errorf("expected output to contain %q when tail exceeds line count, got:\n%s", line, out)
		}
	}
}
