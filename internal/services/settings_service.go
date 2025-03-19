package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
)

type settingsService interface {
	Create(ctx context.Context, settings *models.Settings) error
	GetById(ctx context.Context, settingsId string) (*models.Settings, error)
	GetByUserId(ctx context.Context, userId string) (*models.Settings, error)
	Update(ctx context.Context, settings *models.Settings) error
	Delete(ctx context.Context, settingsId string) error
}

type SettingsService struct {
	repo repository.SettingsRepository
}

func NewSettingsService(repo repository.SettingsRepository) *SettingsService {
	return &SettingsService{
		repo: repo,
	}
}

func (s *SettingsService) Create(ctx context.Context, settings *models.Settings) error {
	return s.repo.Create(ctx, settings)
}

func (s *SettingsService) GetById(ctx context.Context, settingsId string) (*models.Settings, error) {
	return s.repo.GetByID(ctx, settingsId)
}

func (s *SettingsService) GetByUserId(ctx context.Context, userId string) (*models.Settings, error) {
	return s.repo.GetByUserID(ctx, userId)
}

func (s *SettingsService) Update(ctx context.Context, settings *models.Settings) error {
	return s.repo.Update(ctx, settings)
}

func (s *SettingsService) Delete(ctx context.Context, settingsId string) error {
	return s.repo.Delete(ctx, settingsId)
}
