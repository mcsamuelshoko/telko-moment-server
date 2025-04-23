package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/rs/zerolog"
	"time"
)

type IAuthenticationService interface {
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string, tokenDuration time.Duration) error
	GetUserIDFromRefreshToken(ctx context.Context, refreshToken string) (string, error)
	UpdateUserRefreshToken(ctx context.Context, userID string, refreshToken string, tokenDuration time.Duration) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type AuthenticationService struct {
	iName string
	log   *zerolog.Logger
	repo  repository.IAuthenticationRepository
}

func NewAuthenticationService(log *zerolog.Logger, repo repository.IAuthenticationRepository) IAuthenticationService {
	return &AuthenticationService{
		iName: "AuthenticationService",
		log:   log,
		repo:  repo,
	}
}

func (a AuthenticationService) SaveRefreshToken(ctx context.Context, userID string, refreshToken string, tokenDuration time.Duration) error {
	const kName = "SaveRefreshToken"

	err := a.repo.SaveRefreshToken(ctx, userID, refreshToken, tokenDuration)
	if err != nil {
		a.log.Error().Interface(kName, a.iName).Err(err).Msg("Failed to save refresh token")
		return err
	}
	return nil
}

func (a AuthenticationService) GetUserIDFromRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	const kName = "GetUserIDFromRefreshToken"

	userID, err := a.repo.GetUserIDFromRefreshToken(ctx, refreshToken)
	if err != nil {
		a.log.Error().Interface(kName, a.iName).Err(err).Msg("Failed to get user ID from refresh token")
		return "", err
	}
	return userID, nil
}

func (a AuthenticationService) UpdateUserRefreshToken(ctx context.Context, userID string, refreshToken string, tokenDuration time.Duration) error {
	const kName = "UpdateUserRefreshToken"

	err := a.repo.SaveRefreshToken(ctx, userID, refreshToken, tokenDuration)
	if err != nil {
		a.log.Error().Interface(kName, a.iName).Err(err).Msg("Failed to update refresh token")
		return err
	}
	return nil
}

func (a AuthenticationService) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	const kName = "DeleteRefreshToken"

	err := a.repo.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		a.log.Error().Interface(kName, a.iName).Err(err).Msg("Failed to delete refresh token")
		return err
	}
	return nil
}
