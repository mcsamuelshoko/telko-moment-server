package mongodb

import (
	"context"
	"errors"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/services"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	iName                string
	Collection           *mongo.Collection
	Log                  *zerolog.Logger
	EncryptionService    services.IEncryptionService
	SearchKeyHashService services.ISearchKeyService
}

func NewUserRepository(log *zerolog.Logger, db *mongo.Database, encryptSvc services.IEncryptionService, sKeyHashSvc services.ISearchKeyService) repository.IUserRepository {
	return &userRepository{
		iName:                "UserRepository",
		Collection:           db.Collection("users"),
		Log:                  log,
		EncryptionService:    encryptSvc,
		SearchKeyHashService: sKeyHashSvc,
	}
}

func (u userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	//Hash fields used in search
	err := user.HashFields(u.SearchKeyHashService)
	if err != nil {
		u.Log.Error().Err(err).Msg("error hashing user fields")
		return nil, err
	}
	//Encrypt fields before saving
	err = user.EncryptFields(u.EncryptionService)
	if err != nil {
		u.Log.Error().Interface("Create", u.iName).Err(err).Msg("error encrypting user")
		return nil, err
	}

	// Insert User
	res, err := u.Collection.InsertOne(ctx, user)
	if err != nil {
		u.Log.Error().Interface("Create", u.iName).Err(err).Msg("Failed to create user")
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)

	// Decrypt for use
	err = user.DecryptFields(u.EncryptionService)
	if err != nil {
		u.Log.Error().Interface("Create", u.iName).Err(err).Msg("Failed to decrypt new user")
		return nil, err
	}
	return user, nil
}

func (u userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	// user ID to search for
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		u.Log.Error().Interface("GetByID", u.iName).Err(err).Msg("Failed to convert id to object id")
		u.Log.Debug().Interface("GetByID", u.iName).Err(err).Msg("Failed to convert id:" + id)
		return nil, err
	}

	user := &models.User{}
	err = u.Collection.FindOne(ctx, bson.M{"_id": userID}).Decode(user)
	if err != nil {
		u.Log.Error().Interface("GetByID", u.iName).Err(err).Msg("Failed to find user with id: " + id)
		return nil, err
	}
	err = user.DecryptFields(u.EncryptionService)
	if err != nil {
		u.Log.Error().Err(err).Interface("GetByID", u.iName).Msg("Failed to decrypt user with id: " + id)
		return nil, err
	}
	return user, nil

}

func (u userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	// encrypt before search
	hashedEmail, err := u.SearchKeyHashService.GenerateSearchKey(email)
	if err != nil {
		u.Log.Debug().Interface("GetByEmail", u.iName).Err(err).Msg("Failed to hash user Email: " + email)
		u.Log.Error().Interface("GetByEmail", u.iName).Err(err).Msg("Failed to hash user Email")
		return nil, err
	}
	err = u.Collection.FindOne(ctx, bson.M{"emailHash": hashedEmail}).Decode(user)
	if err != nil {
		u.Log.Debug().Interface("GetByEmail", u.iName).Err(err).Msg("Failed to find user with email: " + email + " :::: " + hashedEmail)
		u.Log.Error().Interface("GetByEmail", u.iName).Err(err).Msg("Failed to find user with provided email")
		return nil, err
	}
	err = user.DecryptFields(u.EncryptionService)
	if err != nil {
		u.Log.Error().Interface("GetByEmail", u.iName).Err(err).Msg("Failed to decrypt user with email: " + email)
		return nil, err
	}
	return user, nil
}

func (u userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	// encrypt before search
	hashedUsername, err := u.SearchKeyHashService.GenerateSearchKey(username)
	if err != nil {
		u.Log.Debug().Interface("GetByUsername", u.iName).Err(err).Msg("Failed to hash user with username: " + username + " :::: " + hashedUsername)
		u.Log.Error().Interface("GetByUsername", u.iName).Err(err).Msg("Failed to hash username: ")
		return nil, err
	}
	err = u.Collection.FindOne(ctx, bson.M{"usernameHash": hashedUsername}).Decode(user)
	if err != nil {
		u.Log.Error().Interface("GetByUsername", u.iName).Err(err).Msg("Failed to find user with username: " + username)
		u.Log.Error().Interface("GetByUsername", u.iName).Err(err).Msg("Failed to find user with provided username")
		return nil, err
	}
	err = user.DecryptFields(u.EncryptionService)
	if err != nil {
		u.Log.Error().Interface("GetByUsername", u.iName).Err(err).Msg("Failed to decrypt user with username: " + username)
		return nil, err
	}
	return user, nil
}

func (u userRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	user := &models.User{}
	// encrypt before search
	hashedPhoneNumber, err := u.SearchKeyHashService.GenerateSearchKey(phoneNumber)
	if err != nil {
		u.Log.Error().Interface("GetByPhoneNumber", u.iName).Err(err).
			Msg("Failed to encrypt user with phoneNumber: " + phoneNumber + " :::: " + hashedPhoneNumber)
		return nil, err
	}
	err = u.Collection.FindOne(ctx, bson.M{"phoneNumberHash": hashedPhoneNumber}).Decode(user)
	if err != nil {
		u.Log.Error().Interface("GetByPhoneNumber", u.iName).Err(err).Msg("Failed to find user with phone number: " + phoneNumber)
		return nil, err
	}
	err = user.DecryptFields(u.EncryptionService)
	if err != nil {
		u.Log.Error().Interface("GetByPhoneNumber", u.iName).Err(err).Msg("Failed to decrypt user with phone number: " + phoneNumber)
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

	// decrypt user fields
	for i := 0; i < len(users); i++ {
		err = users[i].DecryptFields(u.EncryptionService)
		if err != nil {
			u.Log.Error().Err(err).Msg("Failed to decrypt user with username: " + users[i].Username)
		}

	}

	return users, nil
}

func (u userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	//// Use the _id field from the user model for the filter
	//// Create an update document with $set to update the user fields
	//// Specify the options
	//filter := bson.D{{"_id", user.ID}}
	//update := bson.D{{"$set", user}}
	//opts := options.Update().SetUpsert(true)
	//
	//// Execute the update operation
	//res, err := u.Collection.UpdateOne(ctx, filter, update, opts)
	//if err != nil {
	//	u.Log.Error().Err(err).Msg("Failed to update user with id: " + user.ID.String())
	//	return nil, err
	//}
	//
	//return nil

	// Create a filter using the _id field
	filter := bson.D{{"_id", user.ID}}

	// Create an update document, excluding the _id field
	update := bson.D{{"$set", user}}

	// Options to return the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)

	// Create a variable to store the updated user
	var updatedUser models.User

	// Execute the update and retrieve the updated document
	err := u.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			u.Log.Warn().Msg("No document found to update with id: " + user.ID.String())
			return nil, err
		}
		u.Log.Error().Err(err).Msg("Failed to update user with id: " + user.ID.String())
		return nil, err
	}

	return &updatedUser, nil
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
