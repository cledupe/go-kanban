package domain

import (
	"context"
	"time"
)

type Column struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ColumnRepository interface {
	Create(ctx context.Context, column Column) (Column, error)
	GetByID(ctx context.Context, id string) (Column, error)
	ListByBoardID(ctx context.Context, boardID string) ([]Column, error)
	Update(ctx context.Context, column Column) (Column, error)
	Delete(ctx context.Context, id string) error
}
