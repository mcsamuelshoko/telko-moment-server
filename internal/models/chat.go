package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chat struct {
	chat_id         primitive.ObjectID   //(unique, primary key)
	participants    []primitive.ObjectID //(array of user IDs)
	chat_type       string               //(direct, group, etc.)
	chat_name       string               //(optional, for group chats)
	created_at      primitive.DateTime
	updated_at      primitive.DateTime
	last_message_id primitive.ObjectID //(references the last message)
}
