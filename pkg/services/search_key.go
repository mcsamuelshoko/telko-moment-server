package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/rs/zerolog"
	"strings"
)

// ISearchKeyService defines the interface for generating deterministic,
// keyed hashes suitable for searchable fields.
type ISearchKeyService interface {
	// GenerateSearchKey creates a hex-encoded HMAC-SHA256 hash of the input string.
	// Ensures consistent standardization (lowercase, trimmed).
	GenerateSearchKey(input string) (string, error)
}

// hmacSearchKeyService implements ISearchKeyService using HMAC-SHA256.
type hmacSearchKeyService struct {
	secretKey []byte          // The secret key for the HMAC function
	log       *zerolog.Logger // Instance of the logger
}

// NewHMACSearchKeyService creates a new instance of the HMAC search key service.
// The secretKeyHex is the hex-encoded string of your securely stored secret key.
// It enforces a minimum key length for security (>= 32 bytes recommended for HMAC-SHA256).
func NewHMACSearchKeyService(log *zerolog.Logger, secretKeyHex string) (ISearchKeyService, error) {
	keyBytes, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode HMAC secret key from hex")
		// Avoid logging the key itself here
		return nil, errors.New("invalid HMAC secret key format")
	}

	// NIST recommends the key length >= hash output size. SHA256 output is 32 bytes.
	const minKeyLength = 32
	if len(keyBytes) < minKeyLength {
		log.Error().Int("key_length", len(keyBytes)).Int("minimum_required", minKeyLength).Msg("HMAC secret key is too short")
		return nil, errors.New("HMAC secret key provided is too short for security requirements")
	}

	// You might want to create a sub-logger specific to this service
	serviceLogger := log.With().Str("service", "HMACSearchKeyService").Logger()

	return &hmacSearchKeyService{
		secretKey: keyBytes,
		log:       &serviceLogger, // Use the sub-logger
	}, nil
}

// GenerateSearchKey implements the ISearchKeyService interface.
// It standardizes the input string (lowercase, trim whitespace) and then computes
// the HMAC-SHA256, returning it as a hex-encoded string.
func (s *hmacSearchKeyService) GenerateSearchKey(input string) (string, error) {
	// Step 1: Consistently standardize the input. Crucial for deterministic results.
	standardizedInput := strings.ToLower(strings.TrimSpace(input))

	// Step 2: Compute the HMAC-SHA256
	mac := hmac.New(sha256.New, s.secretKey)
	// The Write method on hash.Hash implementations (like HMAC) doesn't return an error.
	_, _ = mac.Write([]byte(standardizedInput))
	searchKeyBytes := mac.Sum(nil)

	// Step 3: Encode the resulting hash bytes to a hex string for storage/lookup.
	searchKeyHex := hex.EncodeToString(searchKeyBytes)

	// Note: HMAC generation itself with valid inputs doesn't typically error here.
	// Errors are handled during initialization (key validation).
	s.log.Debug().Msg("Successfully generated search key") // Optional debug log

	return searchKeyHex, nil
}
