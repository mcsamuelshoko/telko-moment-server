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

type highlightRepository struct {
	Collection *mongo.Collection
}

func newHighlightRepository(db *mongo.Database) repository.HighlightRepository {
	return &highlightRepository{
		Collection: db.Collection("highlights"),
	}
}

func (h highlightRepository) Create(ctx context.Context, highlight *models.Highlight) error {
	_, err := h.Collection.InsertOne(ctx, highlight)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting new highlight")
		return err
	}
	return nil
}

func (h highlightRepository) GetByID(ctx context.Context, id string) (*models.Highlight, error) {
	highlight := models.Highlight{}
	err := h.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&highlight)
	if err != nil {
		log.Error().Err(err).Msg("Error finding highlight with id: " + id)
		return nil, err
	}
	return &highlight, nil
}

func (h highlightRepository) List(ctx context.Context, page, limit int) ([]models.Highlight, error) {
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := h.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("Error finding highlights in collection")
		return nil, err
	}
	defer cursor.Close(ctx)
	var highlights []models.Highlight
	for cursor.Next(ctx) {
		var highlight models.Highlight
		if err := cursor.Decode(&highlight); err != nil {
			log.Error().Err(err).Msg("Error decoding highlight")
			return nil, err
		}
		highlights = append(highlights, highlight)
	}
	return highlights, nil
}

func (h highlightRepository) Update(ctx context.Context, highlight *models.Highlight) error {
	filter := bson.M{"id": highlight.Id}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	_, err := h.Collection.UpdateOne(ctx, filter, opts)
	if err != nil {
		log.Error().Err(err).Msg("Error updating highlight with id: " + highlight.Id.String())
		return err
	}
	return nil
}

func (h highlightRepository) Delete(ctx context.Context, id string) error {
	_, err := h.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		log.Error().Err(err).Msg("Error deleting highlight with id: " + id)
		return err
	}
	return nil
}
