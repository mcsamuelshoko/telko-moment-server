package mongodb

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type settingsRepository struct {
	Collection *mongo.Collection
	Logger     *zerolog.Logger
}

func NewSettingsRepository(log *zerolog.Logger, db *mongo.Database) repository.ISettingsRepository {
	return &settingsRepository{
		Collection: db.Collection("settings"),
		Logger:     log,
	}
}

func (s settingsRepository) Create(ctx context.Context, settings *models.Settings) (*models.Settings, error) {
	result, err := s.Collection.InsertOne(ctx, settings)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to create settings")
		return nil, err
	}
	settings.ID = result.InsertedID.(primitive.ObjectID)
	return settings, nil
}

func (s settingsRepository) GetByID(ctx context.Context, id string) (*models.Settings, error) {
	// settings ID to search for
	settingsID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to convert settings id to object id")
		s.Logger.Debug().Err(err).Msg("failed to convert GetByID id:" + id)
		return nil, err
	}

	settings := &models.Settings{}
	err = s.Collection.FindOne(ctx, bson.M{"_id": settingsID}).Decode(settings)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to find settings with id: " + id)
		return nil, err
	}
	return settings, nil

}

func (s settingsRepository) GetByUserID(ctx context.Context, userId string) (*models.Settings, error) {
	// user ID to search for
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to convert userId to object id")
		s.Logger.Debug().Err(err).Msg("failed to convert GetByUserID id:" + userId)
		return nil, err
	}

	settings := &models.Settings{}
	err = s.Collection.FindOne(ctx, bson.M{"userId": userID}).Decode(settings)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to find settings with id: " + userId)
		return nil, err
	}
	return settings, nil
}

func (s settingsRepository) List(ctx context.Context, page, limit int) ([]models.Settings, error) {
	// Calculate how many documents to skip
	skip := (page - 1) * limit

	// Create options with pagination parameters
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Execute the find operation with options
	cursor, err := s.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to query settings")
		return nil, err
	}

	// Don't forget to close the cursor when we're done
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			s.Logger.Error().Err(err).Msg("failed to close cursor")
		}
	}(cursor, ctx)

	// Parse all the documents
	var settingsList []models.Settings
	if err = cursor.All(ctx, &settingsList); err != nil {
		s.Logger.Error().Err(err).Msg("failed to decode settingsList")
		return nil, err
	}

	return settingsList, nil
}

func (s settingsRepository) Update(ctx context.Context, settings *models.Settings) error {
	// Use the _id field from the settings model for the filter
	// Create an update document with $set to update the settings fields
	// Specify the options
	filter := bson.D{{"_id", settings.ID}}
	update := bson.D{{"$set", settings}}
	opts := options.Update().SetUpsert(false)

	// Execute the update operation
	_, err := s.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to update settings with id: " + settings.ID.String())
		return err
	}

	return nil
}

func (s settingsRepository) Delete(ctx context.Context, id string) error {
	// settings ID to search for
	settingsID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to convert Settings, Delete id to object id")
		s.Logger.Debug().Err(err).Msg("failed to convert Settings Delete id:" + id)
		return err
	}
	_, err = s.Collection.DeleteOne(ctx, bson.M{"_id": settingsID})
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to delete settings with id: " + id)
		return err
	}
	return nil
}
