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
	iName              string
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
		iName:              "RoutesHandler",
		logger:             log,
		authMiddleware:     authMiddleware,
		userController:     userController,
		settingsController: settingsController,
		authController:     authController,
		authCtxMiddleware:  authCtxMiddleware,
	}
}

func (r RoutesHandler) SetupRoutes(app *fiber.App) {
	const kName = "SetupRoutes"

	//pointing to handler out of laziness to refactor
	handler := &r

	// Wrap the handler with the generated server interface wrapper
	wrapper := api.ServerInterfaceWrapper{Handler: handler}
	//si := api.ServerInterface(handler)

	// ::: ENTRY
	entry := app.Group("/")
	entry.Get("/", func(c *fiber.Ctx) error {
		r.logger.Debug().Msg("[route]:GET:/")
		return c.Status(fiber.StatusOK).
			JSON(utils.SuccessResponse(
				fiber.Map{
					"name":         "Telko Moment API",
					"version":      "0.1.0-alpha",
					"relativeLink": "/api/v1/"},
				"Hello, World!"))
	})

	// ::: API ROUTING SETUP
	apiRoute := app.Group("/api")
	v1 := apiRoute.Group("/v1")

	//// Register the routes
	//api.RegisterHandlers(v1, si)

	// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: //

	// ::: ENTRY
	v1.Get("/", wrapper.Index)

	// ::: AUTH
	auth := v1.Group("/auth")
	auth.Post("/login", wrapper.AuthLogin)
	auth.Post("/register", wrapper.AuthRegister)
	auth.Post("/refresh-token", func(ctx *fiber.Ctx) error {
		return r.authController.UpdateRefreshToken(ctx)
	})
	auth.Post("/logout", func(ctx *fiber.Ctx) error {
		return r.authController.CancelRefreshToken(ctx)
	})

	// ::: USERS
	users := v1.Group("/users")

	// use middleware that apply
	users.Use(r.authMiddleware.Authenticate()).Use(r.authCtxMiddleware.AddUserContext())

	// add routes that use the middleware
	users.Get("/", wrapper.GetAllUsers)
	users.Get("/:userId", wrapper.GetUserById)

	users.Post("/", wrapper.CreateUser)
	users.Put("/:userId", wrapper.UpdateUser)
	users.Delete("/:userId", wrapper.DeleteUser)

	// ::: SETTINGS
	settings := v1.Group("/settings")
	// use middleware that apply
	settings.Use(r.authMiddleware.Authenticate()).Use(r.authCtxMiddleware.AddUserContext())
	settings.Get("/:userId", wrapper.GetUserSettings)
	settings.Put("/:userId", wrapper.UpdateUserSettings)

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

// ::::::::::::::::::::::::::::::  ROUTES ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// :::: ENTRY -or- INDEX

func (r RoutesHandler) Index(c *fiber.Ctx) error {
	const kName = "Index"

	r.logger.Debug().Interface(kName, r.iName).Msg("[route]:GET:/api/v1/")
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "Ready to serve!"))

}

// :::: AUTH

func (r RoutesHandler) AuthLogin(c *fiber.Ctx) error {
	return r.authController.Login(c)
}

func (r RoutesHandler) AuthRegister(c *fiber.Ctx) error {
	return r.authController.Register(c)
}

// :::: SETTINGS

func (r RoutesHandler) GetUserSettings(c *fiber.Ctx, userId string) error {
	return r.settingsController.GetUserSettings(c, userId)
}

func (r RoutesHandler) UpdateUserSettings(c *fiber.Ctx, userId string) error {
	return r.settingsController.UpdateUserSettings(c, userId)
}

// :::: USERS

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
