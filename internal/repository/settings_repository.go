package repository

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
)

type SettingsRepository interface {
	Create(ctx context.Context, settings *models.Settings) (*models.Settings, error)
	GetByID(ctx context.Context, id string) (*models.Settings, error)
	GetByUserID(ctx context.Context, userId string) (*models.Settings, error)
	List(ctx context.Context, page, limit int) ([]models.Settings, error)
	Update(ctx context.Context, settings *models.Settings) error
	Delete(ctx context.Context, id string) error
}
