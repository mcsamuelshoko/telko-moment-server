package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authentication struct {
	Id           primitive.ObjectID // (references users collection)
	email        string             //(unique)
	passwordHash string
	refreshToken string
	authProvider string //(Google, Facebook, etc.)
	lastLogin    primitive.DateTime
	createdAt    primitive.DateTime
	updatedAt    primitive.DateTime
}
