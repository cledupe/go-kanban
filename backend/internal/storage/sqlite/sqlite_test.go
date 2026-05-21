package sqlite

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/cledupe/go-kanban/backend/internal/domain"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	return db
}

func TestBoardRepositoryCreateAndGet(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	repo := NewBoardRepository(db)

	board, err := repo.Create(ctx, domain.Board{
		Name: "Test Board",
	})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}

	if board.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if board.Name != "Test Board" {
		t.Fatalf("expected name 'Test Board', got %q", board.Name)
	}
	if board.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
	if board.UpdatedAt.IsZero() {
		t.Fatal("expected non-zero UpdatedAt")
	}

	got, err := repo.GetByID(ctx, board.ID)
	if err != nil {
		t.Fatalf("get board: %v", err)
	}
	if got.Name != board.Name {
		t.Fatalf("expected name %q, got %q", board.Name, got.Name)
	}
}

func TestBoardRepositoryGetByIDNotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewBoardRepository(db)

	_, err := repo.GetByID(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent board")
	}
}

func TestBoardRepositoryList(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	repo := NewBoardRepository(db)

	if _, err := repo.Create(ctx, domain.Board{Name: "A"}); err != nil {
		t.Fatalf("create board A: %v", err)
	}
	if _, err := repo.Create(ctx, domain.Board{Name: "B"}); err != nil {
		t.Fatalf("create board B: %v", err)
	}

	boards, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list boards: %v", err)
	}
	if len(boards) != 2 {
		t.Fatalf("expected 2 boards, got %d", len(boards))
	}
}

func TestBoardRepositoryUpdate(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	repo := NewBoardRepository(db)

	board, err := repo.Create(ctx, domain.Board{Name: "Original"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}

	board.Name = "Updated"
	updated, err := repo.Update(ctx, board)
	if err != nil {
		t.Fatalf("update board: %v", err)
	}
	if updated.Name != "Updated" {
		t.Fatalf("expected name 'Updated', got %q", updated.Name)
	}

	got, err := repo.GetByID(ctx, board.ID)
	if err != nil {
		t.Fatalf("get board after update: %v", err)
	}
	if got.Name != "Updated" {
		t.Fatalf("expected name 'Updated', got %q", got.Name)
	}
}

func TestBoardRepositoryDelete(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	repo := NewBoardRepository(db)

	board, err := repo.Create(ctx, domain.Board{Name: "To Delete"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}

	if err := repo.Delete(ctx, board.ID); err != nil {
		t.Fatalf("delete board: %v", err)
	}

	if _, err := repo.GetByID(ctx, board.ID); err == nil {
		t.Fatal("expected error getting deleted board")
	}
}

func TestBoardRepositoryDeleteNotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewBoardRepository(db)

	err := repo.Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent board")
	}
}

func TestColumnRepositoryCreateAndGet(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)

	board, err := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}

	column, err := columnRepo.Create(ctx, domain.Column{
		BoardID:  board.ID,
		Name:     "To Do",
		Position: 0,
	})
	if err != nil {
		t.Fatalf("create column: %v", err)
	}

	if column.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if column.BoardID != board.ID {
		t.Fatalf("expected board ID %q, got %q", board.ID, column.BoardID)
	}

	got, err := columnRepo.GetByID(ctx, column.ID)
	if err != nil {
		t.Fatalf("get column: %v", err)
	}
	if got.Name != "To Do" {
		t.Fatalf("expected name 'To Do', got %q", got.Name)
	}
}

func TestColumnRepositoryListByBoardID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)

	board, err := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}

	if _, err := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "A", Position: 0}); err != nil {
		t.Fatalf("create column A: %v", err)
	}
	if _, err := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "B", Position: 1}); err != nil {
		t.Fatalf("create column B: %v", err)
	}

	columns, err := columnRepo.ListByBoardID(ctx, board.ID)
	if err != nil {
		t.Fatalf("list columns: %v", err)
	}
	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}
	if columns[0].Position != 0 || columns[1].Position != 1 {
		t.Fatal("columns not ordered by position")
	}
}

