package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnterpriseRBAC struct {
	roleId        primitive.ObjectID //(unique, primary key)
	roleName      string             //(e.g., admin, user, manager)
	enterpriseId  primitive.ObjectID
	permissions   string               //(list of permissions like read, write, delete, etc.)
	assignedUsers []primitive.ObjectID //(array of user IDs who have this role)
	createdAt     primitive.DateTime
	updatedAt     primitive.DateTime
}
