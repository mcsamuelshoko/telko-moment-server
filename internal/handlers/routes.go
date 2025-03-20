package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mcsamuelshoko/telko-moment-server/internal/controllers"
)

func SetupRoutes(app *fiber.App, userCtrl *controllers.UserController) {

	// API ROUTING SETUP
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// ENTRY
	entry := app.Group("/")
	entry.Get("/", func(c *fiber.Ctx) error {
		log.Debug("GET:/")
		return c.SendString("Hello, World!")
	})

	// USERS
	users := v1.Group("/users")
	users.Get("/", userCtrl.GetUser)

}
