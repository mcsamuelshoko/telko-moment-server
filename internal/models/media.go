package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Media struct {
	Id              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ChatId          primitive.ObjectID `json:"chatId" bson:"chatId"`
	SenderId        primitive.ObjectID `json:"senderId" bson:"senderId"`
	MediaType       string             `json:"mediaType" bson:"mediaType"`
	FileName        string             `json:"fileName" bson:"fileName"`
	FileSize        int                `json:"fileSize" bson:"fileSize"`
	MediaUrl        string             `json:"mediaUrl" bson:"mediaUrl"`
	UploadTimestamp primitive.DateTime `json:"uploadTimestamp" bson:"uploadTimestamp"`
}
