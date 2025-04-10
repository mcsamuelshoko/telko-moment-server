package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/api"
	"github.com/mcsamuelshoko/telko-moment-server/internal/controllers"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"github.com/rs/zerolog"
)

type RoutesHandler struct {
	Logger             *zerolog.Logger
	UserController     controllers.IUserController
	SettingsController controllers.ISettingsController
	AuthController     controllers.IAuthenticationController
}

// NewRoutesHandler creates a new RoutesHandler instance.
func NewRoutesHandler(
	log *zerolog.Logger,
	userController controllers.IUserController,
	settingsController controllers.ISettingsController,
	authController controllers.IAuthenticationController,
) *RoutesHandler {

	return &RoutesHandler{
		Logger:             log,
		UserController:     userController,
		SettingsController: settingsController,
		AuthController:     authController,
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
		handler.Logger.Debug().Msg("[route]:GET:/")
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
		handler.Logger.Debug().Msg("[route]:GET:/")
		return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "Ready to serve!"))
	})

	// AUTH
	auth := v1.Group("/auth")
	auth.Post("/login", func(ctx *fiber.Ctx) error {
		return r.AuthController.Login(ctx)
	})
	auth.Post("/register", func(ctx *fiber.Ctx) error {
		return r.AuthController.Register(ctx)
	})
	auth.Post("/update-token", func(ctx *fiber.Ctx) error {
		return r.AuthController.UpdateRefreshToken(ctx)
	})
	auth.Post("/logout", func(ctx *fiber.Ctx) error {
		return r.AuthController.CancelRefreshToken(ctx)
	})

	// USERS
	users := v1.Group("/users")
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
	return r.SettingsController.GetUserSettings(c, userId)
}

func (r RoutesHandler) UpdateUserSettings(c *fiber.Ctx, userId string) error {
	return r.SettingsController.UpdateUserSettings(c, userId)
}

func (r RoutesHandler) GetAllUsers(c *fiber.Ctx) error {
	return r.UserController.GetAllUsers(c)
}

func (r RoutesHandler) CreateUser(c *fiber.Ctx) error {
	return r.UserController.CreateUser(c)
}

func (r RoutesHandler) DeleteUser(c *fiber.Ctx, userId string) error {
	return r.UserController.DeleteUser(c, userId)
}

func (r RoutesHandler) GetUserById(c *fiber.Ctx, userId string) error {
	return r.UserController.GetUserById(c, userId)
}

func (r RoutesHandler) UpdateUser(c *fiber.Ctx, userId string) error {
	return r.UserController.UpdateUser(c, userId)
}
