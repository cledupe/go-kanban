package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cledupe/go-kanban/backend/internal/domain"
)

type mockBoardRepo struct {
	domain.BoardRepository
	createFn  func(ctx context.Context, board domain.Board) (domain.Board, error)
	getByIDFn func(ctx context.Context, id string) (domain.Board, error)
	listFn    func(ctx context.Context) ([]domain.Board, error)
	updateFn  func(ctx context.Context, board domain.Board) (domain.Board, error)
	deleteFn  func(ctx context.Context, id string) error
}

func (m *mockBoardRepo) Create(ctx context.Context, board domain.Board) (domain.Board, error) {
	return m.createFn(ctx, board)
}
func (m *mockBoardRepo) GetByID(ctx context.Context, id string) (domain.Board, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockBoardRepo) List(ctx context.Context) ([]domain.Board, error) {
	return m.listFn(ctx)
}
func (m *mockBoardRepo) Update(ctx context.Context, board domain.Board) (domain.Board, error) {
	return m.updateFn(ctx, board)
}
func (m *mockBoardRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

type mockColumnRepo struct {
	domain.ColumnRepository
	createFn       func(ctx context.Context, column domain.Column) (domain.Column, error)
	getByIDFn      func(ctx context.Context, id string) (domain.Column, error)
	listByBoardFn  func(ctx context.Context, boardID string) ([]domain.Column, error)
	updateFn       func(ctx context.Context, column domain.Column) (domain.Column, error)
	deleteFn       func(ctx context.Context, id string) error
}

func (m *mockColumnRepo) Create(ctx context.Context, column domain.Column) (domain.Column, error) {
	return m.createFn(ctx, column)
}
func (m *mockColumnRepo) GetByID(ctx context.Context, id string) (domain.Column, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockColumnRepo) ListByBoardID(ctx context.Context, boardID string) ([]domain.Column, error) {
	return m.listByBoardFn(ctx, boardID)
}
func (m *mockColumnRepo) Update(ctx context.Context, column domain.Column) (domain.Column, error) {
	return m.updateFn(ctx, column)
}
func (m *mockColumnRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

type mockCardRepo struct {
	domain.CardRepository
	createFn        func(ctx context.Context, card domain.Card) (domain.Card, error)
	getByIDFn       func(ctx context.Context, id string) (domain.Card, error)
	listByColumnFn  func(ctx context.Context, columnID string) ([]domain.Card, error)
	updateFn        func(ctx context.Context, card domain.Card) (domain.Card, error)
	deleteFn        func(ctx context.Context, id string) error
	moveFn          func(ctx context.Context, cardID string, targetColumnID string, position int) error
}

func (m *mockCardRepo) Create(ctx context.Context, card domain.Card) (domain.Card, error) {
	return m.createFn(ctx, card)
}
func (m *mockCardRepo) GetByID(ctx context.Context, id string) (domain.Card, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCardRepo) ListByColumnID(ctx context.Context, columnID string) ([]domain.Card, error) {
	return m.listByColumnFn(ctx, columnID)
}
func (m *mockCardRepo) Update(ctx context.Context, card domain.Card) (domain.Card, error) {
	return m.updateFn(ctx, card)
}
func (m *mockCardRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}
func (m *mockCardRepo) Move(ctx context.Context, cardID string, targetColumnID string, position int) error {
	return m.moveFn(ctx, cardID, targetColumnID, position)
}

func TestBoardServiceCreateBoard(t *testing.T) {
	t.Parallel()

	var createdColumns []domain.Column
	boardRepo := &mockBoardRepo{
		createFn: func(_ context.Context, board domain.Board) (domain.Board, error) {
			board.ID = "board-1"
			return board, nil
		},
	}

	columnRepo := &mockColumnRepo{
		createFn: func(_ context.Context, col domain.Column) (domain.Column, error) {
			col.ID = "col-" + col.Name
			createdColumns = append(createdColumns, col)
			return col, nil
		},
	}

	s := NewBoardService(boardRepo, columnRepo, &mockCardRepo{})

	board, err := s.CreateBoard(context.Background(), CreateBoardInput{Name: "My Board"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}
	if board.Name != "My Board" {
		t.Fatalf("expected name 'My Board', got %q", board.Name)
	}
	if len(createdColumns) != 0 {
		t.Fatalf("expected no columns for blank board, got %d", len(createdColumns))
	}
}

func TestBoardServiceCreateBoardFromBasicKanbanTemplate(t *testing.T) {
	t.Parallel()

	expectedColumns := []string{"To Do", "In Progress", "Done"}
	var createdColumns []domain.Column

	boardRepo := &mockBoardRepo{
		createFn: func(_ context.Context, board domain.Board) (domain.Board, error) {
			board.ID = "board-1"
			return board, nil
		},
	}

	columnRepo := &mockColumnRepo{
		createFn: func(_ context.Context, col domain.Column) (domain.Column, error) {
			col.ID = "col-" + col.Name
			createdColumns = append(createdColumns, col)
			return col, nil
		},
	}

	s := NewBoardService(boardRepo, columnRepo, &mockCardRepo{})

	board, err := s.CreateBoard(context.Background(), CreateBoardInput{Name: "ignored", Template: "basic-kanban"})
	if err != nil {
		t.Fatalf("create board from template: %v", err)
	}
	if board.Name != "Basic Kanban" {
		t.Fatalf("expected board name 'Basic Kanban', got %q", board.Name)
	}
	if len(createdColumns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(createdColumns))
	}
	for i, col := range createdColumns {
		if col.Name != expectedColumns[i] {
			t.Fatalf("column %d: expected %q, got %q", i, expectedColumns[i], col.Name)
		}
		if col.Position != i {
			t.Fatalf("column %d: expected position %d, got %d", i, i, col.Position)
		}
		if col.BoardID != "board-1" {
			t.Fatalf("column %d: expected board_id 'board-1', got %q", i, col.BoardID)
		}
	}
}

func TestBoardServiceCreateBoardFromBugTrackerTemplate(t *testing.T) {
	t.Parallel()

	expectedColumns := []string{"Backlog", "Investigating", "Fixing", "Verified"}
	var createdColumns []domain.Column

	boardRepo := &mockBoardRepo{
		createFn: func(_ context.Context, board domain.Board) (domain.Board, error) {
			board.ID = "board-1"
			return board, nil
		},
	}

	columnRepo := &mockColumnRepo{
		createFn: func(_ context.Context, col domain.Column) (domain.Column, error) {
			col.ID = "col-" + col.Name
			createdColumns = append(createdColumns, col)
			return col, nil
		},
	}

	s := NewBoardService(boardRepo, columnRepo, &mockCardRepo{})

	board, err := s.CreateBoard(context.Background(), CreateBoardInput{Name: "ignored", Template: "bug-tracker"})
	if err != nil {
		t.Fatalf("create board from template: %v", err)
	}
	if board.Name != "Bug Tracker" {
		t.Fatalf("expected board name 'Bug Tracker', got %q", board.Name)
	}
	if len(createdColumns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(createdColumns))
	}
	for i, col := range createdColumns {
		if col.Name != expectedColumns[i] {
			t.Fatalf("column %d: expected %q, got %q", i, expectedColumns[i], col.Name)
		}
		if col.Position != i {
			t.Fatalf("column %d: expected position %d, got %d", i, i, col.Position)
		}
	}
}

func TestBoardServiceCreateBoardFromContentPipelineTemplate(t *testing.T) {
	t.Parallel()

	expectedColumns := []string{"Ideas", "Drafting", "Review", "Published"}
	var createdColumns []domain.Column

	boardRepo := &mockBoardRepo{
		createFn: func(_ context.Context, board domain.Board) (domain.Board, error) {
			board.ID = "board-1"
			return board, nil
		},
	}

	columnRepo := &mockColumnRepo{
		createFn: func(_ context.Context, col domain.Column) (domain.Column, error) {
			col.ID = "col-" + col.Name
			createdColumns = append(createdColumns, col)
			return col, nil
		},
	}

	s := NewBoardService(boardRepo, columnRepo, &mockCardRepo{})

	board, err := s.CreateBoard(context.Background(), CreateBoardInput{Name: "ignored", Template: "content-pipeline"})
	if err != nil {
		t.Fatalf("create board from template: %v", err)
	}
	if board.Name != "Content Pipeline" {
		t.Fatalf("expected board name 'Content Pipeline', got %q", board.Name)
	}
	if len(createdColumns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(createdColumns))
	}
	for i, col := range createdColumns {
		if col.Name != expectedColumns[i] {
			t.Fatalf("column %d: expected %q, got %q", i, expectedColumns[i], col.Name)
		}
		if col.Position != i {
			t.Fatalf("column %d: expected position %d, got %d", i, i, col.Position)
		}
	}
}

