package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cledupe/go-kanban/backend/internal/domain"
)

func TestCardServiceCreateCard(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{ID: id, Name: "To Do"}, nil
		},
	}

	cardRepo := &mockCardRepo{
		listByColumnFn: func(_ context.Context, columnID string) ([]domain.Card, error) {
			return nil, nil
		},
		createFn: func(_ context.Context, card domain.Card) (domain.Card, error) {
			card.ID = "card-1"
			return card, nil
		},
	}

	s := NewCardService(columnRepo, cardRepo)

	card, err := s.CreateCard(context.Background(), CreateCardInput{
		ColumnID:    "col-1",
		Title:       "My Card",
		Description: "A description",
	})
	if err != nil {
		t.Fatalf("create card: %v", err)
	}
	if card.Title != "My Card" {
		t.Fatalf("expected title 'My Card', got %q", card.Title)
	}
	if card.Description != "A description" {
		t.Fatalf("expected description 'A description', got %q", card.Description)
	}
}

func TestCardServiceCreateCardRejectsEmptyTitle(t *testing.T) {
	t.Parallel()

	s := NewCardService(&mockColumnRepo{}, &mockCardRepo{})

	_, err := s.CreateCard(context.Background(), CreateCardInput{ColumnID: "c", Title: ""})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCardServiceCreateCardColumnNotFound(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{}, domain.ErrNotFound
		},
	}

	s := NewCardService(columnRepo, &mockCardRepo{})

	_, err := s.CreateCard(context.Background(), CreateCardInput{
		ColumnID: "nonexistent",
		Title:    "Card",
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCardServiceUpdateCard(t *testing.T) {
	t.Parallel()

	cardRepo := &mockCardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Card, error) {
			return domain.Card{ID: id, Title: "Old", Description: "Old desc"}, nil
		},
		updateFn: func(_ context.Context, card domain.Card) (domain.Card, error) {
			return card, nil
		},
	}

	s := NewCardService(&mockColumnRepo{}, cardRepo)

	title := "Updated"
	desc := "New desc"
	card, err := s.UpdateCard(context.Background(), UpdateCardInput{
		ID:          "card-1",
		Title:       &title,
		Description: &desc,
	})
	if err != nil {
		t.Fatalf("update card: %v", err)
	}
	if card.Title != "Updated" {
		t.Fatalf("expected title 'Updated', got %q", card.Title)
	}
	if card.Description != "New desc" {
		t.Fatalf("expected description 'New desc', got %q", card.Description)
	}
}

func TestCardServiceUpdateCardRejectsEmptyTitle(t *testing.T) {
	t.Parallel()

	cardRepo := &mockCardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Card, error) {
			return domain.Card{ID: id, Title: "Card"}, nil
		},
	}

	s := NewCardService(&mockColumnRepo{}, cardRepo)

	title := ""
	_, err := s.UpdateCard(context.Background(), UpdateCardInput{
		ID:    "card-1",
		Title: &title,
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCardServiceUpdateCardNotFound(t *testing.T) {
	t.Parallel()

	cardRepo := &mockCardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Card, error) {
			return domain.Card{}, domain.ErrNotFound
		},
	}

	s := NewCardService(&mockColumnRepo{}, cardRepo)

	_, err := s.UpdateCard(context.Background(), UpdateCardInput{ID: "nonexistent"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCardServiceMoveCard(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{ID: id, Name: "Column"}, nil
		},
	}

	cardRepo := &mockCardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Card, error) {
			return domain.Card{ID: id, Title: "Card", Position: 0}, nil
		},
		moveFn: func(_ context.Context, cardID string, targetColumnID string, position int) error {
			return nil
		},
	}

	s := NewCardService(columnRepo, cardRepo)

	err := s.MoveCard(context.Background(), MoveCardInput{
		CardID:         "card-1",
		TargetColumnID: "col-2",
		Position:       0,
	})
	if err != nil {
		t.Fatalf("move card: %v", err)
	}
}

func TestCardServiceMoveCardCardNotFound(t *testing.T) {
	t.Parallel()

	cardRepo := &mockCardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Card, error) {
			return domain.Card{}, domain.ErrNotFound
		},
	}

	s := NewCardService(&mockColumnRepo{}, cardRepo)

	err := s.MoveCard(context.Background(), MoveCardInput{
		CardID:         "nonexistent",
		TargetColumnID: "col-2",
		Position:       0,
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCardServiceMoveCardColumnNotFound(t *testing.T) {
	t.Parallel()

	cardRepo := &mockCardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Card, error) {
			return domain.Card{ID: id, Title: "Card"}, nil
		},
	}

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{}, domain.ErrNotFound
		},
	}

	s := NewCardService(columnRepo, cardRepo)

	err := s.MoveCard(context.Background(), MoveCardInput{
		CardID:         "card-1",
		TargetColumnID: "nonexistent",
		Position:       0,
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCardServiceMoveCardRejectsNegativePosition(t *testing.T) {
	t.Parallel()

	cardRepo := &mockCardRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Card, error) {
			return domain.Card{ID: id, Title: "Card"}, nil
		},
	}

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{ID: id, Name: "Column"}, nil
		},
	}

	s := NewCardService(columnRepo, cardRepo)

	err := s.MoveCard(context.Background(), MoveCardInput{
		CardID:         "card-1",
		TargetColumnID: "col-2",
		Position:       -1,
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCardServiceDeleteCard(t *testing.T) {
	t.Parallel()

	cardRepo := &mockCardRepo{
		deleteFn: func(_ context.Context, id string) error {
			return nil
		},
	}

	s := NewCardService(&mockColumnRepo{}, cardRepo)

	err := s.DeleteCard(context.Background(), "card-1")
	if err != nil {
		t.Fatalf("delete card: %v", err)
	}
}

func TestCardServiceDeleteCardNotFound(t *testing.T) {
	t.Parallel()

	cardRepo := &mockCardRepo{
		deleteFn: func(_ context.Context, id string) error {
			return domain.ErrNotFound
		},
	}

	s := NewCardService(&mockColumnRepo{}, cardRepo)

	err := s.DeleteCard(context.Background(), "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCardServiceAssignsNextPosition(t *testing.T) {
	t.Parallel()

	columnRepo := &mockColumnRepo{
		getByIDFn: func(_ context.Context, id string) (domain.Column, error) {
			return domain.Column{ID: id}, nil
		},
	}

	cardRepo := &mockCardRepo{
		listByColumnFn: func(_ context.Context, columnID string) ([]domain.Card, error) {
			return []domain.Card{
				{Position: 0},
				{Position: 1},
				{Position: 5},
			}, nil
		},
		createFn: func(_ context.Context, card domain.Card) (domain.Card, error) {
			if card.Position != 6 {
				t.Fatalf("expected position 6, got %d", card.Position)
			}
			return card, nil
		},
	}

	s := NewCardService(columnRepo, cardRepo)

	_, err := s.CreateCard(context.Background(), CreateCardInput{
		ColumnID: "col-1",
		Title:    "New Card",
	})
	if err != nil {
		t.Fatalf("create card: %v", err)
	}
}