package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Defined Chat.Type constants
// for the Chat Model
const (
	ChatTypeDirect  = "direct"
	ChatTypeGroup   = "group"
	ChatTypeChannel = "channel"
	ChatTypeForum   = "forum"
)

type Chat struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	MemberCount   int64              `json:"memberCount,omitempty" bson:"memberCount,omitempty"`
	Type          string             `json:"type" bson:"type"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt     time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	LastMessageID primitive.ObjectID `json:"lastMessageId,omitempty" bson:"lastMessageId,omitempty"`
}
