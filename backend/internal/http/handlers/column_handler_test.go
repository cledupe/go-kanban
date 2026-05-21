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

func TestColumnHandlerCreate(t *testing.T) {
	t.Parallel()

	h := NewColumnHandler(&mockColumnService{
		createFn: func(_ context.Context, input service.CreateColumnInput) (domain.Column, error) {
			return domain.Column{ID: "col-1", BoardID: input.BoardID, Name: input.Name, Position: 0}, nil
		},
	})

	body, _ := json.Marshal(CreateColumnRequest{Name: "To Do"})
	req := httptest.NewRequest(http.MethodPost, "/api/boards/board-1/columns", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("boardId", "board-1")
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var col domain.Column
	if err := json.NewDecoder(rec.Body).Decode(&col); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if col.Name != "To Do" {
		t.Fatalf("expected 'To Do', got %q", col.Name)
	}
}

func TestColumnHandlerCreateBoardNotFound(t *testing.T) {
	t.Parallel()

	h := NewColumnHandler(&mockColumnService{
		createFn: func(_ context.Context, input service.CreateColumnInput) (domain.Column, error) {
			return domain.Column{}, domain.ErrNotFound
		},
	})

	body, _ := json.Marshal(CreateColumnRequest{Name: "Col"})
	req := httptest.NewRequest(http.MethodPost, "/api/boards/nonexistent/columns", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("boardId", "nonexistent")
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestColumnHandlerUpdate(t *testing.T) {
	t.Parallel()

	h := NewColumnHandler(&mockColumnService{
		updateFn: func(_ context.Context, input service.UpdateColumnInput) (domain.Column, error) {
			return domain.Column{ID: input.ID, Name: *input.Name, Position: *input.Position}, nil
		},
	})

	name := "Renamed"
	pos := 2
	body, _ := json.Marshal(UpdateColumnRequest{Name: &name, Position: &pos})
	req := httptest.NewRequest(http.MethodPatch, "/api/columns/col-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "col-1")
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestColumnHandlerUpdateNotFound(t *testing.T) {
	t.Parallel()

	h := NewColumnHandler(&mockColumnService{
		updateFn: func(_ context.Context, input service.UpdateColumnInput) (domain.Column, error) {
			return domain.Column{}, domain.ErrNotFound
		},
	})

	body, _ := json.Marshal(UpdateColumnRequest{})
	req := httptest.NewRequest(http.MethodPatch, "/api/columns/nonexistent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestColumnHandlerDelete(t *testing.T) {
	t.Parallel()

	h := NewColumnHandler(&mockColumnService{
		deleteFn: func(_ context.Context, id string) error {
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/columns/col-1", nil)
	req.SetPathValue("id", "col-1")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestColumnHandlerDeleteNotFound(t *testing.T) {
	t.Parallel()

	h := NewColumnHandler(&mockColumnService{
		deleteFn: func(_ context.Context, id string) error {
			return domain.ErrNotFound
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/columns/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

// -- mocks --

type mockColumnService struct {
	createFn func(context.Context, service.CreateColumnInput) (domain.Column, error)
	updateFn func(context.Context, service.UpdateColumnInput) (domain.Column, error)
	deleteFn func(context.Context, string) error
}

func (m *mockColumnService) CreateColumn(ctx context.Context, input service.CreateColumnInput) (domain.Column, error) {
	return m.createFn(ctx, input)
}
func (m *mockColumnService) UpdateColumn(ctx context.Context, input service.UpdateColumnInput) (domain.Column, error) {
	return m.updateFn(ctx, input)
}
func (m *mockColumnService) DeleteColumn(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

var _ columnService = (*mockColumnService)(nil)