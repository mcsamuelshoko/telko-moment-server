package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chat struct {
	Id            primitive.ObjectID   //(unique, primary key)
	participants  []primitive.ObjectID //(array of user IDs)
	chatType      string               //(direct, group, etc.)
	chatName      string               //(optional, for group chats)
	createdAt     primitive.DateTime
	updatedAt     primitive.DateTime
	lastMessageId primitive.ObjectID //(references the last message)
}
