package service

import (
	"context"
	"errors"
	"strings"

	"github.com/cledupe/go-kanban/backend/internal/domain"
)

type CardService struct {
	columnRepo domain.ColumnRepository
	cardRepo   domain.CardRepository
}

func NewCardService(
	columnRepo domain.ColumnRepository,
	cardRepo domain.CardRepository,
) *CardService {
	return &CardService{
		columnRepo: columnRepo,
		cardRepo:   cardRepo,
	}
}

type CreateCardInput struct {
	ColumnID    string
	Title       string
	Description string
}

func (s *CardService) CreateCard(ctx context.Context, input CreateCardInput) (domain.Card, error) {
	if strings.TrimSpace(input.Title) == "" {
		return domain.Card{}, domain.ErrInvalidInput
	}

	if _, err := s.columnRepo.GetByID(ctx, input.ColumnID); err != nil {
		return domain.Card{}, domain.ErrNotFound
	}

	cards, err := s.cardRepo.ListByColumnID(ctx, input.ColumnID)
	if err != nil {
		return domain.Card{}, err
	}

	nextPos := 0
	for _, c := range cards {
		if c.Position >= nextPos {
			nextPos = c.Position + 1
		}
	}

	return s.cardRepo.Create(ctx, domain.Card{
		ColumnID:    input.ColumnID,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Position:    nextPos,
	})
}

type UpdateCardInput struct {
	ID          string
	Title       *string
	Description *string
}

func (s *CardService) UpdateCard(ctx context.Context, input UpdateCardInput) (domain.Card, error) {
	card, err := s.cardRepo.GetByID(ctx, input.ID)
	if err != nil {
		return domain.Card{}, domain.ErrNotFound
	}

	if input.Title != nil {
		if strings.TrimSpace(*input.Title) == "" {
			return domain.Card{}, domain.ErrInvalidInput
		}
		card.Title = strings.TrimSpace(*input.Title)
	}

	if input.Description != nil {
		card.Description = *input.Description
	}

	return s.cardRepo.Update(ctx, card)
}

type MoveCardInput struct {
	CardID         string
	TargetColumnID string
	Position       int
}

func (s *CardService) MoveCard(ctx context.Context, input MoveCardInput) error {
	if _, err := s.cardRepo.GetByID(ctx, input.CardID); err != nil {
		return domain.ErrNotFound
	}

	if _, err := s.columnRepo.GetByID(ctx, input.TargetColumnID); err != nil {
		return domain.ErrNotFound
	}

	if input.Position < 0 {
		return domain.ErrInvalidInput
	}

	return s.cardRepo.Move(ctx, input.CardID, input.TargetColumnID, input.Position)
}

func (s *CardService) DeleteCard(ctx context.Context, cardID string) error {
	err := s.cardRepo.Delete(ctx, cardID)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.ErrNotFound
	}
	return err
}