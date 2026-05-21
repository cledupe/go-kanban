package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cledupe/go-kanban/backend/internal/service"
)

type ColumnHandler struct {
	columnService columnService
}

func NewColumnHandler(s columnService) *ColumnHandler {
	return &ColumnHandler{columnService: s}
}

func (h *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("boardId")
	if boardID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing board id"})
		return
	}

	var req CreateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	col, err := h.columnService.CreateColumn(r.Context(), service.CreateColumnInput{
		BoardID: boardID,
		Name:    req.Name,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, col)
}

func (h *ColumnHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing column id"})
		return
	}

	var req UpdateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	col, err := h.columnService.UpdateColumn(r.Context(), service.UpdateColumnInput{
		ID:       id,
		Name:     req.Name,
		Position: req.Position,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, col)
}

func (h *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing column id"})
		return
	}

	if err := h.columnService.DeleteColumn(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}