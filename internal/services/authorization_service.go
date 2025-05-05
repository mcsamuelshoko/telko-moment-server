package services

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/rs/zerolog"
)

// Define Effect constants
const (
	EffectAllow = "allow"
	EffectDeny  = "deny"
)

// Define Action constants
const (
	ActionRead   = "read"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"

	ActionPost  = "post" // e.g., post message, post reply
	ActionJoin  = "join"
	ActionLeave = "leave"
	ActionBan   = "ban"
	ActionMute  = "mute"

	ActionModerator = "moderator"
	ActionAdmin     = "admin"
	ActionOwner     = "owner"
	// ... add more specific actions as needed
)

// Define Collection/Resource constants
const (
	ResourceUsers      = "users"
	ResourceSettings   = "settings"
	ResourceMessages   = "messages"
	ResourceChatGroups = "chat_groups"
)

type UserWrapper struct {
	ID string
}

type ObjectWrapper struct {
	UserId string
}

// IAuthorizationService checks if a user can perform an action on a resource.
type IAuthorizationService interface {
	// Can method checks permission. resource can be the actual model struct (*models.ChatMessage)
	// or any struct containing the necessary attributes for the check.
	//Can(ctx context.Context, userId string, action string, resource interface{}) (bool, error)
	Can(ctx context.Context, user *models.User, resource interface{}, action string) (bool, error)
	// LoadPolicies adds policies to the adapter. it adds the default policies defined by code when called at runtime
	LoadPolicies() error
}

// casbinAuthorizationService handles authorization using Casbin
type casbinAuthorizationService struct {
	iName    string
	logger   *zerolog.Logger
	enforcer *casbin.Enforcer
}

func NewCasbinAuthorizationService(log *zerolog.Logger, modelFilePath string, adapter persist.BatchAdapter) (IAuthorizationService, error) {

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(modelFilePath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %v", err)
	}

	// Load Policies
	enforcer.EnableAutoSave(true)

	return &casbinAuthorizationService{
		iName:    "CasbinAuthorizationService",
		logger:   log,
		enforcer: enforcer,
	}, nil
}

// LoadPolicies loads ABAC policies
func (s *casbinAuthorizationService) LoadPolicies() error {
	const kName = "LoadPolicies"

	var err error
	const failedToAddErrMsg = "failed to add policy ::"
	// Department-based access for documents
	//_, err := s.enforcer.AddPolicy("r.sub.Department == 'Engineering'", "engineering_docs", "read")
	//if err != nil {
	//	return err
	//}
	s.logger.Debug().Interface(kName, s.iName).Msg("starting policy additions")
	// User based access
	const userBasedSubrule = "r.sub.ID == r.obj.ID" //"r.sub.ID.Hex() == r.obj.ID.Hex()"
	// :::: Users Collection
	_, err = s.enforcer.AddPolicy(userBasedSubrule, ResourceUsers, ActionRead, EffectAllow)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg(failedToAddErrMsg + ResourceUsers + "-" + ActionRead)
		return err
	}
	_, err = s.enforcer.AddPolicy(userBasedSubrule, ResourceUsers, ActionUpdate, EffectAllow)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg(failedToAddErrMsg + ResourceUsers + "-" + ActionUpdate)
		return err
	}
	_, err = s.enforcer.AddPolicy(userBasedSubrule, ResourceUsers, ActionDelete, EffectAllow)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg(failedToAddErrMsg + ResourceUsers + "-" + ActionDelete)
		return err
	}

	// Owner-based access
	const ownerBasedSubrule = "r.sub.ID == r.obj.UserId"
	// :::: Settings Collection
	_, err = s.enforcer.AddPolicy(ownerBasedSubrule, ResourceSettings, ActionRead, EffectAllow)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg(failedToAddErrMsg + ResourceSettings + "-" + ActionRead)
		return err
	}

	_, err = s.enforcer.AddPolicy(ownerBasedSubrule, ResourceSettings, ActionUpdate, EffectAllow)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg(failedToAddErrMsg + ResourceSettings + "-" + ActionUpdate)
		return err
	}

	// Admin role access
	const adminBasedSubrule = "r.sub.ID in r.obj.AdminIDs"
	// :::: ChatGroups Collection
	_, err = s.enforcer.AddPolicy(adminBasedSubrule, ResourceChatGroups, ActionAdmin, EffectAllow)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg(failedToAddErrMsg + ResourceChatGroups + "-" + ActionAdmin)
		return err
	}

	//_, err = s.enforcer.AddPolicy("r.sub.Role == 'admin'", "any_resource", "admin")
	//if err != nil {
	//	return err
	//}

	// Time-based access
	//timeRule := "time.Since(r.sub.JoinedAt).Hours() > 24 * 30" // User joined more than a month ago
	//_, err = s.enforcer.AddPolicy(timeRule, "premium_content", "access")
	//if err != nil {
	//	return err
	//}
	s.logger.Debug().Interface(kName, s.iName).Msg("finished policy additions")
	return nil
}

// Can CheckPermission checks if a user has permission to perform an action on a resource
func (s *casbinAuthorizationService) Can(ctx context.Context, user *models.User, resource interface{}, action string) (bool, error) {
	const kName = "Can"

	s.logger.Debug().Interface(kName, s.iName).Msg(action + " :: on user ID " + user.ID.Hex())

	// For ABAC, we pass the actual objects rather than just IDs
	allowed, err := s.enforcer.Enforce(UserWrapper{ID: user.ID.Hex()}, resource, action)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %v", err)
	}

	return allowed, nil
}

// AddPolicy adds a new policy rule
func (s *casbinAuthorizationService) AddPolicy(subRule string, obj string, act string) error {
	_, err := s.enforcer.AddPolicy(subRule, obj, act)
	return err
}

// RemovePolicy removes a policy rule
func (s *casbinAuthorizationService) RemovePolicy(subRule string, obj string, act string) error {
	_, err := s.enforcer.RemovePolicy(subRule, obj, act)
	return err
}
