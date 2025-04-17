package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mcsamuelshoko/telko-moment-server/configs"
	"github.com/rs/zerolog"
	"strconv"
	"time"
)

// IJWTService handles encryption/decryption (define this in pkg/utils or a dedicated service)
type IJWTService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	VerifyAccessToken(tokenString string) (*jwt.Token, error)
	VerifyRefreshToken(tokenString string) bool
	GetRefreshTokenDuration() time.Duration
}

type JWTService struct {
	issuer                  string
	jwtSecret               []byte // In production, fetch this from a secure KMS
	jwtTokenDuration        time.Duration
	jwtRefreshTokenSecret   []byte
	jwtRefreshTokenDuration time.Duration
	log                     *zerolog.Logger
}

func NewJWTService(logger *zerolog.Logger, cfg configs.JwtConfig) (IJWTService, error) {
	secret, err := hex.DecodeString(cfg.Secret)
	if err != nil || len(secret) < 32 {
		logger.Error().Err(err).Msg("Invalid JWT Secret")
		return nil, errors.New("invalid JWT secret key: must be a hex-encoded string of at least 32 bytes")
	}

	duration, err := time.ParseDuration(cfg.TokenDuration)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid JWT Token Duration")
		return nil, errors.New("invalid JWT token duration string")
	}

	refreshSecret, err := hex.DecodeString(cfg.RefreshTokenSecret)
	if err != nil || len(refreshSecret) < 32 {
		logger.Error().Err(err).Msg("Invalid JWT Refresh Secret")
		return nil, errors.New("invalid refresh JWT secret key: must be a hex-encoded string of at least 32 bytes")
	}

	refreshDuration, err := time.ParseDuration(cfg.RefreshTokenDuration)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid JWT Refresh Duration")
		return nil, errors.New("invalid JWT refresh-token duration string")
	}

	refreshDaysMultiplier, err := strconv.Atoi(cfg.RefreshTokenDaysMultiplier)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid JWT Refresh Days-Multiplier")
		return nil, errors.New("invalid JWT refresh-token days-multiplier duration string")
	}

	return &JWTService{
		issuer:                  cfg.Issuer,
		jwtSecret:               secret,
		jwtRefreshTokenSecret:   refreshSecret,
		jwtTokenDuration:        duration,
		jwtRefreshTokenDuration: refreshDuration * time.Duration(refreshDaysMultiplier),
		log:                     logger,
	}, nil
}

func (j *JWTService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.jwtTokenDuration).Unix(), // Access token expiration
		"nbf": time.Now().Unix(),                         // Not Before, is the time it should be allowed use after it hs passed
		"iat": time.Now().Unix(),                         // issuedAt time
		"iss": j.issuer,                                  // Token Issuer
		"jti": uuid.New().String(),                       // Nonce , can be used token ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}

func (j *JWTService) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.jwtRefreshTokenDuration).Unix(), // Refresh token expiration
		"nbf": time.Now().Unix(),                                // Not Before, is the time it should be allowed use after it hs passed
		"iat": time.Now().Unix(),                                // issuedAt time
		"iss": j.issuer,                                         // Token Issuer
		"jti": uuid.New().String(),                              // JWT token ID
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

func (j *JWTService) GetRefreshTokenDuration() time.Duration {
	return j.jwtRefreshTokenDuration
}
