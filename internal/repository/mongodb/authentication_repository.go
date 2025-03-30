package mongodb

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthenticationRepository struct {
	Collection        *mongo.Collection
	Logger            *zerolog.Logger
	EncryptionService utils.IEncryptionService
}

func NewAuthenticationRepository(log *zerolog.Logger, db *mongo.Database, encryptSvc utils.IEncryptionService) repository.IAuthenticationRepository {
	return &AuthenticationRepository{
		Collection:        db.Collection("authentications"),
		Logger:            log,
		EncryptionService: encryptSvc,
	}
}

func (a AuthenticationRepository) Create(ctx context.Context, auth *models.Authentication) (*models.Authentication, error) {
	// Encrypt sensitive fields before saving
	if err := auth.EncryptFields(a.EncryptionService); err != nil {
		a.Logger.Error().Err(err).Msg("Failed to encrypt fields in AuthenticationRepository.Create")
		return nil, err
	}

	result, err := a.Collection.InsertOne(ctx, auth)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to create authentication record")
		return nil, err
	}
	auth.ID = result.InsertedID.(primitive.ObjectID)
	return auth, nil
}

func (a AuthenticationRepository) GetList(ctx context.Context) (*[]models.Authentication, error) {
	cursor, err := a.Collection.Find(ctx, bson.M{})
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to get authentication list")
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			a.Logger.Error().Err(err).Msg("Failed to close cursor")
		}
	}(cursor, ctx)

	var authList []models.Authentication
	var decryptedAuthList []models.Authentication
	if err := cursor.All(ctx, &authList); err != nil {
		a.Logger.Error().Err(err).Msg("Failed to decode authentication list")
		return nil, err
	}

	//Decrypt List before sharing to user
	for _, authentication := range authList {
		// Encrypt sensitive fields before saving
		if err := authentication.EncryptFields(a.EncryptionService); err != nil {
			a.Logger.Error().Err(err).Msg("Failed to encrypt fields in AuthenticationRepository.GetList")
			return nil, err
		}
		decryptedAuthList = append(decryptedAuthList, authentication)
	}

	return &decryptedAuthList, nil
}

func (a AuthenticationRepository) GetByUserID(ctx context.Context, userID string) (*models.Authentication, error) {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to convert Auth-GetByUserid to object id")
		a.Logger.Debug().Err(err).Msg("Failed to convert auth-user-id:" + userID)
		return nil, err
	}
	var auth models.Authentication
	err = a.Collection.FindOne(ctx, bson.M{"userId": ID}).Decode(&auth)
	if err != nil {
		a.Logger.Error().Err(err).Str("userID", userID).Msg("Failed to get authentication by user ID")
		return nil, err
	}
	// Decrypt sensitive fields before sharing
	if err := auth.DecryptFields(a.EncryptionService); err != nil {
		a.Logger.Error().Err(err).Msg("Failed to decrypt fields in AuthenticationRepository.GetByUserID")
		return nil, err
	}
	return &auth, nil
}

func (a AuthenticationRepository) UpdateByUserID(ctx context.Context, userID string, auth *models.Authentication) (*models.Authentication, error) {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to convert id to object id")
		a.Logger.Debug().Err(err).Msg("Failed to convert id:" + userID)
		return nil, err
	}
	// Encrypt sensitive fields before saving
	if err := auth.EncryptFields(a.EncryptionService); err != nil {
		a.Logger.Error().Err(err).Msg("Failed to encrypt fields in AuthenticationRepository.UpdateByUserID")
		return nil, err
	}

	filter := bson.M{"userId": ID}
	update := bson.M{"$set": auth}
	_, err = a.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		a.Logger.Error().Err(err).Str("userID", userID).Msg("Failed to update authentication by user ID")
		return nil, err
	}
	return auth, nil
}

func (a AuthenticationRepository) Delete(ctx context.Context, ID string) error {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		a.Logger.Error().Err(err).Str("ID", ID).Msg("Invalid ID format")
		return err
	}
	_, err = a.Collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		a.Logger.Error().Err(err).Str("ID", ID).Msg("Failed to delete authentication")
		return err
	}
	return nil
}

