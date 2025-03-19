package repository

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id string) (*models.Message, error)
	GetByChatID(ctx context.Context, chatId string) (*models.Message, error)
	GetBySenderID(ctx context.Context, userId string) (*models.Message, error)
	List(ctx context.Context, page, limit int) ([]models.Message, error)
	Update(ctx context.Context, message *models.Message) error
	Delete(ctx context.Context, id string) error
}
