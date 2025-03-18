package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h UserHandler) CreateUser(ctx *fiber.Ctx) error {
	user := &models.User{}

	err := h.userService.CreateUser(ctx.Context(), user)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(user)
}

func (h UserHandler) GetUser(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	user, err := h.userService.GetUserByUsername(ctx.Context(), username)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(user)
}

func (h UserHandler) GetUsers(ctx *fiber.Ctx) error {
	//users := []models.User{}
	users, err := h.userService.ListUsers(ctx.Context(), 0, 100)
	if err != nil {
		log.Error().Err(err).Msg("failed to get users")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(users)
}
