package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cledupe/go-kanban/backend/internal/domain"
	"github.com/cledupe/go-kanban/backend/internal/service"
)

func TestCardHandlerCreate(t *testing.T) {
	t.Parallel()

	h := NewCardHandler(&mockCardService{
		createFn: func(_ context.Context, input service.CreateCardInput) (domain.Card, error) {
			return domain.Card{ID: "card-1", ColumnID: input.ColumnID, Title: input.Title, Description: input.Description, Position: 0}, nil
		},
	})

	body, _ := json.Marshal(CreateCardRequest{Title: "My Card", Description: "Desc"})
	req := httptest.NewRequest(http.MethodPost, "/api/columns/col-1/cards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("columnId", "col-1")
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var card domain.Card
	if err := json.NewDecoder(rec.Body).Decode(&card); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if card.Title != "My Card" {
		t.Fatalf("expected 'My Card', got %q", card.Title)
	}
}

func TestCardHandlerCreateInvalidJSON(t *testing.T) {
	t.Parallel()

	h := NewCardHandler(&mockCardService{})
	req := httptest.NewRequest(http.MethodPost, "/api/columns/col-1/cards", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("columnId", "col-1")
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCardHandlerUpdate(t *testing.T) {
	t.Parallel()

	h := NewCardHandler(&mockCardService{
		updateFn: func(_ context.Context, input service.UpdateCardInput) (domain.Card, error) {
			return domain.Card{ID: input.ID, Title: *input.Title}, nil
		},
	})

	title := "Updated"
	body, _ := json.Marshal(UpdateCardRequest{Title: &title})
	req := httptest.NewRequest(http.MethodPatch, "/api/cards/card-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "card-1")
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCardHandlerUpdateNotFound(t *testing.T) {
	t.Parallel()

	h := NewCardHandler(&mockCardService{
		updateFn: func(_ context.Context, input service.UpdateCardInput) (domain.Card, error) {
			return domain.Card{}, domain.ErrNotFound
		},
	})

	body, _ := json.Marshal(UpdateCardRequest{})
	req := httptest.NewRequest(http.MethodPatch, "/api/cards/nonexistent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestCardHandlerMove(t *testing.T) {
	t.Parallel()

	h := NewCardHandler(&mockCardService{
		moveFn: func(_ context.Context, input service.MoveCardInput) error {
			return nil
		},
		updateFn: func(_ context.Context, input service.UpdateCardInput) (domain.Card, error) {
			return domain.Card{ID: input.ID, Title: "Card"}, nil
		},
	})

	body, _ := json.Marshal(MoveCardRequest{TargetColumnID: "col-2", Position: 0})
	req := httptest.NewRequest(http.MethodPost, "/api/cards/card-1/move", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "card-1")
	rec := httptest.NewRecorder()
	h.Move(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCardHandlerMoveNotFound(t *testing.T) {
	t.Parallel()

	h := NewCardHandler(&mockCardService{
		moveFn: func(_ context.Context, input service.MoveCardInput) error {
			return domain.ErrNotFound
		},
	})

	body, _ := json.Marshal(MoveCardRequest{TargetColumnID: "col-2", Position: 0})
	req := httptest.NewRequest(http.MethodPost, "/api/cards/nonexistent/move", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()
	h.Move(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestCardHandlerDelete(t *testing.T) {
	t.Parallel()

	h := NewCardHandler(&mockCardService{
		deleteFn: func(_ context.Context, id string) error {
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/cards/card-1", nil)
	req.SetPathValue("id", "card-1")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestCardHandlerDeleteNotFound(t *testing.T) {
	t.Parallel()

	h := NewCardHandler(&mockCardService{
		deleteFn: func(_ context.Context, id string) error {
			return domain.ErrNotFound
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/cards/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

// -- mocks --

type mockCardService struct {
	createFn func(context.Context, service.CreateCardInput) (domain.Card, error)
	updateFn func(context.Context, service.UpdateCardInput) (domain.Card, error)
	moveFn   func(context.Context, service.MoveCardInput) error
	deleteFn func(context.Context, string) error
}

func (m *mockCardService) CreateCard(ctx context.Context, input service.CreateCardInput) (domain.Card, error) {
	return m.createFn(ctx, input)
}
func (m *mockCardService) UpdateCard(ctx context.Context, input service.UpdateCardInput) (domain.Card, error) {
	return m.updateFn(ctx, input)
}
func (m *mockCardService) MoveCard(ctx context.Context, input service.MoveCardInput) error {
	return m.moveFn(ctx, input)
}
func (m *mockCardService) DeleteCard(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

var _ cardService = (*mockCardService)(nil)