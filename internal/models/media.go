package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Media struct {
	mediaId         primitive.ObjectID //(unique, primary key)
	chatId          primitive.ObjectID //(references chats collection)
	senderId        primitive.ObjectID //(references users collection)
	mediaType       string             //(image, audio, video, document, etc.)
	fileName        string
	fileSize        int
	mediaUrl        string //(URL to the media file)
	uploadTimestamp primitive.DateTime
}
