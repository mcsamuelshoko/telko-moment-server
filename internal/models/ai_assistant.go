package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AIAssistant struct {
	assistantId     primitive.ObjectID // (unique, primary key)
	userId          primitive.ObjectID //(references users collection)
	assistantName   string             //(name of the assistant)
	status          string             //(active, inactive)
	lastInteraction primitive.DateTime
	preferences     primitive.ObjectID //(json AI preferences, like response style, etc.)
	createdAt       primitive.DateTime
	updatedAt       primitive.DateTime
}
