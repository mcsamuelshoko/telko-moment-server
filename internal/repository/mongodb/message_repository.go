package mongodb

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messageRepository struct {
	Collection *mongo.Collection
	iName      string
}

func NewMessageRepository(db *mongo.Database) repository.MessageRepository {
	return &messageRepository{
		iName:      "MessageRepository",
		Collection: db.Collection("messages"),
	}
}

func (m messageRepository) Create(ctx context.Context, message *models.Message) (*models.Message, error) {
	const kName = "Create"

	res, err := m.Collection.InsertOne(ctx, message)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("Error inserting message")
		return nil, err
	}
	message.ID = res.InsertedID.(primitive.ObjectID)
	return message, nil

}

func (m messageRepository) GetByID(ctx context.Context, id string) (*models.Message, error) {
	const kName = "GetByID"
	// message ID to search for
	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("failed to convert Message-GetById-id to object id")
		log.Debug().Interface(kName, m.iName).Err(err).Msg("failed to convert message-GetById id:" + id)
		return nil, err
	}
	message := &models.Message{}
	err = m.Collection.FindOne(ctx, bson.M{"_id": messageID}).Decode(message)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("Error finding message in GetByID")
		return nil, err
	}
	return message, nil

}

func (m messageRepository) GetByChatID(ctx context.Context, chatId string) (*models.Message, error) {
	const kName = "GetByChatID"

	// chat ID to search for
	chatID, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("failed to convert Message-GetByChatId-id to object id")
		log.Debug().Interface(kName, m.iName).Err(err).Msg("failed to convert message-GetByChat id:" + chatId)
		return nil, err
	}
	message := &models.Message{}
	err = m.Collection.FindOne(ctx, bson.M{"chatId": chatID}).Decode(message)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("Error finding message in GetByChatId")
		return nil, err
	}
	return message, nil

}
func (m messageRepository) GetBySenderID(ctx context.Context, userId string) (*models.Message, error) {
	const kName = "GetBySenderID"

	// sender user ID to search for
	senderID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("failed to convert id to object id")
		log.Debug().Interface(kName, m.iName).Err(err).Msg("failed to convert id:" + userId)
		return nil, err
	}
	message := &models.Message{}
	err = m.Collection.FindOne(ctx, bson.M{"senderId": senderID}).Decode(message)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("Error finding message in GetBySenderId")
		return nil, err
	}
	return message, nil

}

func (m messageRepository) List(ctx context.Context, page, limit int) ([]models.Message, error) {
	const kName = "List"

	// Calculate how many documents to skip
	skip := (page - 1) * limit

	// Create options with pagination parameters
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Execute the find operation with options
	cursor, err := m.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("failed to query settings")
		return nil, err
	}

	// Don't forget to close the cursor when we're done
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Error().Interface(kName, m.iName).Err(err).Msg("failed to close cursor in messageRepository.List")
		}
	}(cursor, ctx)

	// Parse all the documents
	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("failed to decode messages")
		return nil, err
	}
	return messages, nil

}

func (m messageRepository) Update(ctx context.Context, message *models.Message) error {
	const kName = "Update"

	_, err := m.Collection.UpdateOne(ctx, bson.M{"id": message.ID}, bson.M{"$set": message})
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("failed to update message with id: " + message.ID.String())
		return err
	}
	return nil
}

func (m messageRepository) Delete(ctx context.Context, id string) error {
	const kName = "Delete"

	// message ID to search for
	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("failed to convert id to object id")
		log.Debug().Interface(kName, m.iName).Err(err).Msg("failed to convert id:" + id)
		return err
	}
	_, err = m.Collection.DeleteOne(ctx, bson.M{"_id": messageID})
	if err != nil {
		log.Error().Interface(kName, m.iName).Err(err).Msg("failed to delete message with id: " + id)
		return err
	}
	return nil
}
