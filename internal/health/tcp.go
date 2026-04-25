package health

import (
	"context"
	"fmt"
	"net"
)

// checkTCP dials target (host:port) and returns nil on a successful
// TCP connection.
func checkTCP(ctx context.Context, target string) error {
	var d net.Dialer

	conn, err := d.DialContext(ctx, "tcp", target)
	if err != nil {
		return fmt.Errorf("tcp probe %q: %w", target, err)
	}

	return conn.Close()
}
