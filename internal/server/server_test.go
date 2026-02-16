package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRunMissingPetstoreUser(t *testing.T) {
	t.Setenv("PETSTORE_USER", "")
	t.Setenv("JWT_SECRET",
		"some-secret-that-is-long-enough-32b")

	err := Run(context.Background())
	if err == nil {
		t.Fatal("expected error for missing PETSTORE_USER")
	}
}

func TestRunMissingJWTSecret(t *testing.T) {
	t.Setenv("PETSTORE_USER", "petstore")
	t.Setenv("PETSTORE_PASSWORD", "secret")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("JWT_SECRET", "")

	err := Run(context.Background())
	if err == nil {
		t.Fatal("expected error for missing JWT_SECRET")
	}
	want := "JWT_SECRET is required"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestBuildShortJWTSecret(t *testing.T) {
	_, err := build(nil, "short", true)
	if err == nil {
		t.Fatal("expected error for short JWT secret")
	}
}

func TestServeGracefulShutdown(t *testing.T) {
	// Find a free port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("finding free port: %v", err)
	}
	addr := ln.Addr().String()
	if err := ln.Close(); err != nil {
		t.Fatalf("closing listener: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	h := http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			if _, err := w.Write([]byte("ok\n")); err != nil {
				return
			}
		},
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- serve(ctx, addr, h)
	}()

	// Wait for the server to be ready.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, dialErr := net.DialTimeout(
			"tcp", addr, 100*time.Millisecond,
		)
		if dialErr == nil {
			if err := conn.Close(); err != nil {
				t.Fatalf("closing conn: %v", err)
			}
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify the server is responding.
	resp, err := http.Get("http://" + addr + "/")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("closing body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	// Cancel context to trigger graceful shutdown.
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("serve returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal(
			"serve did not return after context cancellation",
		)
	}
}
