package repository

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
)

type HighlightRepository interface {
	Create(ctx context.Context, highlight *models.Highlight) error
	GetByID(ctx context.Context, id string) (*models.Highlight, error)
	List(ctx context.Context, page, limit int) ([]models.Highlight, error)
	Update(ctx context.Context, highlight *models.Highlight) error
	Delete(ctx context.Context, id string) error
}
