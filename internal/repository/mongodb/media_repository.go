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
	// media ID to search for
	mediaID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert Media-GetById id to object id")
		log.Debug().Err(err).Msg("failed to convert media getBy-id:" + id)
		return nil, err
	}
	media := &models.Media{}
	err = m.Collection.FindOne(ctx, bson.M{"_id": mediaID}).Decode(media)
	if err != nil {
		log.Error().Err(err).Msg("failed to find media with id: " + id)
		return nil, err
	}
	return media, nil
}

func (m mediaRepository) GetByChatId(ctx context.Context, chatId string, page, limit int) ([]models.Media, error) {
	// chat ID to search for
	chatID, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert Media-GetByChatId-id to object id")
		log.Debug().Err(err).Msg("failed to convert media-ChatId id:" + chatId)
		return nil, err
	}
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := m.Collection.Find(ctx, bson.M{"chatId": chatID}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to find media in collection from mediaRepository.GetByChatId")
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to close cursor in mediaRepository.GetByChatId")
		}
	}(cursor, ctx)
	var results []models.Media
	if err := cursor.All(ctx, &results); err != nil {
		log.Error().Err(err).Msg("failed to decode results in mediaRepository.GetByChatId")
		return nil, err
	}
	return results, nil
}

func (m mediaRepository) GetBySenderId(ctx context.Context, senderId string, page, limit int) ([]models.Media, error) {

	// sender ID to search for
	senderID, err := primitive.ObjectIDFromHex(senderId)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert id to object id")
		log.Debug().Err(err).Msg("failed to convert id:" + senderId)
		return nil, err
	}
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := m.Collection.Find(ctx, bson.M{"senderId": senderID}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to find media in collection mediaRepository.GetBySenderId")
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to close cursor in mediaRepository.GetBySenderId")
		}
	}(cursor, ctx)

	var results []models.Media
	if err := cursor.All(ctx, &results); err != nil {
		log.Error().Err(err).Msg("failed to decode results")
		return nil, err
	}
	return results, nil
}

func (m mediaRepository) List(ctx context.Context, page, limit int) ([]models.Media, error) {
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := m.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to find media in collection from mediaRepository.List")
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to close cursor in mediaRepository.List")
		}
	}(cursor, ctx)
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
	// media ID to search for
	mediaID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert id to object id")
		log.Debug().Err(err).Msg("failed to convert id:" + id)
		return err
	}

	_, err = m.Collection.DeleteOne(ctx, bson.M{"_id": mediaID})
	if err != nil {
		log.Error().Err(err).Msg("failed to delete media with id: " + id)
		return err
	}
	return nil
}
