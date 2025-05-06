package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Defined ChatMembership.Role constants
// for the ChatMembership Model
const (
	ChatMembershipRoleAdmin  = "admin"
	ChatMembershipRoleMember = "member"
	ChatMembershipRoleOwner  = "owner"
)

// Defined ChatMembership.Status constants
// for the ChatMembership Model
const (
	ChatMembershipStatusActive  = "active"
	ChatMembershipStatusBanned  = "banned"
	ChatMembershipStatusDeleted = "deleted"
)

// Defined ChatMembership.BanCode constants
// for the ChatMembership Model
const (
	ChatMembershipBanCodeNone      = 0
	ChatMembershipBanCodePermanent = 1
)

type ChatMembership struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ChatID    primitive.ObjectID `bson:"chatId"`
	UserID    primitive.ObjectID `bson:"userId"`
	Role      string             `bson:"role"` // e.g., "admin", "member","creator"
	JoinedAt  primitive.DateTime `bson:"joinedAt"`
	Status    string             `bson:"status"` // e.g., "active", "banned"
	BanCode   int                `bson:"banCode"`
	Blocked   bool               `bson:"blocked"`
	BlockedAt primitive.DateTime `bson:"blockedAt"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
}
