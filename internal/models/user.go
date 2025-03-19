package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName          string             `json:"firstName" bson:"firstName"`
	LastName           string             `json:"lastName" bson:"lastName"`
	Username           string             `json:"username" bson:"username"`
	Password           string             `json:"password,omitempty" bson:"password"`
	Email              string             `json:"email" bson:"email"`
	PhoneNumber        string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	UserType           string             `json:"userType" bson:"userType"`
	ProfilePicture     string             `json:"profilePicture,omitempty" bson:"profilePicture,omitempty"`
	Status             string             `json:"status" bson:"status"`
	Bio                string             `json:"bio,omitempty" bson:"bio,omitempty"`
	LanguagePreference string             `json:"languagePreference,omitempty" bson:"languagePreference,omitempty"`
	Timezone           string             `json:"timezone,omitempty" bson:"timezone,omitempty"`
	Country            string             `json:"country,omitempty" bson:"country,omitempty"`
	CreatedAt          primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt          primitive.DateTime `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
