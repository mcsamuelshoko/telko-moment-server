package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName          string             `json:"firstName" bson:"firstName"`
	LastName           string             `json:"lastName" bson:"lastName"`
	Username           string             `json:"username" bson:"username"`
	Password           string             `json:"password,omitempty" bson:"password"`
	Email              string             `json:"email" bson:"email"`
	PhoneNumber        string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	UserType           string             `json:"userType" bson:"userType"`
	ProfilePicture     string             `json:"profilePicture,omitempty" bson:"profilePicture,omitempty"`
	Status             string             `json:"status" bson:"status"`
	Bio                string             `json:"bio,omitempty" bson:"bio,omitempty"`
	LanguagePreference string             `json:"languagePreference,omitempty" bson:"languagePreference,omitempty"`
	Timezone           string             `json:"timezone,omitempty" bson:"timezone,omitempty"`
	Country            string             `json:"country,omitempty" bson:"country,omitempty"`
	CreatedAt          primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt          primitive.DateTime `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// CreateUniqueIndexes creates unique indexes for username and email
func CreateUniqueIndexes(collection *mongo.Collection) error {
	// Create unique index for username
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{"username", 1}},
		Options: options.Index().SetUnique(true).SetName("unique_username"),
	}

	// Create unique index for email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true).SetName("unique_email"),
	}

	// Create indexes
	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{usernameIndex, emailIndex})

	return err
}
