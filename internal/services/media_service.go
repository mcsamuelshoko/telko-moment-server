package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
)

type mediaService interface {
	CreateMedia(ctx context.Context, media *models.Media) error
	GetMediaById(ctx context.Context, id string) (*models.Media, error)
	GetAllMediaByChatId(ctx context.Context, chatId string, page, limit int) ([]models.Media, error)
	GetAllMediaBySenderId(ctx context.Context, senderId string, page, limit int) ([]models.Media, error)
	UpdateMedia(ctx context.Context, media *models.Media) error
	DeleteMedia(ctx context.Context, id string) error
	//DeleteMediaByChatId(ctx context.Context, chatId string) error
	//DeleteMediaBySenderId(ctx context.Context, chatId string) error
}

type MediaService struct {
	repo repository.MediaRepository
}

func NewMediaService(repo repository.MediaRepository) *MediaService {
	return &MediaService{
		repo: repo,
	}
}

func (m *MediaService) CreateMedia(ctx context.Context, media *models.Media) error {
	return m.repo.Create(ctx, media)
}

func (m *MediaService) GetMediaById(ctx context.Context, id string) (*models.Media, error) {
	return m.repo.GetByID(ctx, id)
}
func (m *MediaService) GetAllMediaByChatId(ctx context.Context, chatId string, page, limit int) ([]models.Media, error) {
	return m.repo.GetByChatId(ctx, chatId, page, limit)
}
func (m *MediaService) GetAllMediaBySenderId(ctx context.Context, senderId string, page, limit int) ([]models.Media, error) {
	return m.repo.GetBySenderId(ctx, senderId, page, limit)
}
func (m *MediaService) UpdateMedia(ctx context.Context, media *models.Media) error {
	return m.repo.Update(ctx, media)
}
func (m *MediaService) DeleteMedia(ctx context.Context, id string) error {
	return m.repo.Delete(ctx, id)
}
