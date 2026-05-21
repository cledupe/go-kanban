package handlers

import (
	"context"

	"github.com/cledupe/go-kanban/backend/internal/domain"
	"github.com/cledupe/go-kanban/backend/internal/service"
)

type boardService interface {
	ListBoards(ctx context.Context) ([]domain.Board, error)
	CreateBoard(ctx context.Context, input service.CreateBoardInput) (domain.Board, error)
	GetBoard(ctx context.Context, id string) (domain.BoardDetail, error)
	UpdateBoard(ctx context.Context, id string, name string) (domain.Board, error)
	DeleteBoard(ctx context.Context, id string) error
}

type columnService interface {
	CreateColumn(ctx context.Context, input service.CreateColumnInput) (domain.Column, error)
	UpdateColumn(ctx context.Context, input service.UpdateColumnInput) (domain.Column, error)
	DeleteColumn(ctx context.Context, id string) error
}

type cardService interface {
	CreateCard(ctx context.Context, input service.CreateCardInput) (domain.Card, error)
	UpdateCard(ctx context.Context, input service.UpdateCardInput) (domain.Card, error)
	MoveCard(ctx context.Context, input service.MoveCardInput) error
	DeleteCard(ctx context.Context, id string) error
}