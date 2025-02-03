package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Highlight struct {
	highlightId primitive.ObjectID //(unique, primary key)
	userId      primitive.ObjectID //(references users collection)
	quote       string             //(references messages collection)
	item        string             //(e.g., pinned, favorite)
	timestamp   primitive.DateTime
}
