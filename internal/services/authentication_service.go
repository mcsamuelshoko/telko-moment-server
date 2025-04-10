package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/rs/zerolog"
)

type IAuthenticationService interface {
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error
	GetUserIDFromRefreshToken(ctx context.Context, refreshToken string) (string, error)
	UpdateUserRefreshToken(ctx context.Context, userID string, refreshToken string) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type AuthenticationService struct {
	log  *zerolog.Logger
	repo repository.IAuthenticationRepository
}

func NewAuthenticationService(log *zerolog.Logger, repo repository.IAuthenticationRepository) IAuthenticationService {
	return &AuthenticationService{
		log:  log,
		repo: repo,
	}
}

func (a AuthenticationService) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	err := a.repo.SaveRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		a.log.Error().Err(err).Str("userID", userID).Msg("Failed to save refresh token")
		return err
	}
	return nil
}

func (a AuthenticationService) GetUserIDFromRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	userID, err := a.repo.GetUserIDFromRefreshToken(ctx, refreshToken)
	if err != nil {
		a.log.Error().Err(err).Str("refreshToken", refreshToken).Msg("Failed to get user ID from refresh token")
		return "", err
	}
	return userID, nil
}

func (a AuthenticationService) UpdateUserRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	err := a.repo.SaveRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		a.log.Error().Err(err).Str("userID", userID).Msg("Failed to update refresh token")
		return err
	}
	return nil
}

func (a AuthenticationService) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	err := a.repo.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		a.log.Error().Err(err).Str("refreshToken", refreshToken).Msg("Failed to delete refresh token")
		return err
	}
	return nil
}
