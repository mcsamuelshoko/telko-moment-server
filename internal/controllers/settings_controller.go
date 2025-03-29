package controllers

import "github.com/gofiber/fiber/v2"

type ISettingsController interface {

	// GetUserSettings get user settings
	// (GET /settings/{userId}
	GetUserSettings(c *fiber.Ctx, userId string) error

	// UpdateUserSettings Update user settings
	// (PUT /settings/{userId})
	UpdateUserSettings(c *fiber.Ctx, userId string) error
}

type SettingsController struct {
}

func NewSettingsController() ISettingsController {
	return &SettingsController{}
}

func (s *SettingsController) GetUserSettings(c *fiber.Ctx, userId string) error {
	return nil
}

func (s *SettingsController) UpdateUserSettings(c *fiber.Ctx, userId string) error {
	return nil
}
