package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/services"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthenticationRepository struct {
	iName             string
	Collection        *mongo.Collection
	Logger            *zerolog.Logger
	EncryptionService services.IEncryptionService
	SearchKeyHashSvc  services.ISearchKeyService
}

func NewAuthenticationRepository(log *zerolog.Logger, db *mongo.Database, encryptSvc services.IEncryptionService, keyHashSvc services.ISearchKeyService) repository.IAuthenticationRepository {
	return &AuthenticationRepository{
		iName:             "AuthenticationRepository",
		Collection:        db.Collection("authentications"),
		Logger:            log,
		EncryptionService: encryptSvc,
		SearchKeyHashSvc:  keyHashSvc,
	}
}

func (a AuthenticationRepository) Create(ctx context.Context, auth *models.Authentication) (*models.Authentication, error) {
	const kName = "Create"
	// Hash fields
	//err := auth.HashFields(a.SearchKeyHashSvc,"")
	//if err != nil {
	//	a.logger.Error().Err(err).Msg("error hashing authentication fields")
	//	return nil, err
	//}

	// Encrypt sensitive fields before saving
	//if err = auth.EncryptFields(a.EncryptionService); err != nil {
	//	a.logger.Error().Err(err).Msg("Failed to encrypt fields in AuthenticationRepository.Create")
	//	return nil, err

	result, err := a.Collection.InsertOne(ctx, auth)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to create authentication record")
		return nil, err
	}
	auth.ID = result.InsertedID.(primitive.ObjectID)
	return auth, nil
}

func (a AuthenticationRepository) GetList(ctx context.Context) (*[]models.Authentication, error) {
	const kName = "GetList"

	cursor, err := a.Collection.Find(ctx, bson.M{})
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to get authentication list")
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to close cursor")
		}
	}(cursor, ctx)

	var authList []models.Authentication
	//var decryptedAuthList []models.Authentication
	if err := cursor.All(ctx, &authList); err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to decode authentication list")
		return nil, err
	}

	//Decrypt List before sharing to user
	//for _, authentication := range authList {
	//	// Encrypt sensitive fields before saving
	//	if err := authentication.DecryptFields(a.EncryptionService); err != nil {
	//		a.logger.Error().Err(err).Msg("Failed to decrypt fields in AuthenticationRepository.GetList")
	//		return nil, err
	//	}
	//	decryptedAuthList = append(decryptedAuthList, authentication)
	//}

	return &authList, nil
}

func (a AuthenticationRepository) GetByUserID(ctx context.Context, userID string) (*models.Authentication, error) {
	const kName = "GetByUserID"

	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to convert Auth-GetByUserid to object id")
		a.Logger.Debug().Interface(kName, a.iName).Err(err).Msg("Failed to convert auth-user-id:" + userID)
		return nil, err
	}
	var auth models.Authentication
	err = a.Collection.FindOne(ctx, bson.M{"userId": ID}).Decode(&auth)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Str("userID", userID).Msg("Failed to get authentication by user ID")
		return nil, err
	}
	// Decrypt sensitive fields before sharing
	//if err := auth.DecryptFields(a.EncryptionService); err != nil {
	//	a.logger.Error().Err(err).Msg("Failed to decrypt fields in AuthenticationRepository.GetByUserID")
	//	return nil, err
	//}
	return &auth, nil
}

func (a AuthenticationRepository) UpdateByUserID(ctx context.Context, userID string, auth *models.Authentication) (*models.Authentication, error) {
	const kName = "UpdateByUserID"

	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to convert id to object id")
		a.Logger.Debug().Interface(kName, a.iName).Err(err).Msg("Failed to convert id:" + userID)
		return nil, err
	}
	// Encrypt sensitive fields before saving
	//if err := auth.EncryptFields(a.EncryptionService); err != nil {
	//	a.logger.Error().Err(err).Msg("Failed to encrypt fields in AuthenticationRepository.UpdateByUserID")
	//	return nil, err
	//}

	filter := bson.M{"userId": ID}
	update := bson.M{"$set": auth}
	_, err = a.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Str("userID", userID).Msg("Failed to update authentication by user ID")
		return nil, err
	}
	return auth, nil
}

func (a AuthenticationRepository) Delete(ctx context.Context, ID string) error {
	const kName = "Delete"

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Str("ID", ID).Msg("Invalid ID format")
		return err
	}
	_, err = a.Collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Str("ID", ID).Msg("Failed to delete authentication")
		return err
	}
	return nil
}

func (a AuthenticationRepository) DeleteByUserID(ctx context.Context, userID string) error {
	const kName = "DeleteByUserID"

	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to convert id to auth-object id")
		a.Logger.Debug().Interface(kName, a.iName).Err(err).Msg("Failed to convert auth-user-id:" + userID)
		return err
	}
	_, err = a.Collection.DeleteOne(ctx, bson.M{"userId": ID})
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Str("userID", userID).Msg("Failed to delete authentication by user ID")
		return err
	}
	return nil
}

