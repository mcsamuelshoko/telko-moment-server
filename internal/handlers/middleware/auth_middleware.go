package middleware

import "github.com/gofiber/fiber/v2"

type IAuthMiddleware interface {
	Authenticate() fiber.Handler
}
