package repository

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"time"
)

type IAuthenticationRepository interface {
	Create(ctx context.Context, auth *models.Authentication) (*models.Authentication, error)
	GetList(ctx context.Context) (*[]models.Authentication, error)
	GetByUserID(ctx context.Context, userID string) (*models.Authentication, error)
	UpdateByUserID(ctx context.Context, userID string, auth *models.Authentication) (*models.Authentication, error)

	Delete(ctx context.Context, ID string) error
	DeleteByUserID(ctx context.Context, userID string) error

	SaveRefreshToken(ctx context.Context, userID string, refreshToken string, tokenDuration time.Duration) error
	GetUserIDFromRefreshToken(ctx context.Context, refreshToken string) (string, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}
