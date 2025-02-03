package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatGroup struct {
	groupId     primitive.ObjectID //(unique, primary key)
	groupName   string
	members     []primitive.ObjectID //(array of user IDs)
	adminIds    []primitive.ObjectID //(array of admin user IDs)
	description string
	profileUrl  string
	createdAt   primitive.DateTime
	updatedAt   primitive.DateTime
	chatId      primitive.ObjectID //(references chats collection for direct chat links)
}
