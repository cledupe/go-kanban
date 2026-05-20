package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cledupe/go-kanban/backend/internal/domain"
	"github.com/google/uuid"
)

type BoardRepository struct {
	db *DB
}

func NewBoardRepository(db *DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(ctx context.Context, board domain.Board) (domain.Board, error) {
	if board.ID == "" {
		board.ID = uuid.NewString()
	}
	if board.CreatedAt.IsZero() {
		board.CreatedAt = time.Now().UTC()
	}
	if board.UpdatedAt.IsZero() {
		board.UpdatedAt = time.Now().UTC()
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO boards (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		board.ID, board.Name, board.CreatedAt.Format(time.RFC3339), board.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return domain.Board{}, fmt.Errorf("create board: %w", err)
	}

	return board, nil
}

func (r *BoardRepository) GetByID(ctx context.Context, id string) (domain.Board, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, created_at, updated_at FROM boards WHERE id = ?`, id)

	var b domain.Board
	var createdAt, updatedAt string
	if err := row.Scan(&b.ID, &b.Name, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Board{}, fmt.Errorf("board %s: %w", id, err)
		}
		return domain.Board{}, fmt.Errorf("get board %s: %w", id, err)
	}

	b.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	b.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return b, nil
}

func (r *BoardRepository) List(ctx context.Context) ([]domain.Board, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, created_at, updated_at FROM boards ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list boards: %w", err)
	}
	defer rows.Close()

	var boards []domain.Board
	for rows.Next() {
		var b domain.Board
		var createdAt, updatedAt string
		if err := rows.Scan(&b.ID, &b.Name, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan board: %w", err)
		}
		b.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		b.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		boards = append(boards, b)
	}

	return boards, rows.Err()
}

func (r *BoardRepository) Update(ctx context.Context, board domain.Board) (domain.Board, error) {
	board.UpdatedAt = time.Now().UTC()

	result, err := r.db.ExecContext(ctx,
		`UPDATE boards SET name = ?, updated_at = ? WHERE id = ?`,
		board.Name, board.UpdatedAt.Format(time.RFC3339), board.ID,
	)
	if err != nil {
		return domain.Board{}, fmt.Errorf("update board %s: %w", board.ID, err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.Board{}, fmt.Errorf("update board %s: %w", board.ID, sql.ErrNoRows)
	}

	return board, nil
}

func (r *BoardRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM boards WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete board %s: %w", id, err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("delete board %s: %w", id, sql.ErrNoRows)
	}

	return nil
}
