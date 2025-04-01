package models

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Authentication represents the stored authentication information in MongoDB
type Authentication struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	RefreshToken string             `json:"refreshToken,omitempty" bson:"refreshToken,omitempty"`
	IsActive     bool               `json:"isActive" bson:"isActive"`
	ExpiresAt    time.Time          `json:"expiresAt" bson:"expiresAt"`
	AuthProvider string             `json:"authProvider" bson:"authProvider"`
	LastLogin    time.Time          `json:"lastLogin,omitempty" bson:"lastLogin,omitempty"`
	CreatedAt    time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt    time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// CreateUniqueIndexes creates unique indexes for userId and refreshToken
func (a *Authentication) CreateUniqueIndexes(db *mongo.Database) error {
	// Create unique index for userId
	userIdIndex := mongo.IndexModel{
		Keys:    bson.D{{"userId", 1}},
		Options: options.Index().SetUnique(true).SetName("unique_user_id"),
	}

	// Create unique index for refreshToken
	refreshTokenIndex := mongo.IndexModel{
		Keys:    bson.D{{"refreshToken", 1}},
		Options: options.Index().SetUnique(true).SetName("unique_refresh_token"),
	}

	// Create indexes
	_, err := db.Collection("authentications").Indexes().
		CreateMany(context.Background(), []mongo.IndexModel{userIdIndex, refreshTokenIndex})

	return err
}

// EncryptFields Encrypt sensitive fields before saving
func (a *Authentication) EncryptFields(encSvc utils.IEncryptionService) error {
	if a.RefreshToken != "" {
		encrypted, err := encSvc.Encrypt(a.RefreshToken)
		if err != nil {
			return err
		}
		a.RefreshToken = encrypted
	}
	return nil
}

// DecryptFields Decrypt sensitive fields after retrieval
func (a *Authentication) DecryptFields(encSvc utils.IEncryptionService) error {
	if a.RefreshToken != "" {
		decrypted, err := encSvc.Decrypt(a.RefreshToken)
		if err != nil {
			return err
		}
		a.RefreshToken = decrypted
	}
	return nil
}

// LoginRequest represents the data needed for a login attempt
type LoginRequest struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required"`
}

// LoginRequestEmail represents the data needed for a login attempt
type LoginRequestEmail struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the data returned after a successful login
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// RegisterRequest represents the data needed for a registration attempt
type RegisterRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,phoneNumber"`
	Password    string `json:"password" validate:"required"`
}

// RegisterRequestEmail represents the data needed for a registration attempt
type RegisterRequestEmail struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
