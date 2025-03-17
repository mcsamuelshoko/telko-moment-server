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

type userRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) repository.UserRepository {
	return &userRepository{
		Collection: db.Collection("users"),
	}
}

func (u userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := u.Collection.InsertOne(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}
	return nil
}

func (u userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := u.Collection.FindOne(ctx, bson.M{"id": id}).Decode(user)
	if err != nil {
		log.Error().Err(err).Msg("failed to find user with id: " + id)
		return nil, err
	}
	return user, nil

}

func (u userRepository) List(ctx context.Context, page, limit int) ([]models.User, error) {
	cursor, err := u.Collection.Find(ctx, bson.M{})
	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Error().Err(err).Msg("failed to list users")
		return nil, err
	}
	return users, err
}

func (u userRepository) Update(ctx context.Context, user *models.User) error {
	// Use the _id field from the user model for the filter
	// Create an update document with $set to update the user fields
	// Specify the options
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", user}}
	opts := options.Update().SetUpsert(true)

	// Execute the update operation
	_, err := u.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		// Using log.Error() instead of log.Panic() is better for production code
		// as Panic will crash your application
		log.Error().Err(err).Msg("failed to update user with id: " + user.Id.String())
		return err
	}

	return nil
}

func (u userRepository) Delete(ctx context.Context, id string) error {
	_, err := u.Collection.DeleteOne(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete user with id: " + id)
		return err
	}
	return nil
}