func TestBoardServiceCreateBoardRejectsUnknownTemplate(t *testing.T) {
	t.Parallel()

	s := NewBoardService(&mockBoardRepo{}, &mockColumnRepo{}, &mockCardRepo{})

	_, err := s.CreateBoard(context.Background(), CreateBoardInput{Name: "My Board", Template: "nonexistent"})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for unknown template, got %v", err)
	}
}

func TestBoardServiceCreateBoardRejectsEmptyName(t *testing.T) {
	t.Parallel()

	s := NewBoardService(&mockBoardRepo{}, &mockColumnRepo{}, &mockCardRepo{})

	_, err := s.CreateBoard(context.Background(), CreateBoardInput{Name: ""})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	_, err = s.CreateBoard(context.Background(), CreateBoardInput{Name: "   "})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for whitespace-only name, got %v", err)
	}
}

func TestBoardServiceListBoards(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		listFn: func(_ context.Context) ([]domain.Board, error) {
			return []domain.Board{
				{ID: "1", Name: "A"},
				{ID: "2", Name: "B"},
			}, nil
		},
	}

	s := NewBoardService(boardRepo, &mockColumnRepo{}, &mockCardRepo{})

	boards, err := s.ListBoards(context.Background())
	if err != nil {
		t.Fatalf("list boards: %v", err)
	}
	if len(boards) != 2 {
		t.Fatalf("expected 2 boards, got %d", len(boards))
	}
}

