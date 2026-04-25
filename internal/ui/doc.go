// Package ui provides terminal output utilities for the brewer CLI.
//
// It includes:
//
//   - Printer: structured log-style messages for service lifecycle events
//   - Spinner: animated progress indicators for long-running operations
//   - StatusTable: tabular display of service states and metadata
//
// All components support writing to an arbitrary io.Writer, making them
// straightforward to test without capturing os.Stdout directly.
package ui
