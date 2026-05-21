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

func TestBoardHandlerList(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		listFn: func(_ context.Context) ([]domain.Board, error) {
			return []domain.Board{{ID: "1", Name: "A"}}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/boards", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var boards []domain.Board
	if err := json.NewDecoder(rec.Body).Decode(&boards); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(boards) != 1 {
		t.Fatalf("expected 1 board, got %d", len(boards))
	}
}

func TestBoardHandlerCreate(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		createFn: func(_ context.Context, input service.CreateBoardInput) (domain.Board, error) {
			return domain.Board{ID: "new-id", Name: input.Name}, nil
		},
	})

	body, _ := json.Marshal(CreateBoardRequest{Name: "My Board"})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var board domain.Board
	if err := json.NewDecoder(rec.Body).Decode(&board); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if board.Name != "My Board" {
		t.Fatalf("expected 'My Board', got %q", board.Name)
	}
}

func TestBoardHandlerCreateWithTemplate(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		createFn: func(_ context.Context, input service.CreateBoardInput) (domain.Board, error) {
			if input.Template != "basic-kanban" {
				t.Fatalf("expected template 'basic-kanban', got %q", input.Template)
			}
			return domain.Board{ID: "new-id", Name: "Basic Kanban"}, nil
		},
	})

	tmpl := "basic-kanban"
	body, _ := json.Marshal(CreateBoardRequest{Name: "ignored", Template: &tmpl})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var board domain.Board
	if err := json.NewDecoder(rec.Body).Decode(&board); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if board.Name != "Basic Kanban" {
		t.Fatalf("expected 'Basic Kanban', got %q", board.Name)
	}
}

func TestBoardHandlerCreateInvalidJSON(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{})
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestBoardHandlerGet(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		getFn: func(_ context.Context, id string) (domain.BoardDetail, error) {
			return domain.BoardDetail{
				Board:   domain.Board{ID: id, Name: "Board"},
				Columns: []domain.ColumnWithCards{},
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/boards/board-1", nil)
	req.SetPathValue("id", "board-1")
	rec := httptest.NewRecorder()
	h.Get(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestBoardHandlerGetNotFound(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		getFn: func(_ context.Context, id string) (domain.BoardDetail, error) {
			return domain.BoardDetail{}, domain.ErrNotFound
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/boards/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()
	h.Get(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestBoardHandlerUpdate(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		updateFn: func(_ context.Context, id, name string) (domain.Board, error) {
			return domain.Board{ID: id, Name: name}, nil
		},
	})

	body, _ := json.Marshal(UpdateBoardRequest{Name: "Updated"})
	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "board-1")
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestBoardHandlerUpdateInvalidInput(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		updateFn: func(_ context.Context, id, name string) (domain.Board, error) {
			return domain.Board{}, domain.ErrInvalidInput
		},
	})

	body, _ := json.Marshal(UpdateBoardRequest{Name: ""})
	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "board-1")
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestBoardHandlerDelete(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		deleteFn: func(_ context.Context, id string) error {
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/boards/board-1", nil)
	req.SetPathValue("id", "board-1")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestBoardHandlerDeleteNotFound(t *testing.T) {
	t.Parallel()

	h := NewBoardHandler(&mockBoardService{
		deleteFn: func(_ context.Context, id string) error {
			return domain.ErrNotFound
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/boards/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

// -- mocks --

type mockBoardService struct {
	listFn   func(context.Context) ([]domain.Board, error)
	createFn func(context.Context, service.CreateBoardInput) (domain.Board, error)
	getFn    func(context.Context, string) (domain.BoardDetail, error)
	updateFn func(context.Context, string, string) (domain.Board, error)
	deleteFn func(context.Context, string) error
}

func (m *mockBoardService) ListBoards(ctx context.Context) ([]domain.Board, error) {
	return m.listFn(ctx)
}
func (m *mockBoardService) CreateBoard(ctx context.Context, input service.CreateBoardInput) (domain.Board, error) {
	return m.createFn(ctx, input)
}
func (m *mockBoardService) GetBoard(ctx context.Context, id string) (domain.BoardDetail, error) {
	return m.getFn(ctx, id)
}
func (m *mockBoardService) UpdateBoard(ctx context.Context, id, name string) (domain.Board, error) {
	return m.updateFn(ctx, id, name)
}
func (m *mockBoardService) DeleteBoard(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

var _ boardService = (*mockBoardService)(nil)