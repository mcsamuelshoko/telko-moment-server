package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"github.com/rs/zerolog"
)

// ContextKey is a custom type for context keys
type ContextKey string

const (
	UserObjectContextKey ContextKey = "userObject"
	UserIDContextKey     ContextKey = "userID" // Use primitive.ObjectID
)

type AuthContextMiddleware struct {
	iName    string
	userRepo repository.IUserRepository
	logger   *zerolog.Logger // logger
}

func NewAuthContextMiddleware(log *zerolog.Logger, userRepo repository.IUserRepository) *AuthContextMiddleware {
	return &AuthContextMiddleware{
		iName:    "AuthContextMiddleware",
		userRepo: userRepo,
		logger:   log,
	}
}

func (acm *AuthContextMiddleware) AddUserContext() fiber.Handler {
	const kName = "AddUserContext"

	return func(c *fiber.Ctx) error {
		acm.logger.Debug().Interface(kName, acm.iName).Msg("adding user context")

		userIDStr, ok := c.Context().Value(UserIDStrContextKey).(string)
		if !ok || userIDStr == "" {
			//http.Error(w, "Unauthorized: Missing user identifier", http.StatusUnauthorized)
			acm.logger.Error().Interface(kName, acm.iName).Msg("Invalid user id from context")
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Missing user identifier"))
		}

		//userID, err := primitive.ObjectIDFromHex(userIDStr)
		//if err != nil {
		//	//http.Error(w, "Unauthorized: Invalid user identifier", http.StatusUnauthorized)
		//	acm.logger.Error().Err(err).Msg("Invalid user id hex() string value from context")
		//	return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid user id hex string"))
		//}

		// Fetch the user object
		user, err := acm.userRepo.GetByID(c.Context(), userIDStr)
		if err != nil {
			if errors.Is(err, err) { // Use specific errors
				//http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
				acm.logger.Error().Interface(kName, acm.iName).Err(err).Msg("Unauthorized: User not found")
				return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("User not found"))
			} else {
				// Log the actual error
				//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				acm.logger.Error().Interface(kName, acm.iName).Err(err).Msg("Error while getting user")
				return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Error while getting user"))
			}
		}

		// Add user object and validated userID to context
		//ctx := context.WithValue(r.Context(), UserObjectContextKey, user)
		//ctx = context.WithValue(ctx, UserIDContextKey, user.ID) // Add the ObjectID
		c.Locals(UserObjectContextKey, user)
		//c.Locals(UserIDContextKey, userID)

		//next.ServeHTTP(w, r.WithContext(ctx))
		return c.Next()
	}
}
