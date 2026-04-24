package runner

import (
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/brewer/internal/config"
)

// Status represents the current state of a service process.
type Status int

const (
	StatusIdle Status = iota
	StatusRunning
	StatusFailed
	StatusStopped
)

// ServiceProcess tracks a running service and its metadata.
type ServiceProcess struct {
	Service config.Service
	Cmd     *exec.Cmd
	Status  Status
	mu      sync.Mutex
}

// Runner manages the lifecycle of service processes.
type Runner struct {
	processes map[string]*ServiceProcess
	mu        sync.RWMutex
}

// New creates a new Runner instance.
func New() *Runner {
	return &Runner{
		processes: make(map[string]*ServiceProcess),
	}
}

// Start launches a service process using its configured command.
func (r *Runner) Start(ctx context.Context, svc config.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if sp, exists := r.processes[svc.Name]; exists && sp.Status == StatusRunning {
		return fmt.Errorf("service %q is already running", svc.Name)
	}

	if len(svc.Command) == 0 {
		return fmt.Errorf("service %q has no command defined", svc.Name)
	}

	cmd := exec.CommandContext(ctx, svc.Command[0], svc.Command[1:]...)
	cmd.Dir = svc.WorkDir

	sp := &ServiceProcess{
		Service: svc,
		Cmd:     cmd,
		Status:  StatusRunning,
	}

	if err := cmd.Start(); err != nil {
		sp.Status = StatusFailed
		r.processes[svc.Name] = sp
		return fmt.Errorf("failed to start service %q: %w", svc.Name, err)
	}

	r.processes[svc.Name] = sp

	go func() {
		_ = cmd.Wait()
		sp.mu.Lock()
		sp.Status = StatusStopped
		sp.mu.Unlock()
	}()

	return nil
}

// Stop terminates a running service process.
func (r *Runner) Stop(name string) error {
	r.mu.RLock()
	sp, exists := r.processes[name]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %q is not tracked", name)
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.Status != StatusRunning {
		return fmt.Errorf("service %q is not running", name)
	}

	if err := sp.Cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop service %q: %w", name, err)
	}

	sp.Status = StatusStopped
	return nil
}

// GetStatus returns the current status of a service.
func (r *Runner) GetStatus(name string) (Status, bool) {
	r.mu.RLock()
	sp, exists := r.processes[name]
	r.mu.RUnlock()

	if !exists {
		return StatusIdle, false
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()
	return sp.Status, true
}
