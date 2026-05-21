package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cledupe/go-kanban/backend/internal/service"
)

type CardHandler struct {
	cardService cardService
}

func NewCardHandler(s cardService) *CardHandler {
	return &CardHandler{cardService: s}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	columnID := r.PathValue("columnId")
	if columnID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing column id"})
		return
	}

	var req CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	card, err := h.cardService.CreateCard(r.Context(), service.CreateCardInput{
		ColumnID:    columnID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, card)
}

func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing card id"})
		return
	}

	var req UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	card, err := h.cardService.UpdateCard(r.Context(), service.UpdateCardInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, card)
}

func (h *CardHandler) Move(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing card id"})
		return
	}

	var req MoveCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	if err := h.cardService.MoveCard(r.Context(), service.MoveCardInput{
		CardID:         id,
		TargetColumnID: req.TargetColumnID,
		Position:       req.Position,
	}); err != nil {
		writeError(w, err)
		return
	}

	card, err := h.cardService.UpdateCard(r.Context(), service.UpdateCardInput{ID: id})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, card)
}

func (h *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing card id"})
		return
	}

	if err := h.cardService.DeleteCard(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}