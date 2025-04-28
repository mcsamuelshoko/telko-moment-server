package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/internal/handlers/middleware"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"github.com/rs/zerolog"
	"time"
)

type ISettingsController interface {

	// GetUserSettings get user settings
	// (GET /settings/{userId}
	GetUserSettings(c *fiber.Ctx, userId string) error

	// UpdateUserSettings Update user settings
	// (PUT /settings/{userId})
	UpdateUserSettings(c *fiber.Ctx, userId string) error

	// CreateUserSettings creates settings for a user
	CreateUserSettings(c *fiber.Ctx, userId string) error
}

type SettingsController struct {
	iName                string
	settingsService      services.ISettingsService
	authorizationService services.IAuthorizationService
	logger               *zerolog.Logger
}

func NewSettingsController(log *zerolog.Logger, settingsSvc services.ISettingsService, authznSvc services.IAuthorizationService) ISettingsController {
	return &SettingsController{
		iName:                "SettingsController",
		settingsService:      settingsSvc,
		authorizationService: authznSvc,
		logger:               log,
	}
}

func (s *SettingsController) CreateUserSettings(c *fiber.Ctx, userId string) error {
	const kName = "CreateUserSettings"
	// Convert userId to primitive.ObjectID
	objectID, err := utils.StringToObjectID(userId)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Str("userId", userId).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Failed to create userId for settings"))
	}

	// Create default UserSettings
	settings := models.GetSettingsDefaultsFromHeaders(utils.GetHeaderMap(c))
	settings.UserId = objectID

	// create to persist user settings in db
	_, err = s.settingsService.Create(c.Context(), settings)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg("Failed to create user settings")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to create user settings"))
	}
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(settings, "Created user settings"))
}

func (s *SettingsController) GetUserSettings(c *fiber.Ctx, userId string) error {
	const kName = "GetUserSettings"

	// Authorize
	can, status, response, err := s.isAuthorizedForSettingsResource(c, userId, services.ActionRead)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg("Failed to authorize for Settings-Resource")
		return c.Status(status).JSON(response)
	}

	if can {
		userSettings, err := s.settingsService.GetByUserId(c.Context(), userId)
		if err != nil {
			s.logger.Error().Interface(kName, s.iName).Err(err).Msg("Failed to get user settings")
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse("Could not find user settings"))
		}
		return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(userSettings, "User settings"))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to get user settings, unexpected error occurred"))
}

func (s *SettingsController) UpdateUserSettings(c *fiber.Ctx, userId string) error {
	const kName = "UpdateUserSettings"

	can, status, response, err := s.isAuthorizedForSettingsResource(c, userId, services.ActionUpdate)
	if err != nil {
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg("Failed to authorize for Settings-Resource")
		return c.Status(status).JSON(response)
	}
	if can {
		settingsUpdate := new(models.Settings)
		err := c.BodyParser(settingsUpdate)
		if err != nil {
			s.logger.Error().Interface(kName, s.iName).Err(err).Msg("Failed to retrieve user settings from request body")
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Could not parse request body for settings update"))
		}

		settingsUpdate.UpdatedAt = time.Now()
		err = s.settingsService.Update(c.Context(), settingsUpdate)
		if err != nil {
			s.logger.Error().Interface(kName, s.iName).Err(err).Msg("Failed to update user settings")
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to update user settings"))
		}

		return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(settingsUpdate, "Updated user settings"))
	}

	s.logger.Error().Interface(kName, s.iName).Err(err).Msg("Failed to update user settings due to unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to update user settings due to unexpected error"))
}

func (s *SettingsController) isAuthorizedForSettingsResource(c *fiber.Ctx, userId string, action string) (bool, int, fiber.Map, error) {
	const kName = "isAuthorizedForSettingsResource"

	user, ok := c.Context().Value(middleware.UserObjectContextKey).(*models.User)
	if !ok || user.ID.Hex() != userId {
		msg := "Failed to get user object from context"
		s.logger.Error().Interface(kName, s.iName).Msg(msg)
		return false, fiber.StatusInternalServerError, utils.ErrorResponse("Could not determine user context"), errors.New(msg)
	}

	settings, err := s.settingsService.GetByUserId(c.Context(), userId)
	if err != nil {
		msg := "Failed to get user settings"
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg(msg)
		return false, fiber.StatusNotFound, utils.ErrorResponse("Could not find user settings"), errors.New(msg)
	}
	//settings.ID.Hex()
	can, err := s.authorizationService.Can(c.Context(), user, settings, action)
	if err != nil {
		msg := "Failed to " + action + " user settings due to missing permissions"
		s.logger.Error().Interface(kName, s.iName).Err(err).Msg(msg)
		return false,
			fiber.StatusUnauthorized,
			utils.ErrorResponse("Failed to " + action + " user settings due to missing permissions"),
			errors.New(msg)
	}

	return can, fiber.StatusOK, nil, nil
}
