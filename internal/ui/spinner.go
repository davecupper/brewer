package ui

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// spinnerFrames are the characters cycled during a spinner animation.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner displays an animated spinner in the terminal while a task is running.
type Spinner struct {
	message string
	out     io.Writer
	stop    chan struct{}
	done    chan struct{}
	mu      sync.Mutex
	running bool
}

// NewSpinner creates a new Spinner writing to stdout.
func NewSpinner(message string) *Spinner {
	return NewSpinnerTo(message, os.Stdout)
}

// NewSpinnerTo creates a new Spinner writing to the provided writer.
func NewSpinnerTo(message string, out io.Writer) *Spinner {
	return &Spinner{
		message: message,
		out:     out,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

// Start begins the spinner animation in a background goroutine.
func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return
	}
	s.running = true
	go func() {
		defer close(s.done)
		for i := 0; ; i++ {
			select {
			case <-s.stop:
				fmt.Fprintf(s.out, "\r\033[K")
				return
			case <-time.After(80 * time.Millisecond):
				fmt.Fprintf(s.out, "\r%s %s", spinnerFrames[i%len(spinnerFrames)], s.message)
			}
		}
	}()
}

// Stop halts the spinner and clears the line.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	s.running = false
	close(s.stop)
	<-s.done
}

// StopWithMessage halts the spinner and prints a final message.
func (s *Spinner) StopWithMessage(msg string) {
	s.Stop()
	fmt.Fprintln(s.out, msg)
}