func TestColumnRepositoryUpdate(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)

	board, err := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}

	column, err := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "Original", Position: 0})
	if err != nil {
		t.Fatalf("create column: %v", err)
	}

	column.Name = "Renamed"
	column.Position = 5
	updated, err := columnRepo.Update(ctx, column)
	if err != nil {
		t.Fatalf("update column: %v", err)
	}
	if updated.Name != "Renamed" {
		t.Fatalf("expected name 'Renamed', got %q", updated.Name)
	}

	got, err := columnRepo.GetByID(ctx, column.ID)
	if err != nil {
		t.Fatalf("get column after update: %v", err)
	}
	if got.Name != "Renamed" || got.Position != 5 {
		t.Fatalf("expected (Renamed, 5), got (%q, %d)", got.Name, got.Position)
	}
}

func TestColumnRepositoryDelete(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)

	board, err := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}

	column, err := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "To Delete", Position: 0})
	if err != nil {
		t.Fatalf("create column: %v", err)
	}

	if err := columnRepo.Delete(ctx, column.ID); err != nil {
		t.Fatalf("delete column: %v", err)
	}

	if _, err := columnRepo.GetByID(ctx, column.ID); err == nil {
		t.Fatal("expected error getting deleted column")
	}
}

func TestColumnRepositoryDeleteNotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewColumnRepository(db)

	err := repo.Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent column")
	}
}

func TestCardRepositoryCreateAndGet(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)
	cardRepo := NewCardRepository(db)

	board, _ := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	column, _ := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "To Do", Position: 0})

	card, err := cardRepo.Create(ctx, domain.Card{
		ColumnID:    column.ID,
		Title:       "Test Card",
		Description: "A description",
		Position:    0,
	})
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	if card.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if card.Title != "Test Card" {
		t.Fatalf("expected title 'Test Card', got %q", card.Title)
	}
	if card.Description != "A description" {
		t.Fatalf("expected description 'A description', got %q", card.Description)
	}

	got, err := cardRepo.GetByID(ctx, card.ID)
	if err != nil {
		t.Fatalf("get card: %v", err)
	}
	if got.Title != card.Title {
		t.Fatalf("expected title %q, got %q", card.Title, got.Title)
	}
}

func TestCardRepositoryListByColumnID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)
	cardRepo := NewCardRepository(db)

	board, _ := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	col, _ := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "Col", Position: 0})

	if _, err := cardRepo.Create(ctx, domain.Card{ColumnID: col.ID, Title: "A", Position: 0}); err != nil {
		t.Fatalf("create card A: %v", err)
	}
	if _, err := cardRepo.Create(ctx, domain.Card{ColumnID: col.ID, Title: "B", Position: 1}); err != nil {
		t.Fatalf("create card B: %v", err)
	}

	cards, err := cardRepo.ListByColumnID(ctx, col.ID)
	if err != nil {
		t.Fatalf("list cards: %v", err)
	}
	if len(cards) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(cards))
	}
}

func TestCardRepositoryUpdate(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)
	cardRepo := NewCardRepository(db)

	board, _ := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	col, _ := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "Col", Position: 0})

	card, err := cardRepo.Create(ctx, domain.Card{ColumnID: col.ID, Title: "Original", Position: 0})
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	card.Title = "Updated"
	card.Description = "New desc"
	card.Position = 3
	updated, err := cardRepo.Update(ctx, card)
	if err != nil {
		t.Fatalf("update card: %v", err)
	}
	if updated.Title != "Updated" {
		t.Fatalf("expected title 'Updated', got %q", updated.Title)
	}

	got, err := cardRepo.GetByID(ctx, card.ID)
	if err != nil {
		t.Fatalf("get card after update: %v", err)
	}
	if got.Title != "Updated" || got.Description != "New desc" || got.Position != 3 {
		t.Fatalf("expected (Updated, New desc, 3), got (%q, %q, %d)", got.Title, got.Description, got.Position)
	}
}

