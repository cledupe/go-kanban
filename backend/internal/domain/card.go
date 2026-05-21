package domain

import (
	"context"
	"time"
)

type Card struct {
	ID          string    `json:"id"`
	ColumnID    string    `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CardRepository interface {
	Create(ctx context.Context, card Card) (Card, error)
	GetByID(ctx context.Context, id string) (Card, error)
	ListByColumnID(ctx context.Context, columnID string) ([]Card, error)
	Update(ctx context.Context, card Card) (Card, error)
	Delete(ctx context.Context, id string) error
	Move(ctx context.Context, cardID string, targetColumnID string, position int) error
}
