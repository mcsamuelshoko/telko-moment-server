package configs

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type JwtConfig struct {
	Issuer                     string `env:"JWT_ISSUER" envDefault:"telko_moment"`
	Secret                     string `env:"JWT_SECRET" envDefault:"secret"`
	TokenDuration              string `env:"JWT_TOKEN_DURATION" envDefault:"1h"`
	RefreshTokenSecret         string `env:"JWT_REFRESH_TOKEN_SECRET" envDefault:"secret"`
	RefreshTokenDuration       string `env:"JWT_REFRESH_TOKEN_DURATION" envDefault:"24h"`
	RefreshTokenDaysMultiplier string `env:"JWT_REFRESH_TOKEN_DAYS_MULTIPLIER" envDefault:"24h"`
}
type Config struct {
	MongoDB struct {
		URI      string `env:"MONGODB_URI" envDefault:"mongodb://localhost:27017"`
		Database string `env:"MONGODB_DB" envDefault:"tel_mont_db"`
	} `json:"mongodb"`
	Server struct {
		Port string `env:"SERVER_PORT" envDefault:":8080"`
	} `json:"server"`
	Jwt        JwtConfig `json:"jwt"`
	Encryption struct {
		AESKey string `env:"ENC_AES_KEY" envDefault:"0123456789abcdef"`
	} `json:"encryption"`
	Hashing struct {
		HMACSecretKey string `env:"HMAC_SECRET_KEY" envDefault:"0123456789abcdef"`
	} `json:"hashing"`
}

func LoadConfig() (*Config, error) {
	// Determine the root directory of the project
	rootDir, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get working directory")
		return nil, err
	}

	// Load environment variables from .env file
	envPath := filepath.Join(rootDir, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Info().Msg("No .env file found or unable to load, using environment variables")
	}

	var config Config

	// Load the ones left by -> envconfig.Process()
	config.MongoDB.Database = os.Getenv("MONGODB_DB")
	config.MongoDB.URI = os.Getenv("MONGODB_URI")

	config.Server.Port = os.Getenv("SERVER_PORT")

	config.Jwt.Issuer = os.Getenv("JWT_ISSUER")
	config.Jwt.Secret = os.Getenv("JWT_SECRET")
	config.Jwt.TokenDuration = os.Getenv("JWT_TOKEN_DURATION")
	config.Jwt.RefreshTokenSecret = os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	config.Jwt.RefreshTokenDuration = os.Getenv("JWT_REFRESH_TOKEN_DURATION")
	config.Jwt.RefreshTokenDaysMultiplier = os.Getenv("JWT_REFRESH_TOKEN_DAYS_MULTIPLIER")

	config.Encryption.AESKey = os.Getenv("ENC_AES_KEY")

	config.Hashing.HMACSecretKey = os.Getenv("HMAC_SECRET_KEY")

	// Return the loaded configuration
	return &config, nil
}
