package integration

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestContainerTargetsBuildAndComposeStartsBackend(t *testing.T) {
	if os.Getenv("RUN_CONTAINER_TESTS") != "1" {
		t.Skip("set RUN_CONTAINER_TESTS=1 to run container smoke tests")
	}

	repoRoot := repoRoot(t)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	runCommand(t, repoRoot, ctx, "docker", "build", "--target", "test", "-f", "backend/Dockerfile", ".")

	t.Cleanup(func() {
		downCtx, downCancel := context.WithTimeout(context.Background(), time.Minute)
		defer downCancel()
		runCommand(t, repoRoot, downCtx, "docker", "compose", "down", "-v", "--remove-orphans")
	})

	runCommand(t, repoRoot, ctx, "docker", "compose", "up", "-d", "--build", "backend")

	waitForReadiness(t, "http://127.0.0.1:8080/readyz", 30*time.Second)
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve caller path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
}

func runCommand(t *testing.T, workdir string, ctx context.Context, name string, args ...string) {
	t.Helper()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = workdir

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		t.Fatalf("run %s %v: %v\n%s", name, args, err, output.String())
	}
}

func waitForReadiness(t *testing.T, url string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf("readiness endpoint %s was not healthy before timeout", url)
}
