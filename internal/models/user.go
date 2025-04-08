package models

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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
	CreatedAt          time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt          time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// CreateUniqueIndexes creates unique indexes for username, phoneNumber and email
func (u *User) CreateUniqueIndexes(db *mongo.Database) error {
	// Create unique index for username
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{"username", 1}},
		Options: options.Index().SetUnique(true).SetName("unique_username"),
	}

	//// Create unique index for email
	//emailIndex := mongo.IndexModel{
	//	Keys:    bson.D{{"email", 1}},
	//	Options: options.Index().SetUnique(true).SetName("unique_email"),
	//}
	//
	//// Create unique index for phoneNumber
	//phoneNumberIndex := mongo.IndexModel{
	//	Keys:    bson.D{{"phoneNumber", 1}},
	//	Options: options.Index().SetUnique(true).SetName("unique_phone_number"),
	//}

	emailAndPhoneNumberIndex := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}, {"phoneNumber", 1}},
		Options: options.Index().SetUnique(true).SetName("unique_email_and_phone_number"),
	}

	// Create indexes
	_, err := db.Collection("users").Indexes().
		CreateMany(context.Background(), []mongo.IndexModel{usernameIndex, emailAndPhoneNumberIndex})

	return err
}

// EncryptFields Encrypt sensitive fields before saving
func (u *User) EncryptFields(encSvc services.IEncryptionService) error {
	if u.FirstName != "" {
		encrypted, err := encSvc.Encrypt(u.FirstName)
		if err != nil {
			return err
		}
		u.FirstName = encrypted
	}
	if u.LastName != "" {
		encrypted, err := encSvc.Encrypt(u.LastName)
		if err != nil {
			return err
		}
		u.LastName = encrypted
	}
	if u.Username != "" {
		encrypted, err := encSvc.Encrypt(u.Username)
		if err != nil {
			return err
		}
		u.Username = encrypted
	}
	if u.Email != "" {
		encrypted, err := encSvc.Encrypt(u.Email)
		if err != nil {
			return err
		}
		u.Email = encrypted
	}
	if u.PhoneNumber != "" {
		encrypted, err := encSvc.Encrypt(u.PhoneNumber)
		if err != nil {
			return err
		}
		u.PhoneNumber = encrypted
	}
	if u.UserType != "" {
		encrypted, err := encSvc.Encrypt(u.UserType)
		if err != nil {
			return err
		}
		u.UserType = encrypted
	}
	if u.Password != "" {
		encrypted, err := encSvc.Encrypt(u.Password)
		if err != nil {
			return err
		}
		u.Password = encrypted
	}
	if u.Bio != "" {
		encrypted, err := encSvc.Encrypt(u.Bio)
		if err != nil {
			return err
		}
		u.Bio = encrypted
	}
	if u.ProfilePicture != "" {
		encrypted, err := encSvc.Encrypt(u.ProfilePicture)
		if err != nil {
			return err
		}
		u.ProfilePicture = encrypted
	}
	if u.Country != "" {
		encrypted, err := encSvc.Encrypt(u.Country)
		if err != nil {
			return err
		}
		u.Country = encrypted
	}

	return nil
}

// DecryptFields Decrypt sensitive fields after retrieval
func (u *User) DecryptFields(encSvc services.IEncryptionService) error {
	if u.FirstName != "" {
		decrypted, err := encSvc.Decrypt(u.FirstName)
		if err != nil {
			return err
		}
		u.FirstName = decrypted
	}
	if u.LastName != "" {
		decrypted, err := encSvc.Decrypt(u.LastName)
		if err != nil {
			return err
		}
		u.LastName = decrypted
	}
	if u.Username != "" {
		decrypted, err := encSvc.Decrypt(u.Username)
		if err != nil {
			return err
		}
		u.Username = decrypted
	}
	if u.Email != "" {
		decrypted, err := encSvc.Decrypt(u.Email)
		if err != nil {
			return err
		}
		u.Email = decrypted
	}
	if u.PhoneNumber != "" {
		decrypted, err := encSvc.Decrypt(u.PhoneNumber)
		if err != nil {
			return err
		}
		u.PhoneNumber = decrypted
	}
	if u.UserType != "" {
		decrypted, err := encSvc.Decrypt(u.UserType)
		if err != nil {
			return err
		}
		u.UserType = decrypted
	}
	//if u.Password != "" {
	//	decrypted, err := encSvc.Decrypt(u.Password)
	//	if err != nil {
	//		return err
	//	}
	//	u.Password = decrypted
	//}
	if u.Bio != "" {
		decrypted, err := encSvc.Decrypt(u.Bio)
		if err != nil {
			return err
		}
		u.Bio = decrypted
	}
	if u.ProfilePicture != "" {
		decrypted, err := encSvc.Decrypt(u.ProfilePicture)
		if err != nil {
			return err
		}
		u.ProfilePicture = decrypted
	}
	if u.Country != "" {
		decrypted, err := encSvc.Decrypt(u.Country)
		if err != nil {
			return err
		}
		u.Country = decrypted
	}

	return nil
}

// :::: SANITIZER HELPER FUNCTIONS

// Sanitize Helper function to remove sensitive fields from user data
func (u *User) Sanitize() map[string]interface{} {
	return map[string]interface{}{
		"id":                 u.ID,
		"firstName":          u.FirstName,
		"lastName":           u.LastName,
		"username":           u.Username,
		"email":              u.Email,
		"phoneNumber":        u.PhoneNumber,
		"userType":           u.UserType,
		"profilePicture":     u.ProfilePicture,
		"status":             u.Status,
		"bio":                u.Bio,
		"languagePreference": u.LanguagePreference,
		"timezone":           u.Timezone,
		"country":            u.Country,
		"createdAt":          u.CreatedAt,
		"updatedAt":          u.UpdatedAt,
	}
}

// :::: DEFAULTS FUNCTION(S)

// GetUserDefaultsFromHeaders Get user defaults from request headers
func GetUserDefaultsFromHeaders(headers map[string]string) *User {
	user := &User{
		Status:             "active",
		LanguagePreference: "en",
		UserType:           "regular",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Get timezone from header or default to UTC
	timezone := headers["Timezone"]
	if timezone == "" {
		timezone = "UTC"
	}
	user.Timezone = timezone

	// Attempt to get country from request headers
	country := headers["CF-IPCountry"]
	if country == "" {
		country = headers["X-Country"]
	}
	user.Country = country

	return user
}