// SaveRefreshToken updates refresh token and adds a fresh one if it does not exist for the user
func (a AuthenticationRepository) SaveRefreshToken(ctx context.Context, userID string, refreshToken string, tokenDuration time.Duration) error {
	const kName = "SaveRefreshToken"

	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to convert id to object id")
		a.Logger.Debug().Interface(kName, a.iName).Err(err).Msg("Failed to convert id:" + userID)
		return err
	}
	// search if user already exists
	var auth models.Authentication
	err = a.Collection.FindOne(ctx, bson.M{"userId": ID}).Decode(&auth)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Str("userID", userID).Msg("Failed to get authentication by user ID")
		// creating new auth object for user
		auth = *models.GetAuthenticationDefaults()
		a.Logger.Info().Interface(kName, a.iName).Str("userID", userID).Msg("Created new authentication from defaults")
	}

	// Hash field(s) for search
	refreshTokenHash, err := a.SearchKeyHashSvc.GenerateSearchKey(refreshToken)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to generate search key")
		return err
	}
	// Encrypt Refresh Token before saving
	//encRefreshToken, err := a.EncryptionService.Encrypt(refreshToken)
	//if err != nil {
	//	a.logger.Error().Err(err).Msg("Failed to encrypt refresh token in AuthenticationRepository.SaveRefreshToken")
	//	return err
	//}

	// Assign to Fields
	//auth.RefreshToken = encRefreshToken
	auth.RefreshTokenHash = refreshTokenHash
	auth.UpdatedAt = time.Now()
	auth.LastLogin = time.Now()
	auth.ExpiresAt = time.Now().Add(tokenDuration)

	// Prepare query and update
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"userId": ID}
	update := bson.D{{"$set", auth}}

	_, err = a.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Str("userID", userID).Msg("Failed to save refresh token")
		return err
	}
	// ########### dumping log #################
	//a.Logger.Debug().Interface(kName, a.iName).
	//	Str("hashedRefreshToken", refreshTokenHash).
	//	Str("refreshToken", refreshToken).
	//	Msg("dumped refresh token & its hash")

	return nil
}

func (a AuthenticationRepository) GetUserIDFromRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	const kName = "GetUserIDFromRefreshToken"

	// Input validation
	if refreshToken == "" {
		a.Logger.Error().Interface(kName, a.iName).Msg("RefreshToken is empty")
		return "", fmt.Errorf("refresh token cannot be empty")
	}

	// Hash token for search
	hashedRefreshToken, err := a.SearchKeyHashSvc.GenerateSearchKey(refreshToken)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to hash refresh token")
		return "", err
	}
	// ########### dumping log #################
	//a.Logger.Debug().Interface(kName, a.iName).
	//	Str("hashedRefreshToken", hashedRefreshToken).
	//	Str("refreshToken", refreshToken).
	//	Msg("dumped refresh token & its hash")

	var result models.Authentication

	// Find the token document
	err = a.Collection.FindOne(ctx, bson.M{
		"refreshTokenHash": hashedRefreshToken,
		"isActive":         true, // Only match active tokens
		"expiresAt": bson.M{
			"$gt": time.Now(), // Only match non-expired tokens
		},
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.Logger.Warn().Interface(kName, a.iName).Msg("No active refresh token found")
			return "", fmt.Errorf("invalid or expired refresh token")
		}
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to get user ID from refresh token")
		return "", err
	}

	if result.UserID.Hex() == "" {
		a.Logger.Warn().Interface(kName, a.iName).Msg("Refresh token found but user ID is empty")
		return "", fmt.Errorf("invalid refresh token: no user associated")
	}

	return result.UserID.Hex(), nil
}

func (a AuthenticationRepository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	const kName = "RevokeRefreshToken"

	// Input validation
	if refreshToken == "" {
		a.Logger.Error().Interface(kName, a.iName).Msg("RefreshToken is empty")
		return fmt.Errorf("refresh token cannot be empty")
	}

	// Hash token for search
	hashedRefreshToken, err := a.SearchKeyHashSvc.GenerateSearchKey(refreshToken)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to hash refresh token")
		return err
	}
	// ########### dumping log #################
	//a.Logger.Debug().Interface(kName, a.iName).
	//	Str("hashedRefreshToken", hashedRefreshToken).
	//	Str("refreshToken", refreshToken).
	//	Msg("dumped refresh token & its hash")

	var result models.Authentication

	// Find the token document
	err = a.Collection.FindOne(ctx, bson.M{
		"refreshTokenHash": hashedRefreshToken,
		"isActive":         true, // Only match active tokens
		"expiresAt": bson.M{
			"$gt": time.Now(), // Only match non-expired tokens
		},
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.Logger.Warn().Interface(kName, a.iName).Msg("No active refresh token found")
			return fmt.Errorf("invalid or expired refresh token")
		}
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to get user ID from refresh token")
		return err
	}

	// deactivate token & expire it
	result.IsActive = false
	result.ExpiresAt = time.Now()
	a.Logger.Debug().Interface(kName, a.iName).Msg("Deactivating refreshToken")
	opts := options.FindOneAndUpdate().SetUpsert(false)
	filter := bson.D{{"_id", result.ID}}
	update := bson.D{{"$set", result}}
	err = a.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to deactivate refresh token")
		return err
	}

	return nil

}

func (a AuthenticationRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	const kName = "DeleteRefreshToken"

	// Hash token for search
	hashedRefreshToken, err := a.SearchKeyHashSvc.GenerateSearchKey(refreshToken)
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to hash refresh token")
		return err
	}
	_, err = a.Collection.DeleteOne(ctx, bson.M{"refreshTokenHash": hashedRefreshToken})
	if err != nil {
		a.Logger.Error().Interface(kName, a.iName).Err(err).Msg("Failed to delete refresh token")
		return err
	}
	return nil
}
