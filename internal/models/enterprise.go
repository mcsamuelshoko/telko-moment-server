package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Enterprise struct {
	Id                 primitive.ObjectID
	name               string
	enterpriseOwner    []primitive.ObjectID
	websiteUrl         string
	verificationUrl    string
	verified           bool
	createdAt          primitive.Timestamp
	industry           string
	country            string
	headquarterCountry string
}
