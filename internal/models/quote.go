package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Quote struct {
	Id        primitive.ObjectID //(unique, primary key)
	userId    primitive.ObjectID //(references users collection)
	messageId primitive.ObjectID //(references messages collection)
	item      string             //(e.g., pinned, favorite)
	timestamp primitive.DateTime
}
