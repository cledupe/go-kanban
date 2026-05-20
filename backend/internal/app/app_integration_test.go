package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestAppServesReadinessEndpoint(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Host:   "127.0.0.1",
		Port:   "8080",
		DBPath: filepath.Join(t.TempDir(), "test.db"),
	}

	application, err := New(cfg)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	t.Cleanup(func() { application.db.Close() })

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recorder := httptest.NewRecorder()

	application.Handler().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(recorder.Body).Decode(&payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}

	if payload.Status != "ok" {
		t.Fatalf("expected payload status ok, got %q", payload.Status)
	}
}
