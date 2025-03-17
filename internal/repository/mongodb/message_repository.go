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

type messageRepository struct {
	Collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) repository.MessageRepository {
	return &messageRepository{
		Collection: db.Collection("messages"),
	}
}

func (m messageRepository) Create(ctx context.Context, message *models.Message) error {
	_, err := m.Collection.InsertOne(ctx, message)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting message")
		return err
	}
	return nil

}

func (m messageRepository) GetByID(ctx context.Context, id string) (*models.Message, error) {
	message := &models.Message{}
	err := m.Collection.FindOne(ctx, bson.M{"id": id}).Decode(message)
	if err != nil {
		log.Error().Err(err).Msg("Error finding message")
		return nil, err
	}
	return message, nil

}

func (m messageRepository) List(ctx context.Context, page, limit int) ([]models.Message, error) {
	// Calculate how many documents to skip
	skip := (page - 1) * limit

	// Create options with pagination parameters
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Execute the find operation with options
	cursor, err := m.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to query settings")
		return nil, err
	}

	// Don't forget to close the cursor when we're done
	defer cursor.Close(ctx)

	// Parse all the documents
	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		log.Error().Err(err).Msg("failed to decode messages")
		return nil, err
	}
	return messages, nil

}

func (m messageRepository) Update(ctx context.Context, message *models.Message) error {
	_, err := m.Collection.UpdateOne(ctx, bson.M{"id": message.Id}, bson.M{"$set": message})
	if err != nil {
		log.Error().Err(err).Msg("failed to update message with id: " + message.Id.String())
		return err
	}
	return nil
}

func (m messageRepository) Delete(ctx context.Context, id string) error {
	_, err := m.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		log.Error().Err(err).Msg("failed to delete message with id: " + id)
		return err
	}
	return nil
}
