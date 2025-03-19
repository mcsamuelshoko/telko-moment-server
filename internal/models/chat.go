package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chat struct {
	Id            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Participants  []primitive.ObjectID `json:"participants" bson:"participants"`
	ChatType      string               `json:"chatType" bson:"chatType"`
	ChatName      string               `json:"chatName,omitempty" bson:"chatName,omitempty"`
	CreatedAt     primitive.DateTime   `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt     primitive.DateTime   `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	LastMessageId primitive.ObjectID   `json:"lastMessageId,omitempty" bson:"lastMessageId,omitempty"`
}
