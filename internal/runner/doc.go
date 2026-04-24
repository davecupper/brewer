// Package runner provides process lifecycle management for brewer services.
//
// It is responsible for starting, stopping, and monitoring OS-level processes
// corresponding to services defined in a brewer.yaml configuration file.
//
// Usage:
//
//	r := runner.New()
//	err := r.Start(ctx, svc)
//	err  = r.Stop(svc.Name)
//	status, ok := r.GetStatus(svc.Name)
//
// The Runner is safe for concurrent use across goroutines. Each service is
// tracked by name and its status transitions through Idle → Running → Stopped
// (or Failed if startup fails).
package runner
