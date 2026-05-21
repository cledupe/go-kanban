package service

import (
	"context"
	"errors"
	"strings"

	"github.com/cledupe/go-kanban/backend/internal/domain"
)

type BoardService struct {
	boardRepo  domain.BoardRepository
	columnRepo domain.ColumnRepository
	cardRepo   domain.CardRepository
}

func NewBoardService(
	boardRepo domain.BoardRepository,
	columnRepo domain.ColumnRepository,
	cardRepo domain.CardRepository,
) *BoardService {
	return &BoardService{
		boardRepo:  boardRepo,
		columnRepo: columnRepo,
		cardRepo:   cardRepo,
	}
}

type CreateBoardInput struct {
	Name string
}

func (s *BoardService) ListBoards(ctx context.Context) ([]domain.Board, error) {
	return s.boardRepo.List(ctx)
}

func (s *BoardService) CreateBoard(ctx context.Context, input CreateBoardInput) (domain.Board, error) {
	if strings.TrimSpace(input.Name) == "" {
		return domain.Board{}, domain.ErrInvalidInput
	}

	return s.boardRepo.Create(ctx, domain.Board{Name: strings.TrimSpace(input.Name)})
}

func (s *BoardService) GetBoard(ctx context.Context, boardID string) (domain.BoardDetail, error) {
	board, err := s.boardRepo.GetByID(ctx, boardID)
	if err != nil {
		return domain.BoardDetail{}, domain.ErrNotFound
	}

	columns, err := s.columnRepo.ListByBoardID(ctx, boardID)
	if err != nil {
		return domain.BoardDetail{}, err
	}

	detail := domain.BoardDetail{Board: board}
	for _, col := range columns {
		cards, err := s.cardRepo.ListByColumnID(ctx, col.ID)
		if err != nil {
			return domain.BoardDetail{}, err
		}
		detail.Columns = append(detail.Columns, domain.ColumnWithCards{
			Column: col,
			Cards:  cards,
		})
	}

	return detail, nil
}

func (s *BoardService) UpdateBoard(ctx context.Context, boardID string, name string) (domain.Board, error) {
	if strings.TrimSpace(name) == "" {
		return domain.Board{}, domain.ErrInvalidInput
	}

	board, err := s.boardRepo.GetByID(ctx, boardID)
	if err != nil {
		return domain.Board{}, domain.ErrNotFound
	}

	board.Name = strings.TrimSpace(name)
	return s.boardRepo.Update(ctx, board)
}

func (s *BoardService) DeleteBoard(ctx context.Context, boardID string) error {
	err := s.boardRepo.Delete(ctx, boardID)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.ErrNotFound
	}
	return err
}