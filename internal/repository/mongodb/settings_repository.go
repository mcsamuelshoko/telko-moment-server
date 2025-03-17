package mongodb

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type settingsRepository struct {
	Collection *mongo.Collection
}

func NewSettingsRepository(db *mongo.Database) repository.SettingsRepository {
	return &settingsRepository{
		Collection: db.Collection("settings"),
	}
}

func (u settingsRepository) Create(ctx context.Context, settings *models.Settings) error {
	_, err := u.Collection.InsertOne(ctx, settings)
	if err != nil {
		log.Error().Err(err).Msg("failed to create settings")
		return err
	}
	return nil
}

func (u settingsRepository) GetByID(ctx context.Context, id string) (*models.Settings, error) {
	settings := &models.Settings{}
	err := u.Collection.FindOne(ctx, bson.M{"id": id}).Decode(settings)
	if err != nil {
		log.Error().Err(err).Msg("failed to find settings with id: " + id)
		return nil, err
	}
	return settings, nil

}

func (u settingsRepository) List(ctx context.Context, page, limit int) ([]models.Settings, error) {
	// Calculate how many documents to skip
	skip := (page - 1) * limit

	// Create options with pagination parameters
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Execute the find operation with options
	cursor, err := u.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to query settings")
		return nil, err
	}

	// Don't forget to close the cursor when we're done
	defer cursor.Close(ctx)

	// Parse all the documents
	var settingsList []models.Settings
	if err = cursor.All(ctx, &settingsList); err != nil {
		log.Error().Err(err).Msg("failed to decode settingsList")
		return nil, err
	}

	return settingsList, nil
}

func (u settingsRepository) Update(ctx context.Context, settings *models.Settings) error {
	// Use the _id field from the settings model for the filter
	// Create an update document with $set to update the settings fields
	// Specify the options
	filter := bson.D{{"_id", settings.Id}}
	update := bson.D{{"$set", settings}}
	opts := options.Update().SetUpsert(true)

	// Execute the update operation
	_, err := u.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Error().Err(err).Msg("failed to update settings with id: " + settings.Id.String())
		return err
	}

	return nil
}

func (u settingsRepository) Delete(ctx context.Context, id string) error {
	_, err := u.Collection.DeleteOne(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete settings with id: " + id)
		return err
	}
	return nil
}
