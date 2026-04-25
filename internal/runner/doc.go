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
// tracked by name and its status transitions through the following states:
//
//	Idle → Running → Stopped
//	Idle → Failed   (if startup fails)
//	Running → Failed (if the process exits unexpectedly)
//
// Callers can use GetStatus to poll the current state of a service, or
// subscribe to status changes via the Watch method if supported by the
// implementation.
package runner
