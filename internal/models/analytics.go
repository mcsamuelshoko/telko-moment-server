package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Analytics struct {
	analyticsId primitive.ObjectID //(unique, primary key)
	userId      primitive.ObjectID //(references users collection)
	metrics     string             //(performance metrics like message count, active hours, etc.)
	data        string             //(raw data about messages, calls, activity, etc.)
	createdAt   primitive.DateTime
	updatedAt   primitive.DateTime
}
