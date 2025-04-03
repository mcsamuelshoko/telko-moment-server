package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mcsamuelshoko/telko-moment-server/configs"
	"github.com/rs/zerolog"
	"time"
)

// IJWTService handles encryption/decryption (define this in pkg/utils or a dedicated service)
type IJWTService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	VerifyAccessToken(tokenString string) (*jwt.Token, error)
	VerifyRefreshToken(tokenString string) bool
}

type JWTService struct {
	jwtSecret             []byte // In production, fetch this from a secure KMS
	jwtRefreshTokenSecret []byte
	jwtTokenDuration      time.Duration
	log                   *zerolog.Logger
}

func NewJWTService(logger *zerolog.Logger, cfg configs.JwtConfig) (*JWTService, error) {
	secret, err := hex.DecodeString(cfg.Secret)
	if err != nil || len(secret) < 32 {
		logger.Error().Err(err).Msg("Invalid JWT Secret")
		return nil, errors.New("invalid JWT secret key: must be a hex-encoded string of at least 32 bytes")
	}

	refreshSecret, err := hex.DecodeString(cfg.RefreshTokenSecret)
	if err != nil || len(refreshSecret) < 32 {
		logger.Error().Err(err).Msg("Invalid JWT Refresh Secret")
		return nil, errors.New("invalid refresh JWT secret key: must be a hex-encoded string of at least 32 bytes")
	}

	duration, err := time.ParseDuration(cfg.RefreshTokenDuration)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid JWT Refresh Duration")
		return nil, errors.New("invalid JWT duration string")
	}

	return &JWTService{
		jwtSecret:             secret,
		jwtRefreshTokenSecret: refreshSecret,
		jwtTokenDuration:      duration,
		log:                   logger,
	}, nil
}

func (j *JWTService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(), // Access token expires in 15 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}

func (j *JWTService) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtRefreshTokenSecret)
}

func (j *JWTService) VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			j.log.Error().Msgf("Failed to validate token, Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.jwtSecret, nil
	})

	if err != nil {
		j.log.Error().Msgf("Failed to validate token: %v", err)
		return nil, err
	}

	return token, nil
}

func (j *JWTService) VerifyRefreshToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			j.log.Error().Msgf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.jwtRefreshTokenSecret, nil
	})

	if err != nil || !token.Valid {
		j.log.Error().Err(err).Msg("Failed to verify refresh token")
		return false
	}
	return true
}
