package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Call struct {
	callId           primitive.ObjectID   //(unique, primary key)
	callerId         primitive.ObjectID   //(references users collection)
	receiverId       []primitive.ObjectID //(references users collection)
	callType         string               //(audio, video)
	callStatus       string               //(initiated, ongoing, ended, missed)
	startTime        primitive.DateTime
	endTime          primitive.DateTime //(optional)
	duration         int                //(seconds)
	callRecordingUrl string             //(optional, for saved calls)
}
