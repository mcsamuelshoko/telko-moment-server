package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"github.com/rs/zerolog"
)

type IUserController interface {

	// CreateUser Create a new user
	// (POST /users)
	CreateUser(c *fiber.Ctx) error

	// GetAllUsers Get all users
	// (GET /users)
	GetAllUsers(c *fiber.Ctx) error

	// GetUserById Get a user by ID
	// (GET /users/{userId})
	GetUserById(c *fiber.Ctx, userId string) error

	// UpdateUser Update a user
	// (PUT /users/{userId})
	UpdateUser(c *fiber.Ctx, userId string) error

	// DeleteUser Delete a user
	// (DELETE /users/{userId})
	DeleteUser(c *fiber.Ctx, userId string) error
}

type UserController struct {
	log             *zerolog.Logger
	userService     services.IUserService
	settingsService services.ISettingsService
}

//var _ api.ServerInterface = (*UserController)(nil)

func NewUserController(log *zerolog.Logger, service services.IUserService, settingsSvc services.ISettingsService) IUserController {
	return &UserController{
		userService:     service,
		log:             log,
		settingsService: settingsSvc,
	}
}

func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		ctrl.log.Error().Err(err).Msg("Failed to parse user body")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	// Hash user password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		ctrl.log.Error().Err(err).Msg("Failed to hash password")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to create user"))
	} else {
		user.Password = hashedPassword
	}

	// Create User in DB
	createdUser, err := ctrl.userService.CreateUser(c.Context(), user)
	if err != nil {
		msg := "Failed to create user"
		ctrl.log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(msg))
	}

	// Create Default UserSettings
	settings := models.GetSettingsDefaultsFromHeaders(utils.GetHeaderMap(c)) // Using a method in models package.
	settings.UserId = createdUser.ID

	_, err = ctrl.settingsService.Create(c.Context(), settings)
	if err != nil {
		ctrl.log.Error().Err(err).Msg("Failed to create user settings")
		// Handle settings creation error, perhaps delete the user that was created.
		// Rollback user creation if settings creation fails.
		deleteErr := ctrl.userService.DeleteUser(c.Context(), createdUser.ID.String())
		if deleteErr != nil {
			ctrl.log.Error().Err(deleteErr).Msg("Failed to delete user after settings creation error")
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to create user settings"))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(createdUser.Sanitize(), "User created"))
}

func (ctrl *UserController) GetAllUsers(c *fiber.Ctx) error {
	//TODO: make sure the page and the limit come from the request and not solid values
	users, err := ctrl.userService.ListUsers(c.Context(), 0, 50)
	if err != nil {
		ctrl.log.Error().Err(err).Msg("Failed to get list of users")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to get list of user"))
	}
	// sanitizeUsers
	sanitizedUsers := make([]map[string]interface{}, 0)
	for _, user := range users {
		sanitizedUsers = append(sanitizedUsers, user.Sanitize())
	}
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(sanitizedUsers, "Users Listed"))
}

func (ctrl *UserController) GetUserById(c *fiber.Ctx, userId string) error {
	user, err := ctrl.userService.GetUserByID(c.Context(), userId)
	if err != nil {
		ctrl.log.Error().Err(err).Str("userID", userId).Msg("Failed to get user")
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse("Failed to get user"))
	}
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(user.Sanitize(), "User found")) // Return the user object directly
}

func (ctrl *UserController) UpdateUser(c *fiber.Ctx, userId string) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		ctrl.log.Error().Err(err).Msg("Failed to parse user body")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	updatedUser, err := ctrl.userService.UpdateUser(c.Context(), user)
	if err != nil {
		ctrl.log.Error().Err(err).Str("userID", userId).Msg("Failed to update user")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to update user"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(updatedUser.Sanitize(), "User updated"))
}

func (ctrl *UserController) DeleteUser(c *fiber.Ctx, userId string) error {
	err := ctrl.userService.DeleteUser(c.Context(), userId)
	if err != nil {
		ctrl.log.Error().Err(err).Str("userID", userId).Msg("Failed to delete user")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to delete user"))
	}

	return c.Status(fiber.StatusNoContent).JSON(utils.SuccessResponse(nil, "User deleted"))
}