func (a AuthenticationRepository) DeleteByUserID(ctx context.Context, userID string) error {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to convert id to auth-object id")
		a.Logger.Debug().Err(err).Msg("Failed to convert auth-user-id:" + userID)
		return err
	}
	_, err = a.Collection.DeleteOne(ctx, bson.M{"userId": ID})
	if err != nil {
		a.Logger.Error().Err(err).Str("userID", userID).Msg("Failed to delete authentication by user ID")
		return err
	}
	return nil
}

func (a AuthenticationRepository) GenerateUniqueRefreshToken(ctx context.Context, userID string) (string, error) {
	const maxAttempts = 5

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Generate a secure random token
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			a.Logger.Error().Err(err).Msg("Failed to generate refresh token")
			return "", err
		}
		refreshToken := hex.EncodeToString(tokenBytes)
		// Encrypt token for saving to db
		encryptedRefreshToken, err := a.EncryptionService.Encrypt(refreshToken)
		if err != nil {
			a.Logger.Error().Err(err).Msg("Failed to encrypt refresh token in AuthenticationRepository.GenerateUniqueRefreshToken")
			return "", err
		}

		// Try to insert the new token
		_, err = a.Collection.InsertOne(ctx, bson.M{
			"userId":       userID,
			"refreshToken": encryptedRefreshToken,
			"expiresAt":    time.Now().Add(7 * 24 * time.Hour), // Example: 7 days expiration
			"isActive":     true,
			"createdAt":    time.Now(),
		})

		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				// Token already exists, try again
				continue
			}
			a.Logger.Error().Err(err).Str("userID", userID).Msg("Failed to store refresh token")
			return "", err
		}

		return refreshToken, nil
	}

	return "", fmt.Errorf("failed to generate unique refresh token after %d attempts", maxAttempts)
}

func (a AuthenticationRepository) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to convert id to object id")
		a.Logger.Debug().Err(err).Msg("Failed to convert id:" + userID)
		return err
	}
	// Encrypt Refresh Token before saving
	refreshToken, err = a.EncryptionService.Encrypt(refreshToken)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to encrypt refresh token in AuthenticationRepository.SaveRefreshToken")
		return err
	}

	_, err = a.Collection.InsertOne(ctx, bson.M{"userId": ID, "token": refreshToken, "created_at": time.Now()})
	if err != nil {
		a.Logger.Error().Err(err).Str("userID", userID).Msg("Failed to save refresh token")
		return err
	}
	return nil
}

func (a AuthenticationRepository) GetUserIDFromRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Input validation
	if refreshToken == "" {
		return "", fmt.Errorf("refresh token cannot be empty")
	}

	// Encrypt token for search
	encryptedRefreshToken, err := a.EncryptionService.Encrypt(refreshToken)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to encrypt refresh token")
		return "", err
	}
	var result models.Authentication

	// Find the token document
	err = a.Collection.FindOne(ctx, bson.M{
		"refreshToken": encryptedRefreshToken,
		"is_active":    true, // Only match active tokens
		"expires_at": bson.M{
			"$gt": time.Now(), // Only match non-expired tokens
		},
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.Logger.Warn().Str("refreshToken", refreshToken).Msg("No active refresh token found")
			return "", fmt.Errorf("invalid or expired refresh token")
		}
		a.Logger.Error().Err(err).Str("refreshToken", refreshToken).Msg("Failed to get user ID from refresh token")
		return "", err
	}

	if result.UserID.String() == "" {
		a.Logger.Warn().Str("refreshToken", refreshToken).Msg("Refresh token found but user ID is empty")
		return "", fmt.Errorf("invalid refresh token: no user associated")
	}

	return result.UserID.String(), nil
}

func (a AuthenticationRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	// Encrypt token for search
	encryptedRefreshToken, err := a.EncryptionService.Encrypt(refreshToken)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Failed to encrypt refresh token")
		return err
	}
	_, err = a.Collection.DeleteOne(ctx, bson.M{"token": encryptedRefreshToken})
	if err != nil {
		a.Logger.Error().Err(err).Str("refreshToken", refreshToken).Msg("Failed to delete refresh token")
		return err
	}
	return nil
}
