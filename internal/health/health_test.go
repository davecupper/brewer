package health

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/patrickward/brewer/internal/config"
)

func svc(hcType, target, timeout, interval string) config.Service {
	return config.Service{
		Name: "test",
		HealthCheck: config.HealthCheck{
			Type:     hcType,
			Target:   target,
			Timeout:  timeout,
			Interval: interval,
		},
	}
}

func TestWait_NoHealthCheck(t *testing.T) {
	c := New(svc("", "", "", ""))
	if err := c.Wait(context.Background()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWait_HTTP_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New(svc("http", ts.URL, "5s", "100ms"))
	if err := c.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_HTTP_Timeout(t *testing.T) {
	c := New(svc("http", "http://127.0.0.1:19999", "300ms", "100ms"))
	if err := c.Wait(context.Background()); err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestWait_TCP_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	c := New(svc("tcp", ln.Addr().String(), "5s", "100ms"))
	if err := c.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_TCP_Timeout(t *testing.T) {
	c := New(svc("tcp", "127.0.0.1:19998", "300ms", "100ms"))
	if err := c.Wait(context.Background()); err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestWait_Exec_Success(t *testing.T) {
	c := New(svc("exec", "exit 0", "5s", "100ms"))
	if err := c.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_Exec_Timeout(t *testing.T) {
	c := New(svc("exec", "exit 1", "300ms", "100ms"))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.Wait(ctx); err == nil {
		t.Fatal("expected error")
	}
}

func TestWait_UnknownType(t *testing.T) {
	c := New(svc("grpc", "localhost:50051", "1s", "100ms"))
	if err := c.Wait(context.Background()); err == nil {
		t.Fatal("expected error for unknown type")
	}
}
