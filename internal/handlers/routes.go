package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/api"
	"github.com/mcsamuelshoko/telko-moment-server/internal/controllers"
	"github.com/mcsamuelshoko/telko-moment-server/internal/handlers/middleware"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"github.com/rs/zerolog"
)

type RoutesHandler struct {
	logger             *zerolog.Logger
	authMiddleware     middleware.IAuthMiddleware
	authCtxMiddleware  *middleware.AuthContextMiddleware
	userController     controllers.IUserController
	settingsController controllers.ISettingsController
	authController     controllers.IAuthenticationController
}

// NewRoutesHandler creates a new RoutesHandler instance.
func NewRoutesHandler(
	log *zerolog.Logger,
	authMiddleware middleware.IAuthMiddleware,
	authCtxMiddleware *middleware.AuthContextMiddleware,
	userController controllers.IUserController,
	settingsController controllers.ISettingsController,
	authController controllers.IAuthenticationController,
) *RoutesHandler {

	return &RoutesHandler{
		logger:             log,
		authMiddleware:     authMiddleware,
		userController:     userController,
		settingsController: settingsController,
		authController:     authController,
		authCtxMiddleware:  authCtxMiddleware,
	}
}

func (r RoutesHandler) SetupRoutes(app *fiber.App) {

	//pointing to handler out of laziness to refactor
	handler := &r

	// Wrap the handler with the generated server interface wrapper
	wrapper := api.ServerInterfaceWrapper{Handler: handler}
	//si := api.ServerInterface(handler)

	// ENTRY
	entry := app.Group("/")
	entry.Get("/", func(c *fiber.Ctx) error {
		r.logger.Debug().Msg("[route]:GET:/")
		return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.Map{"name": "Telko Moment API", "version": "0.1.0-alpha", "relativeLink": "/api/v1/"}, "Hello, World!"))
	})

	// API ROUTING SETUP
	_api := app.Group("/api")
	v1 := _api.Group("/v1")

	//// Register the routes
	//api.RegisterHandlers(v1, si)

	// ##################################################################### //
	// ENTRY
	v1.Get("/", func(c *fiber.Ctx) error {
		r.logger.Debug().Msg("[route]:GET:/")
		return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "Ready to serve!"))
	})

	// AUTH
	auth := v1.Group("/auth")
	auth.Post("/login", func(ctx *fiber.Ctx) error {
		return r.authController.Login(ctx)
	})
	auth.Post("/register", func(ctx *fiber.Ctx) error {
		return r.authController.Register(ctx)
	})
	auth.Post("/update-token", func(ctx *fiber.Ctx) error {
		return r.authController.UpdateRefreshToken(ctx)
	})
	auth.Post("/logout", func(ctx *fiber.Ctx) error {
		return r.authController.CancelRefreshToken(ctx)
	})

	// USERS
	users := v1.Group("/users")
	// use middleware that apply
	users.Use(r.authMiddleware.Authenticate())
	users.Use(r.authCtxMiddleware.AddUserContext())
	users.Get("/", wrapper.GetAllUsers)
	users.Get("/:id", wrapper.GetUserById)

	users.Post("/", wrapper.CreateUser)
	users.Put("/:id", wrapper.UpdateUser)
	users.Delete("/:id", wrapper.DeleteUser)

	// SETTINGS
	settings := v1.Group("/settings")
	settings.Get("/:id", wrapper.GetUserSettings)
	settings.Put("/:id", wrapper.UpdateUserSettings)

	// ... Setup routes for other resources (chats, messages, etc.)

}

func (r RoutesHandler) GetAllChatGroups(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) CreateChatGroup(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) DeleteChatGroup(c *fiber.Ctx, chatGroupId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetChatGroupById(c *fiber.Ctx, chatGroupId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) UpdateChatGroup(c *fiber.Ctx, chatGroupId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetAllChats(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) CreateChat(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) DeleteChat(c *fiber.Ctx, chatId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetChatById(c *fiber.Ctx, chatId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) UpdateChat(c *fiber.Ctx, chatId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetAllHighlights(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) CreateHighlight(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) DeleteHighlightById(c *fiber.Ctx, highlightId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetHighlightById(c *fiber.Ctx, highlightId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) UpdateHighlightById(c *fiber.Ctx, highlightId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) UploadMedia(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) DeleteMedia(c *fiber.Ctx, mediaId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetMediaById(c *fiber.Ctx, mediaId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetMessagesByChatId(c *fiber.Ctx, params api.GetMessagesByChatIdParams) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) SendMessage(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) DeleteMessage(c *fiber.Ctx, messageId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetMessageById(c *fiber.Ctx, messageId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) UpdateMessage(c *fiber.Ctx, messageId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RoutesHandler) GetUserSettings(c *fiber.Ctx, userId string) error {
	return r.settingsController.GetUserSettings(c, userId)
}

func (r RoutesHandler) UpdateUserSettings(c *fiber.Ctx, userId string) error {
	return r.settingsController.UpdateUserSettings(c, userId)
}

func (r RoutesHandler) GetAllUsers(c *fiber.Ctx) error {
	return r.userController.GetAllUsers(c)
}

func (r RoutesHandler) CreateUser(c *fiber.Ctx) error {
	return r.userController.CreateUser(c)
}

func (r RoutesHandler) DeleteUser(c *fiber.Ctx, userId string) error {
	return r.userController.DeleteUser(c, userId)
}

func (r RoutesHandler) GetUserById(c *fiber.Ctx, userId string) error {
	return r.userController.GetUserById(c, userId)
}

func (r RoutesHandler) UpdateUser(c *fiber.Ctx, userId string) error {
	return r.userController.UpdateUser(c, userId)
}
