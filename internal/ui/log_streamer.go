package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

// LogStreamer streams log lines from a reader to an output writer,
// prefixing each line with a timestamp and service name.
type LogStreamer struct {
	service string
	out     io.Writer
	tail    int
}

// NewLogStreamer creates a LogStreamer that writes to stdout.
func NewLogStreamer(service string, tail int) *LogStreamer {
	return NewLogStreamerTo(os.Stdout, service, tail)
}

// NewLogStreamerTo creates a LogStreamer that writes to the provided writer.
func NewLogStreamerTo(out io.Writer, service string, tail int) *LogStreamer {
	return &LogStreamer{
		service: service,
		out:     out,
		tail:    tail,
	}
}

// Stream reads lines from r and writes them formatted to the output.
// It returns when r is exhausted or closed.
func (l *LogStreamer) Stream(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		timestamp := time.Now().Format("15:04:05")
		_, err := fmt.Fprintf(l.out, "%s  [%s]  %s\n", timestamp, l.service, line)
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

// StreamLines reads up to tail lines from the provided slice and writes them.
func (l *LogStreamer) StreamLines(lines []string) {
	start := 0
	if l.tail > 0 && len(lines) > l.tail {
		start = len(lines) - l.tail
	}
	for _, line := range lines[start:] {
		timestamp := time.Now().Format("15:04:05")
		fmt.Fprintf(l.out, "%s  [%s]  %s\n", timestamp, l.service, line)
	}
}
