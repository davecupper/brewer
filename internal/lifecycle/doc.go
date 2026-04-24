// Package lifecycle provides high-level orchestration for starting and
// stopping groups of services in dependency order.
//
// It combines the dependency graph resolved by [graph.Build] with the
// process management provided by [runner.Runner], emitting human-readable
// status messages via [ui.Printer] as each service transitions state.
//
// Typical usage:
//
//	m := lifecycle.New(cfg, runner, printer)
//	if err := m.StartAll(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer m.StopAll(ctx)
package lifecycle
