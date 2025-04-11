package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"github.com/rs/zerolog"
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
	settingsService services.ISettingsService
	logger          *zerolog.Logger
}

func NewSettingsController(log *zerolog.Logger, settingsSvc services.ISettingsService) ISettingsController {
	return &SettingsController{
		settingsService: settingsSvc,
		logger:          log,
	}
}

func (s *SettingsController) CreateUserSettings(c *fiber.Ctx, userId string) error {
	// Convert userId to primitive.ObjectID
	objectID, err := utils.StringToObjectID(userId)
	if err != nil {
		s.logger.Error().Err(err).Str("userId", userId).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Failed to create userId for settings"))
	}

	// Create default UserSettings
	settings := models.GetSettingsDefaultsFromHeaders(utils.GetHeaderMap(c))
	settings.UserId = objectID

	// create to persist user settings in db
	_, err = s.settingsService.Create(c.Context(), settings)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create user settings")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to create user settings"))
	}
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(settings, "Created user settings"))
}

func (s *SettingsController) GetUserSettings(c *fiber.Ctx, userId string) error {
	return nil
}

func (s *SettingsController) UpdateUserSettings(c *fiber.Ctx, userId string) error {
	return nil
}
