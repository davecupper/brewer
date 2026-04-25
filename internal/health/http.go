package health

import (
	"context"
	"fmt"
	"net/http"
)

// checkHTTP performs a single HTTP GET against target and returns nil
// if the response status code is 2xx.
func checkHTTP(ctx context.Context, target string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http probe: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http probe: unexpected status %d", resp.StatusCode)
	}

	return nil
}
