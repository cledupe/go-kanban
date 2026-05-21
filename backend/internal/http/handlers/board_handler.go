package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cledupe/go-kanban/backend/internal/service"
)

type BoardHandler struct {
	boardService boardService
}

func NewBoardHandler(s boardService) *BoardHandler {
	return &BoardHandler{boardService: s}
}

func (h *BoardHandler) List(w http.ResponseWriter, r *http.Request) {
	boards, err := h.boardService.ListBoards(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, boards)
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	input := service.CreateBoardInput{Name: req.Name}
	if req.Template != nil {
		input.Template = *req.Template
	}

	board, err := h.boardService.CreateBoard(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, board)
}

func (h *BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing board id"})
		return
	}

	detail, err := h.boardService.GetBoard(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (h *BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing board id"})
		return
	}

	var req UpdateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	board, err := h.boardService.UpdateBoard(r.Context(), id, req.Name)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, board)
}

func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing board id"})
		return
	}

	if err := h.boardService.DeleteBoard(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}