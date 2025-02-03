package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	userId             primitive.ObjectID //(unique, primary key)
	firstName          string
	lastName           string
	email              string //(unique)
	phoneNumber        string //(optional)
	userType           string //(regular, enterprise)
	profilePicture     string //(URL to image)
	status             string //(online, offline, busy, etc.)
	bio                string
	languagePreference string
	timezone           string
	country            string
	createdAt          primitive.DateTime
	updatedAt          primitive.DateTime
}
