// Package health provides health check utilities for brewer services.
//
// A Checker is created for a given service configuration and its Wait method
// blocks until the configured health check passes or the timeout is exceeded.
//
// Supported health check types:
//
//   - http  — performs an HTTP GET and expects a 2xx response
//   - tcp   — dials a TCP address and expects a successful connection
//   - exec  — runs a shell command and expects exit code 0
//
// If no health check is configured for a service, Wait returns immediately.
package health
