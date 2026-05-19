package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouterServesReadinessEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recorder := httptest.NewRecorder()

	NewRouter().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}
