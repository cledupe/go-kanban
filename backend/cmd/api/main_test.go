package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cledupe/go-kanban/backend/internal/app"
)

func TestRealMainPassesConfigToRunner(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	getenv := func(key string) string {
		switch key {
		case "HOST":
			return "127.0.0.1"
		case "PORT":
			return "8080"
		case "DB_PATH":
			return filepath.Join(t.TempDir(), "test.db")
		default:
			return ""
		}
	}

	var got app.Config
	run := func(gotCtx context.Context, cfg app.Config) error {
		if gotCtx != ctx {
			t.Fatal("expected context to be passed through")
		}

		got = cfg
		return nil
	}

	if err := realMain(ctx, getenv, run); err != nil {
		t.Fatalf("realMain returned error: %v", err)
	}

	if got.Host != "127.0.0.1" || got.Port != "8080" || got.DBPath == "" {
		t.Fatalf("unexpected config passed to runner: %+v", got)
	}
}

func TestRealMainReturnsConfigurationError(t *testing.T) {
	t.Parallel()

	runCalled := false
	run := func(context.Context, app.Config) error {
		runCalled = true
		return nil
	}

	err := realMain(context.Background(), func(string) string { return "" }, run)
	if err == nil {
		t.Fatal("expected configuration error")
	}

	if runCalled {
		t.Fatal("runner should not be called when config is invalid")
	}
}

func TestRealMainReturnsRunnerError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("boom")
	getenv := func(key string) string {
		switch key {
		case "HOST":
			return "127.0.0.1"
		case "PORT":
			return "8080"
		case "DB_PATH":
			return "/tmp/test.db"
		default:
			return ""
		}
	}

	run := func(context.Context, app.Config) error {
		return expectedErr
	}

	err := realMain(context.Background(), getenv, run)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestMainExitsWhenConfigurationIsInvalid(t *testing.T) {
	t.Parallel()

	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcessMain")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "HOST=", "PORT=", "DB_PATH=")

	err := cmd.Run()
	if err == nil {
		t.Fatal("expected helper process to fail")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected exit error, got %v", err)
	}

	if exitErr.ExitCode() == 0 {
		t.Fatal("expected non-zero exit code")
	}
}

func TestHelperProcessMain(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	main()
}
