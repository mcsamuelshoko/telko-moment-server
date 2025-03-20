package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	"github.com/rs/zerolog"
)

type UserController struct {
	log     *zerolog.Logger
	service *services.UserService
}

func NewUserController(log *zerolog.Logger, service *services.UserService) *UserController {
	return &UserController{
		service: service,
		log:     log,
	}
}

func (ctrl *UserController) GetUser(c *fiber.Ctx) error {
	user, err := ctrl.service.GetUserByID(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(500).SendString("Error")
	}
	return c.JSON(fiber.Map{"name": user})
}
