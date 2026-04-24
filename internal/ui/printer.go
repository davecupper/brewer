package ui

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// Printer writes formatted status messages to an output writer.
type Printer struct {
	out io.Writer
	color bool
}

// NewPrinter returns a Printer writing to os.Stdout with color auto-detected.
func NewPrinter() *Printer {
	return &Printer{out: os.Stdout, color: isTerminal(os.Stdout)}
}

// NewPrinterTo returns a Printer writing to w without color.
func NewPrinterTo(w io.Writer) *Printer {
	return &Printer{out: w, color: false}
}

// ServiceStarting prints a "starting" notice for a service.
func (p *Printer) ServiceStarting(name string) {
	p.printf(colorCyan, "→", name, "starting...")
}

// ServiceRunning prints a "running" notice with PID.
func (p *Printer) ServiceRunning(name string, pid int) {
	p.printf(colorGreen, "✓", name, fmt.Sprintf("running (pid %d)", pid))
}

// ServiceStopped prints a "stopped" notice.
func (p *Printer) ServiceStopped(name string) {
	p.printf(colorGray, "■", name, "stopped")
}

// ServiceFailed prints a failure message.
func (p *Printer) ServiceFailed(name string, err error) {
	p.printf(colorRed, "✗", name, fmt.Sprintf("failed: %v", err))
}

// ServiceSkipped prints a skip message.
func (p *Printer) ServiceSkipped(name string, reason string) {
	p.printf(colorYellow, "⚠", name, fmt.Sprintf("skipped: %s", reason))
}

// Uptime prints a service uptime line.
func (p *Printer) Uptime(name string, d time.Duration) {
	p.printf(colorGray, " ", name, fmt.Sprintf("uptime %s", d.Round(time.Second)))
}

func (p *Printer) printf(clr, icon, name, msg string) {
	if p.color {
		fmt.Fprintf(p.out, "%s%s%s  %-20s %s\n", clr, icon, colorReset, name, msg)
	} else {
		fmt.Fprintf(p.out, "%s  %-20s %s\n", icon, name, msg)
	}
}

// isTerminal reports whether f is a character device (terminal).
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
