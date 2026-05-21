package service

import (
	"context"
	"errors"
	"strings"

	"github.com/cledupe/go-kanban/backend/internal/domain"
)

type ColumnService struct {
	boardRepo  domain.BoardRepository
	columnRepo domain.ColumnRepository
}

func NewColumnService(
	boardRepo domain.BoardRepository,
	columnRepo domain.ColumnRepository,
) *ColumnService {
	return &ColumnService{
		boardRepo:  boardRepo,
		columnRepo: columnRepo,
	}
}

type CreateColumnInput struct {
	BoardID string
	Name    string
}

func (s *ColumnService) CreateColumn(ctx context.Context, input CreateColumnInput) (domain.Column, error) {
	if strings.TrimSpace(input.Name) == "" {
		return domain.Column{}, domain.ErrInvalidInput
	}

	if _, err := s.boardRepo.GetByID(ctx, input.BoardID); err != nil {
		return domain.Column{}, domain.ErrNotFound
	}

	columns, err := s.columnRepo.ListByBoardID(ctx, input.BoardID)
	if err != nil {
		return domain.Column{}, err
	}

	nextPos := 0
	for _, col := range columns {
		if col.Position >= nextPos {
			nextPos = col.Position + 1
		}
	}

	return s.columnRepo.Create(ctx, domain.Column{
		BoardID:  input.BoardID,
		Name:     strings.TrimSpace(input.Name),
		Position: nextPos,
	})
}

type UpdateColumnInput struct {
	ID       string
	Name     *string
	Position *int
}

func (s *ColumnService) UpdateColumn(ctx context.Context, input UpdateColumnInput) (domain.Column, error) {
	col, err := s.columnRepo.GetByID(ctx, input.ID)
	if err != nil {
		return domain.Column{}, domain.ErrNotFound
	}

	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			return domain.Column{}, domain.ErrInvalidInput
		}
		col.Name = strings.TrimSpace(*input.Name)
	}

	if input.Position != nil {
		col.Position = *input.Position
	}

	return s.columnRepo.Update(ctx, col)
}

func (s *ColumnService) DeleteColumn(ctx context.Context, columnID string) error {
	err := s.columnRepo.Delete(ctx, columnID)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.ErrNotFound
	}
	return err
}