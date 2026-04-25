package health_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nickcorin/brewer/internal/config"
	"github.com/nickcorin/brewer/internal/health"
)

func svc(hc *config.HealthCheck) config.Service {
	return config.Service{Name: "test", HealthCheck: hc}
}

func TestWait_NoHealthCheck(t *testing.T) {
	c := health.New(svc(nil))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := c.Wait(ctx); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWait_HTTP_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := health.New(svc(&config.HealthCheck{Type: "http", Target: ts.URL}))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Wait(ctx); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWait_HTTP_Timeout(t *testing.T) {
	c := health.New(svc(&config.HealthCheck{Type: "http", Target: "http://127.0.0.1:19999"}))
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	if err := c.Wait(ctx); err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestWait_TCP_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	c := health.New(svc(&config.HealthCheck{Type: "tcp", Target: ln.Addr().String()}))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Wait(ctx); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWait_UnknownType(t *testing.T) {
	c := health.New(svc(&config.HealthCheck{Type: "grpc", Target: "localhost:50051"}))
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	if err := c.Wait(ctx); err == nil {
		t.Fatal("expected error for unknown probe type, got nil")
	}
}
