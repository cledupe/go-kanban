package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestLoadConfigAcceptsValidSettings(t *testing.T) {
	t.Parallel()

	cfg, err := LoadConfig(func(key string) string {
		switch key {
		case "HOST":
			return "127.0.0.1"
		case "PORT":
			return "8080"
		default:
			return ""
		}
	})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Address() != "127.0.0.1:8080" {
		t.Fatalf("expected address 127.0.0.1:8080, got %q", cfg.Address())
	}
}

func TestNewRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	if _, err := New(Config{}); err == nil {
		t.Fatal("expected invalid config error")
	}
}

func TestNewCreatesHandler(t *testing.T) {
	t.Parallel()

	application, err := New(Config{Host: "127.0.0.1", Port: "8080"})
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	if application.Handler() == nil {
		t.Fatal("expected handler to be initialized")
	}
}

func TestServeDelegatesToInjectedServer(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("serve error")
	called := false

	application := &App{
		serve: func(net.Listener) error {
			called = true
			return expectedErr
		},
	}

	err := application.Serve(nil)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if !called {
		t.Fatal("expected serve to be called")
	}
}

func TestShutdownDelegatesToInjectedServer(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("shutdown error")
	called := false

	application := &App{
		shutdown: func(context.Context) error {
			called = true
			return expectedErr
		},
	}

	err := application.Shutdown(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if !called {
		t.Fatal("expected shutdown to be called")
	}
}

func TestRunReturnsServeError(t *testing.T) {
	originalNewApp := newApp
	t.Cleanup(func() {
		newApp = originalNewApp
	})

	expectedErr := errors.New("listen failure")
	newApp = func(Config) (*App, error) {
		return &App{
			listenAndServe: func() error { return expectedErr },
			shutdown: func(context.Context) error {
				t.Fatal("shutdown should not be called when listen fails immediately")
				return nil
			},
		}, nil
	}

	err := Run(context.Background(), Config{Host: "127.0.0.1", Port: "8080"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestRunReturnsNilOnServerClosed(t *testing.T) {
	originalNewApp := newApp
	t.Cleanup(func() {
		newApp = originalNewApp
	})

	newApp = func(Config) (*App, error) {
		return &App{
			listenAndServe: func() error { return http.ErrServerClosed },
			shutdown: func(context.Context) error {
				t.Fatal("shutdown should not be called when server already closed")
				return nil
			},
		}, nil
	}

	if err := Run(context.Background(), Config{Host: "127.0.0.1", Port: "8080"}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRunShutsDownOnContextCancellation(t *testing.T) {
	originalNewApp := newApp
	t.Cleanup(func() {
		newApp = originalNewApp
	})

	listenStarted := make(chan struct{})
	releaseListen := make(chan struct{})
	shutdownCalled := make(chan struct{})

	newApp = func(Config) (*App, error) {
		return &App{
			listenAndServe: func() error {
				close(listenStarted)
				<-releaseListen
				return http.ErrServerClosed
			},
			shutdown: func(context.Context) error {
				close(shutdownCalled)
				close(releaseListen)
				return nil
			},
		}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(ctx, Config{Host: "127.0.0.1", Port: "8080"})
	}()

	select {
	case <-listenStarted:
	case <-time.After(time.Second):
		t.Fatal("listen function did not start")
	}

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Run did not return after context cancellation")
	}

	select {
	case <-shutdownCalled:
	case <-time.After(time.Second):
		t.Fatal("expected shutdown to be called")
	}
}
