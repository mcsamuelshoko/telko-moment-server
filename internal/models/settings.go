package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Settings struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	Settings  string             `json:"settings" bson:"settings"`
	CreatedAt primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt primitive.DateTime `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
