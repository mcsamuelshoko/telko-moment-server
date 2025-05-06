package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Define Settings.Preferences.AutoDownloadMedia constants
// for the Settings Model
const (
	SettingsPrefAutoDownloadMediaWifi     = "wifi"
	SettingsPrefAutoDownloadMediaCellular = "cellular"
	SettingsPrefAutoDownloadMediaNever    = "never"
)

// Defined Settings.Preferences.LastActiveVisibility constants
// for the Settings Model
const (
	SettingsPrefLastActiveVisibilityEveryone = "everyone"
	SettingsPrefLastActiveVisibilityContacts = "contacts"
	SettingsPrefLastActiveVisibilityNobody   = "nobody"
)

// Defined Settings.Preferences.HighlightVisibility constants
// for the Settings Model
const (
	SettingsPrefHighlightVisibilityEveryone = "everyone"
	SettingsPrefHighlightVisibilityContacts = "contacts"
	SettingsPrefHighlightVisibilityNobody   = "nobody"
)

// Note: For other constants there is a consideration to give breathing space for the frontend/UI team
// so they can be flexible and creative, since the constants do not affect backend or have effects on
// the backend's functionality e.g. theme, emoji etc

type UserPreferences struct {
	Theme                string `json:"theme"`                // "light", "dark", "system"
	Notifications        bool   `json:"notifications"`        // Enable/disable notifications
	Sound                bool   `json:"sound"`                // Enable/disable message sounds
	Vibration            bool   `json:"vibration"`            // Enable/disable message vibration
	FontSize             int    `json:"fontSize"`             // Font size for chat messages
	Language             string `json:"language"`             // User's preferred language (e.g., "en", "es", "fr")
	ShowPreviews         bool   `json:"showPreviews"`         // Show message previews in notifications
	AutoDownloadMedia    string `json:"autoDownloadMedia"`    // "wifi", "cellular", "never"
	ReadReceipts         bool   `json:"readReceipts"`         // Enable/disable read receipts
	LastActiveVisibility string `json:"lastActiveVisibility"` // "everyone", "contacts", "nobody"
	HighlightVisibility  string `json:"highlightVisibility"`  // "everyone", "contacts", "nobody"
	ChatWallpaper        string `json:"chatWallpaper"`        // Path or ID of the chat wallpaper
	EmojiStyle           string `json:"emojiStyle"`           // "system", "apple", "google", etc.
	Accessibility        struct {
		HighContrast bool `json:"highContrast"`
		TextToSpeech bool `json:"textToSpeech"`
	} `json:"accessibility"`
	Privacy struct {
		ProfilePictureVisibility string `json:"profilePictureVisibility"` // "everyone", "contacts", "nobody"
		PhoneNumberVisibility    string `json:"phoneNumberVisibility"`    // "everyone", "contacts", "nobody"
		AddToGroups              string `json:"addToGroups"`              // "everyone", "contacts", "contacts-of-contacts", "nobody"
	} `json:"privacy"`
	// Add more preferences as needed
}

type Settings struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId"`
	Preferences UserPreferences    `json:"preferences" bson:"preferences"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// :::: UNIQUE INDEXES

// CreateUniqueIndexes creates unique indexes for username, phoneNumber and email
func (s *Settings) CreateUniqueIndexes(db *mongo.Database) error {
	// Create unique index for userId
	userIdIndex := mongo.IndexModel{
		Keys:    bson.D{{"userId", 1}},
		Options: options.Index().SetUnique(true).SetName("unique_user_id"),
	}

	// Create indexes
	_, err := db.Collection("settings").Indexes().
		CreateMany(context.Background(), []mongo.IndexModel{userIdIndex})

	return err
}

// :::: DEFAULTS FUNCTION(S)

// GetSettingsDefaultsFromHeaders Get user defaults from request headers
func GetSettingsDefaultsFromHeaders(headers map[string]string) *Settings {
	settings := &Settings{
		Preferences: UserPreferences{
			AutoDownloadMedia:    SettingsPrefAutoDownloadMediaWifi,
			EmojiStyle:           "system",
			Language:             "en",
			LastActiveVisibility: SettingsPrefLastActiveVisibilityEveryone,
			HighlightVisibility:  SettingsPrefHighlightVisibilityEveryone,
			Notifications:        true,
			ReadReceipts:         true,
			Sound:                true,
			Theme:                "system",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return settings
}