func TestBoardServiceGetBoard(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Board, error) {
			return domain.Board{ID: id, Name: "My Board"}, nil
		},
	}

	columnRepo := &mockColumnRepo{
		listByBoardFn: func(_ context.Context, boardID string) ([]domain.Column, error) {
			return []domain.Column{
				{ID: "col-1", BoardID: boardID, Name: "To Do", Position: 0},
				{ID: "col-2", BoardID: boardID, Name: "Done", Position: 1},
			}, nil
		},
	}

	cardRepo := &mockCardRepo{
		listByColumnFn: func(_ context.Context, columnID string) ([]domain.Card, error) {
			if columnID == "col-1" {
				return []domain.Card{
					{ID: "card-1", ColumnID: columnID, Title: "Task A", Position: 0},
				}, nil
			}
			return nil, nil
		},
	}

	s := NewBoardService(boardRepo, columnRepo, cardRepo)

	detail, err := s.GetBoard(context.Background(), "board-1")
	if err != nil {
		t.Fatalf("get board: %v", err)
	}
	if detail.Name != "My Board" {
		t.Fatalf("expected board name 'My Board', got %q", detail.Name)
	}
	if len(detail.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(detail.Columns))
	}
	if len(detail.Columns[0].Cards) != 1 {
		t.Fatalf("expected 1 card in first column, got %d", len(detail.Columns[0].Cards))
	}
}

func TestBoardServiceGetBoardNotFound(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Board, error) {
			return domain.Board{}, domain.ErrNotFound
		},
	}

	s := NewBoardService(boardRepo, &mockColumnRepo{}, &mockCardRepo{})

	_, err := s.GetBoard(context.Background(), "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestBoardServiceUpdateBoard(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Board, error) {
			return domain.Board{ID: id, Name: "Old"}, nil
		},
		updateFn: func(_ context.Context, board domain.Board) (domain.Board, error) {
			return board, nil
		},
	}

	s := NewBoardService(boardRepo, &mockColumnRepo{}, &mockCardRepo{})

	board, err := s.UpdateBoard(context.Background(), "board-1", "New Name")
	if err != nil {
		t.Fatalf("update board: %v", err)
	}
	if board.Name != "New Name" {
		t.Fatalf("expected name 'New Name', got %q", board.Name)
	}
}

func TestBoardServiceUpdateBoardRejectsEmptyName(t *testing.T) {
	t.Parallel()

	s := NewBoardService(&mockBoardRepo{}, &mockColumnRepo{}, &mockCardRepo{})

	_, err := s.UpdateBoard(context.Background(), "board-1", "")
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestBoardServiceDeleteBoard(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		deleteFn: func(_ context.Context, id string) error {
			return nil
		},
	}

	s := NewBoardService(boardRepo, &mockColumnRepo{}, &mockCardRepo{})

	err := s.DeleteBoard(context.Background(), "board-1")
	if err != nil {
		t.Fatalf("delete board: %v", err)
	}
}

func TestBoardServiceDeleteBoardNotFound(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		deleteFn: func(_ context.Context, id string) error {
			return domain.ErrNotFound
		},
	}

	s := NewBoardService(boardRepo, &mockColumnRepo{}, &mockCardRepo{})

	err := s.DeleteBoard(context.Background(), "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestBoardServiceUpdateBoardNotFound(t *testing.T) {
	t.Parallel()

	boardRepo := &mockBoardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Board, error) {
			return domain.Board{}, domain.ErrNotFound
		},
	}

	s := NewBoardService(boardRepo, &mockColumnRepo{}, &mockCardRepo{})

	_, err := s.UpdateBoard(context.Background(), "nonexistent", "Name")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}