package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
)

type messageService interface {
	Create(ctx context.Context, message *models.Message) error
	Update(ctx context.Context, message *models.Message) error
	GetById(ctx context.Context, messageId string) (*models.Message, error)
	GetBySenderId(ctx context.Context, userId string) (*models.Message, error)
	GetByChatId(ctx context.Context, chatId string) (*models.Message, error)
	Delete(ctx context.Context, messageId string) error
}

type MessageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (m *MessageService) Create(ctx context.Context, message *models.Message) error {
	return m.repo.Create(ctx, message)
}

func (m *MessageService) Update(ctx context.Context, message *models.Message) error {
	return m.repo.Update(ctx, message)
}

func (m *MessageService) GetById(ctx context.Context, messageId string) (*models.Message, error) {
	return m.repo.GetByID(ctx, messageId)
}

func (m *MessageService) GetBySenderId(ctx context.Context, userId string) (*models.Message, error) {
	return m.repo.GetBySenderID(ctx, userId)
}

func (m *MessageService) GetByChatId(ctx context.Context, chatId string) (*models.Message, error) {
	return m.repo.GetByChatID(ctx, chatId)
}

func (m *MessageService) Delete(ctx context.Context, messageId string) error {
	return m.repo.Delete(ctx, messageId)
}
