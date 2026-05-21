package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cledupe/go-kanban/backend/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrInvalidInput) {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if errors.Is(err, domain.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	if errors.Is(err, domain.ErrConflict) {
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
}