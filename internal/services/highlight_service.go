package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
)

type highlightService interface {
	Create(ctx context.Context, highlight *models.Highlight) error
	GetHighlightById(ctx context.Context, id string) (*models.Highlight, error)
	GetAllHighlightByUserId(ctx context.Context, userId string, page, limit int) ([]*models.Highlight, error)
	UpdateHighlight(ctx context.Context, highlight *models.Highlight) error
	DeleteHighlightById(ctx context.Context, id string) error
}

type HighlightService struct {
	repo repository.HighlightRepository
}

func NewHighlightService(repo repository.HighlightRepository) *HighlightService {
	return &HighlightService{repo: repo}
}

func (h *HighlightService) Create(ctx context.Context, highlight *models.Highlight) error {
	return h.repo.Create(ctx, highlight)
}

func (h *HighlightService) GetHighlightById(ctx context.Context, id string) (*models.Highlight, error) {
	return h.repo.GetByID(ctx, id)
}

func (h *HighlightService) GetAllHighlightByUserId(ctx context.Context, userId string, page, limit int) ([]models.Highlight, error) {
	return h.repo.GetByUserId(ctx, userId, page, limit)
}

func (h *HighlightService) UpdateHighlight(ctx context.Context, highlight *models.Highlight) error {
	return h.repo.Update(ctx, highlight)
}

func (h *HighlightService) DeleteHighlightById(ctx context.Context, id string) error {
	return h.repo.Delete(ctx, id)
}
