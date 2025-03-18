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

type mediaRepository struct {
	Collection *mongo.Collection
}

func NewMediaRepository(db *mongo.Database) repository.MediaRepository {
	return &mediaRepository{
		Collection: db.Collection("medias"),
	}
}

func (m mediaRepository) Create(ctx context.Context, media *models.Media) error {
	_, err := m.Collection.InsertOne(ctx, media)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert media")
		return err
	}
	return nil
}

func (m mediaRepository) GetByID(ctx context.Context, id string) (*models.Media, error) {
	media := &models.Media{}
	err := m.Collection.FindOne(ctx, bson.M{"id": id}).Decode(media)
	if err != nil {
		log.Error().Err(err).Msg("failed to find media with id: " + id)
		return nil, err
	}
	return media, nil
}

func (m mediaRepository) List(ctx context.Context, page, limit int) ([]models.Media, error) {
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := m.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to find media in collection")
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []models.Media
	if err := cursor.All(ctx, &results); err != nil {
		log.Error().Err(err).Msg("failed to decode results")
		return nil, err
	}
	return results, nil
}

func (m mediaRepository) Update(ctx context.Context, media *models.Media) error {
	filter := bson.D{{"_id", media.Id}}
	update := bson.D{{"$set", media}}
	opts := options.Update().SetUpsert(false)
	_, err := m.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Error().Err(err).Msg("failed to update media with id: " + media.Id.String())
		return err
	}
	return nil
}

func (m mediaRepository) Delete(ctx context.Context, id string) error {
	_, err := m.Collection.DeleteOne(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete media with id: " + id)
		return err
	}
	return nil
}
