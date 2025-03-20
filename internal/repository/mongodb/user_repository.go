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

type userRepository struct {
	Collection *mongo.Collection
	Log        *zerolog.Logger
}

func NewUserRepository(log *zerolog.Logger, db *mongo.Database) repository.UserRepository {
	return &userRepository{
		Collection: db.Collection("users"),
		Log:        log,
	}
}

func (u userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := u.Collection.InsertOne(ctx, user)
	if err != nil {
		u.Log.Error().Err(err).Msg("Failed to create user")
		return err
	}
	return nil
}

func (u userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	// user ID to search for
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		u.Log.Error().Err(err).Msg("Failed to convert id to object id")
		u.Log.Debug().Err(err).Msg("Failed to convert id:" + id)
		return nil, err
	}

	user := &models.User{}
	err = u.Collection.FindOne(ctx, bson.M{"_id": userID}).Decode(user)
	if err != nil {
		u.Log.Error().Err(err).Msg("Failed to find user with id: " + id)
		return nil, err
	}
	return user, nil

}

func (u userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := u.Collection.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		u.Log.Error().Err(err).Msg("Failed to find user with email: " + email)
		return nil, err
	}
	return user, nil
}

func (u userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := u.Collection.FindOne(ctx, bson.M{"username": username}).Decode(user)
	if err != nil {
		u.Log.Error().Err(err).Msg("Failed to find user with username: " + username)
		return nil, err
	}
	return user, nil
}

func (u userRepository) List(ctx context.Context, page, limit int) ([]models.User, error) {
	// Calculate how many documents to skip
	skip := (page - 1) * limit

	// Create options with pagination parameters
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Execute the find operation with options
	cursor, err := u.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		u.Log.Error().Err(err).Msg("Failed to query users")
		return nil, err
	}

	// Don't forget to close the cursor when we're done
	defer cursor.Close(ctx)

	// Parse all the documents
	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		u.Log.Error().Err(err).Msg("Failed to decode users")
		return nil, err
	}

	return users, nil
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
		u.Log.Error().Err(err).Msg("Failed to update user with id: " + user.Id.String())
		return err
	}

	return nil
}

func (u userRepository) Delete(ctx context.Context, id string) error {
	// user ID to search for
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		u.Log.Error().Err(err).Msg("Failed to convert id to object id")
		u.Log.Debug().Err(err).Msg("Failed to convert id:" + id)
		return err
	}

	_, err = u.Collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		u.Log.Error().Err(err).Msg("Failed to delete user with id: " + id)
		return err
	}
	return nil
}
