package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscription struct {
	subscriptionId primitive.ObjectID //(unique, primary key)
	userId         primitive.ObjectID //(references users collection)
	planType       string             //(Free, Basic, Premium, Enterprise)
	startDate      primitive.DateTime
	endDate        primitive.DateTime //(optional)
	status         string             //(active, expired, canceled)
	paymentMethod  string             //(e.g., credit card, PayPal)
	paymentStatus  string             //(paid, pending, failed)
	createdAt      primitive.DateTime
	updatedAt      primitive.DateTime
}
