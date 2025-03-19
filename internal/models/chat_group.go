package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatGroup struct {
	Id          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	GroupName   string               `json:"groupName" bson:"groupName"`
	Members     []primitive.ObjectID `json:"members" bson:"members"`
	AdminIds    []primitive.ObjectID `json:"adminIds" bson:"adminIds"`
	Description string               `json:"description,omitempty" bson:"description,omitempty"`
	ProfileUrl  string               `json:"profileUrl,omitempty" bson:"profileUrl,omitempty"`
	CreatedAt   primitive.DateTime   `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   primitive.DateTime   `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	ChatId      primitive.ObjectID   `json:"chatId" bson:"chatId"`
}
