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

type chatGroupRepository struct {
	Collection *mongo.Collection
}

func NewChatGroupRepository(db *mongo.Database) repository.ChatGroupRepository {
	return &chatGroupRepository{
		Collection: db.Collection("chat_groups"),
	}
}

func (c chatGroupRepository) Create(ctx context.Context, chatGroup *models.ChatGroup) error {
	_, err := c.Collection.InsertOne(ctx, chatGroup)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert chat_group")
		return err
	}
	return nil
}

func (c chatGroupRepository) GetByID(ctx context.Context, id string) (*models.ChatGroup, error) {
	// chat group ID to search for
	cgID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert id to object id")
		log.Debug().Err(err).Msg("failed to convert id:" + id)
		return nil, err
	}

	chatGroup := &models.ChatGroup{}
	err = c.Collection.FindOne(ctx, bson.M{"_id": cgID}).Decode(chatGroup)
	if err != nil {
		log.Error().Err(err).Msg("failed to find chat_group with id: " + id)
		return nil, err
	}
	return chatGroup, nil
}

func (c chatGroupRepository) List(ctx context.Context, page, limit int) ([]models.ChatGroup, error) {
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := c.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Error().Err(err).Msg("failed to find chat_groups from collection")
		return nil, err
	}
	var chatGroups []models.ChatGroup
	err = cursor.All(ctx, &chatGroups)
	if err != nil {
		log.Error().Err(err).Msg("failed to find chat_groups through cursor")
		return nil, err
	}
	return chatGroups, nil
}

func (c chatGroupRepository) UpdateWithFilter(ctx context.Context, chatGroupId string, updateData map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(chatGroupId)
	if err != nil {
		log.Error().Err(err).Msg("invalid chat group ID format")
		return err
	}

	bsonUpdate := bson.M{}
	for key, value := range updateData {
		bsonUpdate[key] = value
	}

	_, err = c.Collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bsonUpdate})
	if err != nil {
		log.Error().Err(err).Msg("failed to update chat_group")
		return err
	}
	return nil
}

func (c chatGroupRepository) Update(ctx context.Context, chatGroup *models.ChatGroup) error {
	filter := bson.M{"id": chatGroup.Id}
	update := bson.M{"$set": bson.M{}}
	_, err := c.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("failed to update chat_group with id: " + chatGroup.Id.String())
		return err
	}
	return nil
}

func (c chatGroupRepository) Delete(ctx context.Context, id string) error {
	// chat group ID to search for
	cgID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert id to object id")
		log.Debug().Err(err).Msg("failed to convert id:" + id)
		return err
	}

	_, err = c.Collection.DeleteOne(ctx, bson.M{"_id": cgID})
	if err != nil {
		log.Error().Err(err).Msg("failed to delete chat_group with id: " + id)
		return err
	}
	return nil
}
