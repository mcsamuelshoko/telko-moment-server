package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type JwtConfig struct {
	Secret               string `env:"JWT_SECRET" envDefault:"secret"`
	RefreshTokenSecret   string `env:"JWT_REFRESH_TOKEN_SECRET" envDefault:"secret"`
	RefreshTokenDuration string `env:"JWT_REFRESH_TOKEN_DURATION" envDefault:"1h"`
}
type Config struct {
	MongoDB struct {
		URI      string `env:"MONGODB_URI" envDefault:"mongodb://localhost:27017"`
		Database string `env:"MONGODB_DB" envDefault:"telko_moment_db"`
	} `json:"mongodb"`
	Server struct {
		Port string `env:"SERVER_PORT" envDefault:":8080"`
	} `json:"server"`
	Jwt        JwtConfig `json:"jwt"`
	Encryption struct {
		AESKey string `env:"ENC_AES_KEY" envDefault:"0123456789abcdef"`
	} `json:"encryption"`
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

	// Load configuration from environment variables or default values
	err = envconfig.Process("", &config)
	if err != nil {
		log.Error().Err(err).Msg("Error loading configuration")
		return nil, err
	}
	// Load the ones left by -> envconfig.Process()
	config.Jwt.RefreshTokenSecret = os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	config.Jwt.RefreshTokenDuration = os.Getenv("JWT_REFRESH_TOKEN_DURATION")

	config.Encryption.AESKey = os.Getenv("ENC_AES_KEY")

	// Return the loaded configuration
	return &config, nil
}
