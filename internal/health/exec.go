package health

import (
	"context"
	"os/exec"
)

// checkExec runs a shell command and returns nil if it exits with code 0.
func checkExec(ctx context.Context, command string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	return cmd.Run()
}
