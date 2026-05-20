package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cledupe/go-kanban/backend/internal/domain"
	"github.com/google/uuid"
)

type CardRepository struct {
	db *DB
}

func NewCardRepository(db *DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(ctx context.Context, card domain.Card) (domain.Card, error) {
	if card.ID == "" {
		card.ID = uuid.NewString()
	}
	if card.CreatedAt.IsZero() {
		card.CreatedAt = time.Now().UTC()
	}
	if card.UpdatedAt.IsZero() {
		card.UpdatedAt = time.Now().UTC()
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO cards (id, column_id, title, description, position, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		card.ID, card.ColumnID, card.Title, card.Description, card.Position,
		card.CreatedAt.Format(time.RFC3339), card.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return domain.Card{}, fmt.Errorf("create card: %w", err)
	}

	return card, nil
}

func (r *CardRepository) GetByID(ctx context.Context, id string) (domain.Card, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, column_id, title, description, position, created_at, updated_at FROM cards WHERE id = ?`, id)

	var c domain.Card
	var createdAt, updatedAt string
	if err := row.Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Position, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Card{}, fmt.Errorf("card %s: %w", id, err)
		}
		return domain.Card{}, fmt.Errorf("get card %s: %w", id, err)
	}

	c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return c, nil
}

func (r *CardRepository) ListByColumnID(ctx context.Context, columnID string) ([]domain.Card, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, column_id, title, description, position, created_at, updated_at FROM cards WHERE column_id = ? ORDER BY position`, columnID)
	if err != nil {
		return nil, fmt.Errorf("list cards by column %s: %w", columnID, err)
	}
	defer rows.Close()

	var cards []domain.Card
	for rows.Next() {
		var c domain.Card
		var createdAt, updatedAt string
		if err := rows.Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Position, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan card: %w", err)
		}
		c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		cards = append(cards, c)
	}

	return cards, rows.Err()
}

func (r *CardRepository) Update(ctx context.Context, card domain.Card) (domain.Card, error) {
	card.UpdatedAt = time.Now().UTC()

	result, err := r.db.ExecContext(ctx,
		`UPDATE cards SET title = ?, description = ?, position = ?, updated_at = ? WHERE id = ?`,
		card.Title, card.Description, card.Position, card.UpdatedAt.Format(time.RFC3339), card.ID,
	)
	if err != nil {
		return domain.Card{}, fmt.Errorf("update card %s: %w", card.ID, err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.Card{}, fmt.Errorf("update card %s: %w", card.ID, sql.ErrNoRows)
	}

	return card, nil
}

func (r *CardRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM cards WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete card %s: %w", id, err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("delete card %s: %w", id, sql.ErrNoRows)
	}

	return nil
}
