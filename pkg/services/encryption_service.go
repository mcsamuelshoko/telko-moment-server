package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/rs/zerolog"
	"io"
)

// IEncryptionService handles encryption/decryption (define this in pkg/utils or a dedicated service)
type IEncryptionService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type AESEncryptionService struct {
	key []byte // In production, fetch this from a secure KMS
	log *zerolog.Logger
}

func NewAESEncryptionService(key string, log *zerolog.Logger) (IEncryptionService, error) {
	keyBytes, err := hex.DecodeString(key) // Key should be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256
	if err != nil || len(keyBytes) < 16 {
		log.Err(err).Str("key", key).Msg("Failed to decode encryption key")
		return nil, errors.New("invalid encryption key")
	}
	return &AESEncryptionService{key: keyBytes, log: log}, nil
}

func (s *AESEncryptionService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to create AES cipher")
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		s.log.Error().Err(err).Msg("Failed to generate IV")
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	return hex.EncodeToString(ciphertext), nil
}

func (s *AESEncryptionService) Decrypt(ciphertext string) (string, error) {
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to decode ciphertext")
		return "", err
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to create AES cipher")
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		s.log.Error().Msg("ciphertext is too short, increase size")
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
