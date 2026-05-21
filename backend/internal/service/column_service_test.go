package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cledupe/go-kanban/backend/internal/domain"
)

func TestColumnServiceCreateColumn(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Board, error) {
			return domain.Board{ID: id, Name: "Board"}, nil
		},
	}

	columnRepo := &mockColumnRepo{
		listByBoardFn: func(_ context.Context, boardID string) ([]domain.Column, error) {
			return nil, nil
		},
		createFn: func(_ context.Context, col domain.Column) (domain.Column, error) {
			col.ID = "col-1"
			return col, nil
		},
	}

	s := NewColumnService(boardRepo, columnRepo)

	col, err := s.CreateColumn(context.Background(), CreateColumnInput{
		BoardID: "board-1",
		Name:    "To Do",
	})
	if err != nil {
		t.Fatalf("create column: %v", err)
	}
	if col.Name != "To Do" {
		t.Fatalf("expected name 'To Do', got %q", col.Name)
	}
}

func TestColumnServiceCreateColumnRejectsEmptyName(t *testing.T) {
	t.Parallel()

	s := NewColumnService(&mockBoardRepo{}, &mockColumnRepo{})

	_, err := s.CreateColumn(context.Background(), CreateColumnInput{BoardID: "b", Name: ""})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestColumnServiceCreateColumnBoardNotFound(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Board, error) {
			return domain.Board{}, domain.ErrNotFound
		},
	}

	s := NewColumnService(boardRepo, &mockColumnRepo{})

	_, err := s.CreateColumn(context.Background(), CreateColumnInput{
		BoardID: "nonexistent",
		Name:    "To Do",
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestColumnServiceUpdateColumn(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{ID: id, Name: "Old", Position: 0}, nil
		},
		updateFn: func(_ context.Context, col domain.Column) (domain.Column, error) {
			return col, nil
		},
	}

	s := NewColumnService(&mockBoardRepo{}, columnRepo)

	name := "Renamed"
	pos := 5
	col, err := s.UpdateColumn(context.Background(), UpdateColumnInput{
		ID:       "col-1",
		Name:     &name,
		Position: &pos,
	})
	if err != nil {
		t.Fatalf("update column: %v", err)
	}
	if col.Name != "Renamed" {
		t.Fatalf("expected name 'Renamed', got %q", col.Name)
	}
	if col.Position != 5 {
		t.Fatalf("expected position 5, got %d", col.Position)
	}
}

func TestColumnServiceUpdateColumnRejectsEmptyName(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{ID: id, Name: "Col"}, nil
		},
	}

	s := NewColumnService(&mockBoardRepo{}, columnRepo)

	name := ""
	_, err := s.UpdateColumn(context.Background(), UpdateColumnInput{
		ID:   "col-1",
		Name: &name,
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestColumnServiceUpdateColumnNotFound(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{}, domain.ErrNotFound
		},
	}

	s := NewColumnService(&mockBoardRepo{}, columnRepo)

	_, err := s.UpdateColumn(context.Background(), UpdateColumnInput{ID: "nonexistent"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestColumnServiceDeleteColumn(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		deleteFn: func(_ context.Context, id string) error {
			return nil
		},
	}

	s := NewColumnService(&mockBoardRepo{}, columnRepo)

	err := s.DeleteColumn(context.Background(), "col-1")
	if err != nil {
		t.Fatalf("delete column: %v", err)
	}
}

func TestColumnServiceDeleteColumnNotFound(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		deleteFn: func(_ context.Context, id string) error {
			return domain.ErrNotFound
		},
	}

	s := NewColumnService(&mockBoardRepo{}, columnRepo)

	err := s.DeleteColumn(context.Background(), "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestColumnServiceAssignsNextPosition(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Board, error) {
			return domain.Board{ID: id}, nil
		},
	}

	columnRepo := &mockColumnRepo{
		listByBoardFn: func(_ context.Context, boardID string) ([]domain.Column, error) {
			return []domain.Column{
				{Position: 0},
				{Position: 1},
				{Position: 2},
			}, nil
		},
		createFn: func(_ context.Context, col domain.Column) (domain.Column, error) {
			if col.Position != 3 {
				t.Fatalf("expected position 3, got %d", col.Position)
			}
			return col, nil
		},
	}

	s := NewColumnService(boardRepo, columnRepo)

	_, err := s.CreateColumn(context.Background(), CreateColumnInput{
		BoardID: "board-1",
		Name:    "New Column",
	})
	if err != nil {
		t.Fatalf("create column: %v", err)
	}
}