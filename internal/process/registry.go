package process

import (
	"fmt"
	"sync"
)

// Registry tracks all named processes for a brewer session.
type Registry struct {
	mu       sync.RWMutex
	processes map[string]*Process
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{processes: make(map[string]*Process)}
}

// Register adds a Process under its name. Returns an error if the name is
// already registered.
func (r *Registry) Register(p *Process) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.processes[p.name]; exists {
		return fmt.Errorf("process %q already registered", p.name)
	}
	r.processes[p.name] = p
	return nil
}

// Get returns the Process registered under name, or an error if not found.
func (r *Registry) Get(name string) (*Process, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.processes[name]
	if !ok {
		return nil, fmt.Errorf("process %q not found", name)
	}
	return p, nil
}

// Snapshots returns a map of name → Snapshot for all registered processes.
func (r *Registry) Snapshots() map[string]Snapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]Snapshot, len(r.processes))
	for name, p := range r.processes {
		out[name] = p.Snapshot()
	}
	return out
}

// StopAll stops every running process in the registry.
func (r *Registry) StopAll() []error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var errs []error
	for _, p := range r.processes {
		if err := p.Stop(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
