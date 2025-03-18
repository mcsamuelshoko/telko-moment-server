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

type chatRepository struct {
	Collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) repository.ChatRepository {
	return &chatRepository{
		Collection: db.Collection("chats"),
	}
}

func (c chatRepository) Create(ctx context.Context, chat *models.Chat) error {
	_, err := c.Collection.InsertOne(ctx, chat)
	if err != nil {
		log.Error().Err(err).Msg("failed to create chat")
		return err
	}
	return nil
}

func (c chatRepository) GetByID(ctx context.Context, id string) (*models.Chat, error) {
	chat := &models.Chat{}
	err := c.Collection.FindOne(ctx, bson.M{"id": id}).Decode(chat)
	if err != nil {
		log.Error().Err(err).Msg("failed to find chat with id: " + chat.Id.String())
		return nil, err
	}
	return chat, nil
}

func (c chatRepository) List(ctx context.Context, page, limit int) ([]models.Chat, error) {
	skip := (page - 1) * limit
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	var chats []models.Chat

	cursor, err := c.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to find chats")
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var chat models.Chat
		if err := cursor.Decode(&chat); err != nil {
			log.Error().Err(err).Msg("failed to decode chat")
			return nil, err
		}
		chats = append(chats, chat)

	}
	if err := cursor.Err(); err != nil {
		log.Error().Err(err).Msg("failed to find chats")
		return nil, err
	}
	return chats, nil
}

func (c chatRepository) Update(ctx context.Context, chat *models.Chat) error {
	filter := bson.M{"id": chat.Id}
	opts := options.Update().SetUpsert(false)
	_, err := c.Collection.UpdateOne(ctx, filter, bson.M{"$set": chat}, opts)
	if err != nil {
		log.Error().Err(err).Msg("failed to update chat with id: " + chat.Id.String())
		return err
	}
	return nil

}

func (c chatRepository) Delete(ctx context.Context, id string) error {
	_, err := c.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		log.Error().Err(err).Msg("failed to delete chat with id: " + id)
		return err
	}
	return nil
}
