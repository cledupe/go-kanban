package httpapi

import (
	"encoding/json"
	"net/http"
)

type readinessResponse struct {
	Status string `json:"status"`
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(readinessResponse{Status: "ok"})
}
