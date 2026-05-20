package domain

import (
	"context"
	"time"
)

type Board struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BoardRepository interface {
	Create(ctx context.Context, board Board) (Board, error)
	GetByID(ctx context.Context, id string) (Board, error)
	List(ctx context.Context) ([]Board, error)
	Update(ctx context.Context, board Board) (Board, error)
	Delete(ctx context.Context, id string) error
}
