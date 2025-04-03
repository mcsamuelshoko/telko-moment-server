package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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
	jwtTokenDuration      string
	log                   *zerolog.Logger
}

func NewJWTService(logger *zerolog.Logger, secret, refreshSecret []byte, duration string) IJWTService {
	return &JWTService{
		log:                   logger,
		jwtSecret:             secret,
		jwtRefreshTokenSecret: refreshSecret,
		jwtTokenDuration:      duration,
	}
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
