package ui

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/jasonuc/brewer/internal/process"
)

// StatusTable renders a tabular view of process snapshots.
type StatusTable struct {
	out io.Writer
}

// NewStatusTable returns a StatusTable writing to os.Stdout.
func NewStatusTable() *StatusTable {
	return &StatusTable{out: os.Stdout}
}

// NewStatusTableTo returns a StatusTable writing to w.
func NewStatusTableTo(w io.Writer) *StatusTable {
	return &StatusTable{out: w}
}

// Render writes a formatted table of snapshots to the writer.
func (t *StatusTable) Render(snapshots []process.Snapshot) {
	w := tabwriter.NewWriter(t.out, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\tPID\tUPTIME")
	fmt.Fprintln(w, "----\t------\t---\t------")
	for _, s := range snapshots {
		pid := "-"
		if s.PID > 0 {
			pid = fmt.Sprintf("%d", s.PID)
		}
		uptime := "-"
		if s.Status == "running" && !s.StartedAt.IsZero() {
			uptime = time.Since(s.StartedAt).Round(time.Second).String()
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, s.Status, pid, uptime)
	}
	w.Flush()
}
