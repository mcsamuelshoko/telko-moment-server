package repository

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
)

type MediaRepository interface {
	Create(ctx context.Context, media *models.Media) error
	GetByID(ctx context.Context, id string) (*models.Media, error)
	List(ctx context.Context, page, limit int) ([]models.Media, error)
	Update(ctx context.Context, media *models.Media) error
	Delete(ctx context.Context, id string) error
}
