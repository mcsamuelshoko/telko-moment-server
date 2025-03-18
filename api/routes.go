package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/internal/handlers"
)

func SetupRoutes(app *fiber.App, handlers *handlers.UserHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	entry := app.Group("/")
	entry.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	users := v1.Group("/users")
	users.Post("/", handlers.CreateUser)
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)

	// other routes...

}
