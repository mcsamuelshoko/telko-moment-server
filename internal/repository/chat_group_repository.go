package repository

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
)

type ChatGroupRepository interface {
	Create(ctx context.Context, chatGroup *models.ChatGroup) error
	GetByID(ctx context.Context, id string) (*models.ChatGroup, error)
	List(ctx context.Context, page, limit int) ([]models.ChatGroup, error)
	Update(ctx context.Context, chatGroup *models.ChatGroup) error
	Delete(ctx context.Context, id string) error
}
