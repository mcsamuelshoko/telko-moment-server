package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id                 primitive.ObjectID //(unique, primary key)
	FirstName          string
	LastName           string
	Username           string
	Password           string
	Email              string //(unique)
	PhoneNumber        string //(optional)
	UserType           string //(regular, enterprise)
	ProfilePicture     string //(URL to image)
	Status             string //(online, offline, busy, etc.)
	Bio                string
	LanguagePreference string
	Timezone           string
	Country            string
	CreatedAt          primitive.DateTime
	UpdatedAt          primitive.DateTime
}
