package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	Id                 primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	ChatId             primitive.ObjectID   `json:"chatId" bson:"chatId"`
	SenderId           primitive.ObjectID   `json:"senderId" bson:"senderId"`
	MessageType        string               `json:"messageType" bson:"messageType"`
	Content            string               `json:"content,omitempty" bson:"content,omitempty"`
	MediaUrls          []string             `json:"mediaUrls,omitempty" bson:"mediaUrls,omitempty"`
	Timestamp          primitive.DateTime   `json:"timestamp" bson:"timestamp"`
	EditedTimestamp    primitive.DateTime   `json:"editedTimestamp,omitempty" bson:"editedTimestamp,omitempty"`
	EditedMessage      bool                 `json:"editedMessage" bson:"editedMessage"`
	Status             string               `json:"status" bson:"status"`
	Mentions           []primitive.ObjectID `json:"mentions,omitempty" bson:"mentions,omitempty"`
	RepliedToMessageId primitive.ObjectID   `json:"repliedToMessageId,omitempty" bson:"repliedToMessageId,omitempty"`
}
