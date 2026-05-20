package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cledupe/go-kanban/backend/internal/domain"
	"github.com/google/uuid"
)

type ColumnRepository struct {
	db *DB
}

func NewColumnRepository(db *DB) *ColumnRepository {
	return &ColumnRepository{db: db}
}

func (r *ColumnRepository) Create(ctx context.Context, column domain.Column) (domain.Column, error) {
	if column.ID == "" {
		column.ID = uuid.NewString()
	}
	if column.CreatedAt.IsZero() {
		column.CreatedAt = time.Now().UTC()
	}
	if column.UpdatedAt.IsZero() {
		column.UpdatedAt = time.Now().UTC()
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO columns (id, board_id, name, position, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		column.ID, column.BoardID, column.Name, column.Position,
		column.CreatedAt.Format(time.RFC3339), column.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return domain.Column{}, fmt.Errorf("create column: %w", err)
	}

	return column, nil
}

func (r *ColumnRepository) GetByID(ctx context.Context, id string) (domain.Column, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, board_id, name, position, created_at, updated_at FROM columns WHERE id = ?`, id)

	var c domain.Column
	var createdAt, updatedAt string
	if err := row.Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Column{}, fmt.Errorf("column %s: %w", id, err)
		}
		return domain.Column{}, fmt.Errorf("get column %s: %w", id, err)
	}

	c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return c, nil
}

func (r *ColumnRepository) ListByBoardID(ctx context.Context, boardID string) ([]domain.Column, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, board_id, name, position, created_at, updated_at FROM columns WHERE board_id = ? ORDER BY position`, boardID)
	if err != nil {
		return nil, fmt.Errorf("list columns by board %s: %w", boardID, err)
	}
	defer rows.Close()

	var columns []domain.Column
	for rows.Next() {
		var c domain.Column
		var createdAt, updatedAt string
		if err := rows.Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan column: %w", err)
		}
		c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		columns = append(columns, c)
	}

	return columns, rows.Err()
}

func (r *ColumnRepository) Update(ctx context.Context, column domain.Column) (domain.Column, error) {
	column.UpdatedAt = time.Now().UTC()

	result, err := r.db.ExecContext(ctx,
		`UPDATE columns SET name = ?, position = ?, updated_at = ? WHERE id = ?`,
		column.Name, column.Position, column.UpdatedAt.Format(time.RFC3339), column.ID,
	)
	if err != nil {
		return domain.Column{}, fmt.Errorf("update column %s: %w", column.ID, err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.Column{}, fmt.Errorf("update column %s: %w", column.ID, sql.ErrNoRows)
	}

	return column, nil
}

func (r *ColumnRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM columns WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete column %s: %w", id, err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("delete column %s: %w", id, sql.ErrNoRows)
	}

	return nil
}
