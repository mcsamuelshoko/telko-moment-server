package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	messageId          primitive.ObjectID //(unique, primary key)
	chatId             primitive.ObjectID //(references chats collection)
	senderId           primitive.ObjectID //(references users collection)
	messageType        string             //(text, image, video, file, etc.)
	content            string             //(text content)
	mediaUrls          []string           //(for multimedia messages)
	timestamp          primitive.DateTime
	editedTimestamp    primitive.DateTime
	editedMessage      bool
	status             string               //(sent, delivered, read)
	mentions           []primitive.ObjectID //(array of user IDs mentioned in the message)
	repliedToMessageId primitive.ObjectID   //(optional, references another message_id for thread)
}
