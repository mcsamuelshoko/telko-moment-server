package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Settings struct {
	Id        primitive.ObjectID //(references users collection)
	settings  string             //(JSON structure of various settings, e.g., notification preferences, theme, language, etc.)
	createdAt primitive.DateTime
	updatedAt primitive.DateTime
}
