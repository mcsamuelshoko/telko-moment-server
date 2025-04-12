package services

import (
	"context"
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
	Can(ctx context.Context, userId string, action string, resource interface{}) (bool, error)
}

type AuthorizationService struct {
	Logger *zerolog.Logger
}

func NewAuthorizationService(Logger *zerolog.Logger) IAuthorizationService {
	return &AuthorizationService{
		Logger: Logger,
	}
}

func (a *AuthorizationService) Can(ctx context.Context, userId string, action string, resource interface{}) (bool, error) {
	//TODO implement me
	panic("implement me")
}
