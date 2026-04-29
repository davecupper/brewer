package health

import (
	"context"
	"net"
)

// checkTCP dials the given address and returns nil on a successful connection.
func checkTCP(ctx context.Context, addr string) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	return conn.Close()
}
