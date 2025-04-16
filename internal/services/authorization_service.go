package services

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/rs/zerolog"
)

// Define Action constants
const (
	ActionRead   = "read"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionPost   = "post" // e.g., post message, post reply
	ActionJoin   = "join"
	ActionLeave  = "leave"
	ActionBan    = "ban"
	ActionMute   = "mute"
	// ... add more specific actions as needed
)

// IAuthorizationService checks if a user can perform an action on a resource.
type IAuthorizationService interface {
	// Can method checks permission. resource can be the actual model struct (*models.ChatMessage)
	// or any struct containing the necessary attributes for the check.
	//Can(ctx context.Context, userId string, action string, resource interface{}) (bool, error)
	Can(ctx context.Context, user *models.User, resource interface{}, action string) (bool, error)
}

func NewAuthorizationService(log *zerolog.Logger, modelFilePath string, adapter persist.BatchAdapter) (IAuthorizationService, error) {

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(modelFilePath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %v", err)
	}

	return &casbinService{
		logger:   log,
		enforcer: enforcer,
	}, nil
}

// casbinService handles authorization using Casbin
type casbinService struct {
	logger   *zerolog.Logger
	enforcer *casbin.Enforcer
}

// User represents a user with attributes
//type User struct {
//	ID         primitive.ObjectID `bson:"_id,omitempty"`
//	Username   string             `bson:"username"`
//	Role       string             `bson:"role"`
//	Department string             `bson:"department"`
//	JoinedAt   time.Time          `bson:"joinedAt"`
//}
//
//// Resource represents a resource with attributes
//type Resource struct {
//	ID        primitive.ObjectID `bson:"_id,omitempty"`
//	Type      string             `bson:"type"`
//	OwnerID   primitive.ObjectID `bson:"ownerId"`
//	Public    bool               `bson:"public"`
//	CreatedAt time.Time          `bson:"createdAt"`
//}

// LoadPolicies loads ABAC policies
func (s *casbinService) LoadPolicies() error {
	// Department-based access for documents
	_, err := s.enforcer.AddPolicy("r.sub.Department == 'Engineering'", "engineering_docs", "read")
	if err != nil {
		return err
	}

	// Owner-based access
	_, err = s.enforcer.AddPolicy("r.sub.ID.Hex() == r.obj.UserId.Hex()", "any_resource", "edit")
	if err != nil {
		return err
	}

	// Admin role access
	_, err = s.enforcer.AddPolicy("r.sub.Role == 'admin'", "any_resource", "admin")
	if err != nil {
		return err
	}

	// Time-based access
	timeRule := "time.Since(r.sub.JoinedAt).Hours() > 24 * 30" // User joined more than a month ago
	_, err = s.enforcer.AddPolicy(timeRule, "premium_content", "access")
	if err != nil {
		return err
	}

	return nil
}

// Can CheckPermission checks if a user has permission to perform an action on a resource
func (s *casbinService) Can(ctx context.Context, user *models.User, resource interface{}, action string) (bool, error) {
	// For ABAC, we pass the actual objects rather than just IDs
	allowed, err := s.enforcer.Enforce(user, resource, action)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %v", err)
	}

	return allowed, nil
}

// AddPolicy adds a new policy rule
func (s *casbinService) AddPolicy(subRule string, obj string, act string) error {
	_, err := s.enforcer.AddPolicy(subRule, obj, act)
	return err
}

// RemovePolicy removes a policy rule
func (s *casbinService) RemovePolicy(subRule string, obj string, act string) error {
	_, err := s.enforcer.RemovePolicy(subRule, obj, act)
	return err
}
