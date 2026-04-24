package process

import "time"

// Snapshot is an immutable point-in-time view of a Process.
type Snapshot struct {
	Name      string
	PID       int
	Status    Status
	StartedAt time.Time
	StoppedAt time.Time
}

// Uptime returns the duration the process has been running.
// Returns zero if the process is not running.
func (s Snapshot) Uptime() time.Duration {
	if s.Status != StatusRunning {
		return 0
	}
	return time.Since(s.StartedAt).Truncate(time.Second)
}

// IsRunning reports whether the process is currently running.
func (s Snapshot) IsRunning() bool {
	return s.Status == StatusRunning
}
