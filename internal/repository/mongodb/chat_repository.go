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

type chatRepository struct {
	iName      string
	logger     *zerolog.Logger
	Collection *mongo.Collection
}

func NewChatRepository(log *zerolog.Logger, db *mongo.Database) repository.ChatRepository {
	return &chatRepository{
		iName:      "ChatRepository",
		logger:     log,
		Collection: db.Collection("chats"),
	}
}

func (c chatRepository) Create(ctx context.Context, chat *models.Chat) (*models.Chat, error) {
	const kName = "CreateChat"

	res, err := c.Collection.InsertOne(ctx, chat)
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to create chat")
		return nil, err
	}
	chat.ID = res.InsertedID.(primitive.ObjectID)
	return chat, nil
}

func (c chatRepository) GetByID(ctx context.Context, id string) (*models.Chat, error) {
	const kName = "GetByID"

	// chat ID to search for
	chatID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to convert id to object id")
		c.logger.Debug().Interface(kName, c.iName).Err(err).Msg("failed to convert Chat GetById -- id:" + id)
		return nil, err
	}

	chat := &models.Chat{}
	err = c.Collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(chat)
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to find chat with id: " + chatID.String())
		return nil, err
	}
	return chat, nil
}

func (c chatRepository) List(ctx context.Context, page, limit int) ([]models.Chat, error) {
	const kName = "List"

	skip := (page - 1) * limit
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	var chats []models.Chat

	cursor, err := c.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to find all chats list")
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to close cursor")
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var chat models.Chat
		if err := cursor.Decode(&chat); err != nil {
			c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to decode chat")
			return nil, err
		}
		chats = append(chats, chat)

	}
	if err := cursor.Err(); err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to find chats")
		return nil, err
	}
	return chats, nil
}

func (c chatRepository) ListByUserId(ctx context.Context, id string, page, limit int) ([]models.Chat, error) {
	const kName = "ListByUserId"

	// participant ID to search for
	participantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to convert id to object id for listing")
		c.logger.Debug().Interface(kName, c.iName).Err(err).Msg("failed to convert ListByUser id:" + id)
		return nil, err
	}
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	// Filter to find chats where the participant ID is in the participants array
	filter := bson.M{"participants": bson.M{"$in": []primitive.ObjectID{participantID}}}

	var chats []models.Chat
	cursor, err := c.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to find chat listByUserId")
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to close cursor")
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var chat models.Chat
		if err := cursor.Decode(&chat); err != nil {
			c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to decode chat")
			return nil, err
		}
		chats = append(chats, chat)

	}
	if err := cursor.Err(); err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to find chats by userId")
		return nil, err
	}
	return chats, nil

}

func (c chatRepository) Update(ctx context.Context, chat *models.Chat) error {
	const kName = "Update"

	filter := bson.M{"_id": chat.ID}
	opts := options.Update().SetUpsert(false)
	_, err := c.Collection.UpdateOne(ctx, filter, bson.M{"$set": chat}, opts)
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to update chat with id: " + chat.ID.String())
		return err
	}
	return nil

}

func (c chatRepository) Delete(ctx context.Context, id string) error {
	const kName = "Delete"

	// chat ID to search for
	chatID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to convert id to object id for deletion")
		c.logger.Debug().Interface(kName, c.iName).Err(err).Msg("failed to convert Delete id:" + id)
		return err
	}

	_, err = c.Collection.DeleteOne(ctx, bson.M{"_id": chatID})
	if err != nil {
		c.logger.Error().Interface(kName, c.iName).Err(err).Msg("failed to delete chat with id: " + id)
		return err
	}
	return nil
}
