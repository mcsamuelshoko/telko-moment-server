package repository

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) error
	GetByID(ctx context.Context, id string) (*models.Chat, error)
	List(ctx context.Context, page, limit int) ([]models.Chat, error)
	Update(ctx context.Context, chat *models.Chat) error
	Delete(ctx context.Context, id string) error
}
