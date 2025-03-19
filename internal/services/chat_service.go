package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
)

type chatService interface {
	CreateChat(ctx context.Context, chat *models.Chat) error
	GetChat(ctx context.Context, chatId string) (*models.Chat, error)
	ListByUserId(ctx context.Context, userId string) ([]models.Chat, error)
	//UpdateChat(ctx context.Context, chat *models.Chat) error
	DeleteChat(ctx context.Context, chatId string) error
}

type ChatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) *ChatService {
	return &ChatService{
		repo: repo,
	}
}

func (c ChatService) CreateChat(ctx context.Context, chat *models.Chat) error {
	return c.repo.Create(ctx, chat)
}

func (c ChatService) GetChat(ctx context.Context, chatId string) (*models.Chat, error) {
	return c.repo.GetByID(ctx, chatId)
}

func (c ChatService) ListByUserId(ctx context.Context, userId string, page, limit int) ([]models.Chat, error) {
	return c.repo.ListByUserId(ctx, userId, page, limit)
}

//func (c ChatService) UpdateChat(ctx context.Context, chat *models.Chat) error {
//	panic("implement me")
//}

func (c ChatService) DeleteChat(ctx context.Context, chatId string) error {
	return c.repo.Delete(ctx, chatId)
}
