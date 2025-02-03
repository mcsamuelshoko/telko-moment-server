package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Settings struct {
	userId     primitive.ObjectID //(references users collection)
	settings   string             //(JSON structure of various settings, e.g., notification preferences, theme, language, etc.)
	created_at primitive.DateTime
	updated_at primitive.DateTime
}
