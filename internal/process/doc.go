// Package process provides a thin wrapper around os/exec.Cmd that tracks
// the lifecycle of a single managed service process.
//
// A Process is created via New, started with Start, and stopped with Stop.
// At any point a read-only Snapshot can be retrieved without holding a lock
// in the caller.
//
// This package is intentionally low-level; higher-level orchestration lives
// in the runner package.
package process
