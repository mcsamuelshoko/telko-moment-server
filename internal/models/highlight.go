package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Highlight struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	Quote     string             `json:"quote" bson:"quote"`
	Item      string             `json:"item" bson:"item"`
	Timestamp primitive.DateTime `json:"timestamp" bson:"timestamp"`
}
