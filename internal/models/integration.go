package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Integration struct {
	integrationId   primitive.ObjectID //(unique, primary key)
	name            string             //(e.g., Google Drive, Slack)
	userId          primitive.ObjectID //(references users collection)
	integrationType string             //(OAuth, Webhook, etc.)
	integrationData string             //(JSON storing keys, access tokens, etc.)
	status          string             //(active, inactive)
	createdAt       primitive.DateTime
	updatedAt       primitive.DateTime
}