func TestCardRepositoryDelete(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)
	cardRepo := NewCardRepository(db)

	board, _ := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	col, _ := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "Col", Position: 0})

	card, err := cardRepo.Create(ctx, domain.Card{ColumnID: col.ID, Title: "To Delete", Position: 0})
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	if err := cardRepo.Delete(ctx, card.ID); err != nil {
		t.Fatalf("delete card: %v", err)
	}

	if _, err := cardRepo.GetByID(ctx, card.ID); err == nil {
		t.Fatal("expected error getting deleted card")
	}
}

func TestCardRepositoryDeleteNotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewCardRepository(db)

	err := repo.Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent card")
	}
}

func TestCascadeDeleteBoardRemovesColumnsAndCards(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)
	cardRepo := NewCardRepository(db)

	board, _ := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	col, _ := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "Col", Position: 0})
	cardRepo.Create(ctx, domain.Card{ColumnID: col.ID, Title: "Card", Position: 0})

	if err := boardRepo.Delete(ctx, board.ID); err != nil {
		t.Fatalf("delete board: %v", err)
	}

	cols, _ := columnRepo.ListByBoardID(ctx, board.ID)
	if len(cols) != 0 {
		t.Fatal("expected columns to be cascade deleted")
	}

	cards, _ := cardRepo.ListByColumnID(ctx, col.ID)
	if len(cards) != 0 {
		t.Fatal("expected cards to be cascade deleted")
	}
}

func TestCascadeDeleteColumnRemovesCards(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)
	cardRepo := NewCardRepository(db)

	board, _ := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	col, _ := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "Col", Position: 0})
	cardRepo.Create(ctx, domain.Card{ColumnID: col.ID, Title: "Card", Position: 0})

	if err := columnRepo.Delete(ctx, col.ID); err != nil {
		t.Fatalf("delete column: %v", err)
	}

	cards, _ := cardRepo.ListByColumnID(ctx, col.ID)
	if len(cards) != 0 {
		t.Fatal("expected cards to be cascade deleted")
	}
}

func TestCreateBoardGeneratesID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewBoardRepository(db)

	b1, _ := repo.Create(context.Background(), domain.Board{Name: "A"})
	b2, _ := repo.Create(context.Background(), domain.Board{Name: "B"})

	if b1.ID == b2.ID {
		t.Fatal("expected different IDs")
	}
}

func TestCardRepositoryMove(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	ctx := context.Background()
	boardRepo := NewBoardRepository(db)
	columnRepo := NewColumnRepository(db)
	cardRepo := NewCardRepository(db)

	board, _ := boardRepo.Create(ctx, domain.Board{Name: "Board"})
	col1, _ := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "Col 1", Position: 0})
	col2, _ := columnRepo.Create(ctx, domain.Column{BoardID: board.ID, Name: "Col 2", Position: 1})

	card, err := cardRepo.Create(ctx, domain.Card{ColumnID: col1.ID, Title: "Card", Position: 0})
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	if err := cardRepo.Move(ctx, card.ID, col2.ID, 5); err != nil {
		t.Fatalf("move card: %v", err)
	}

	got, err := cardRepo.GetByID(ctx, card.ID)
	if err != nil {
		t.Fatalf("get card after move: %v", err)
	}
	if got.ColumnID != col2.ID {
		t.Fatalf("expected column %q, got %q", col2.ID, got.ColumnID)
	}
	if got.Position != 5 {
		t.Fatalf("expected position 5, got %d", got.Position)
	}
}

func TestCardRepositoryMoveNotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	repo := NewCardRepository(db)

	err := repo.Move(context.Background(), "nonexistent", "col-1", 0)
	if err == nil {
		t.Fatal("expected error for nonexistent card")
	}
}
