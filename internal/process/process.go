package process

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Status represents the current state of a managed process.
type Status int

const (
	StatusStopped Status = iota
	StatusStarting
	StatusRunning
	StatusFailed
)

func (s Status) String() string {
	switch s {
	case StatusStopped:
		return "stopped"
	case StatusStarting:
		return "starting"
	case StatusRunning:
		return "running"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// Process wraps an os/exec.Cmd with lifecycle tracking.
type Process struct {
	mu        sync.RWMutex
	name      string
	cmd       *exec.Cmd
	status    Status
	pid       int
	startedAt time.Time
	stoppedAt time.Time
	stdout     io.Writer
	stderr     io.Writer
}

// New creates a new Process for the given service name and command.
func New(name string, args []string, stdout, stderr io.Writer) *Process {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	return &Process{
		name:   name,
		status: StatusStopped,
		stdout: stdout,
		stderr: stderr,
	}
}

// Start launches the process and updates internal state.
func (p *Process) Start(args []string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.status == StatusRunning || p.status == StatusStarting {
		return fmt.Errorf("process %q is already running", p.name)
	}
	if len(args) == 0 {
		return fmt.Errorf("process %q has no command", p.name)
	}

	p.cmd = exec.Command(args[0], args[1:]...)
	p.cmd.Stdout = p.stdout
	p.cmd.Stderr = p.stderr
	p.status = StatusStarting

	if err := p.cmd.Start(); err != nil {
		p.status = StatusFailed
		return fmt.Errorf("start %q: %w", p.name, err)
	}

	p.pid = p.cmd.Process.Pid
	p.startedAt = time.Now()
	p.status = StatusRunning
	return nil
}

// Stop sends an interrupt signal and waits for the process to exit.
func (p *Process) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.status != StatusRunning {
		return nil
	}
	if err := p.cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("signal %q: %w", p.name, err)
	}
	_ = p.cmd.Wait()
	p.status = StatusStopped
	p.stoppedAt = time.Now()
	return nil
}

// Snapshot returns a point-in-time view of the process state.
func (p *Process) Snapshot() Snapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return Snapshot{
		Name:      p.name,
		PID:       p.pid,
		Status:    p.status,
		StartedAt: p.startedAt,
		StoppedAt: p.stoppedAt,
	}
}
