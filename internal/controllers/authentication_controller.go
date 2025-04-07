package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	pkgservices "github.com/mcsamuelshoko/telko-moment-server/pkg/services"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"

	"github.com/rs/zerolog"
)

type IAuthenticationController interface {
	CreateRefreshToken(c *fiber.Ctx) error
	UpdateRefreshToken(c *fiber.Ctx) error
	CancelRefreshToken(c *fiber.Ctx) error

	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

type AuthenticationController struct {
	authService     services.IAuthenticationService
	log             *zerolog.Logger
	userService     services.IUserService
	settingsService services.ISettingsService
	jwtService      pkgservices.IJWTService
}

func NewAuthController(log *zerolog.Logger, userSvc services.IUserService, authSvc services.IAuthenticationService, settingsSvc services.ISettingsService, jwtSvc pkgservices.IJWTService) IAuthenticationController {
	return &AuthenticationController{
		log:             log,
		authService:     authSvc,
		settingsService: settingsSvc,
		userService:     userSvc,
		jwtService:      jwtSvc,
	}
}

func (a *AuthenticationController) CreateRefreshToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token) // Get the user from the context (assuming you have middleware)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string) // Assuming "sub" claim holds the user ID.

	refreshToken, err := a.jwtService.GenerateRefreshToken(userID) // Generate a refresh token using your utility function
	if err != nil {
		msg := "Failed to generate refresh token"
		a.log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg})
	}

	// Store the refresh token in the database, associating it with the user.
	err = a.authService.SaveRefreshToken(c.Context(), userID, refreshToken) // Implement this function in your database layer.
	if err != nil {
		msg := "Failed to save refresh token"
		a.log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg})
	}

	return c.JSON(fiber.Map{"refreshToken": refreshToken})
}

func (a *AuthenticationController) UpdateRefreshToken(c *fiber.Ctx) error {
	refreshToken := c.FormValue("refreshToken")

	userID, err := a.authService.GetUserIDFromRefreshToken(c.Context(), refreshToken) // Implement this function in your database layer.
	if err != nil {
		a.log.Error().Err(err).Msg("Invalid or expired refresh token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired refresh token"})
	}

	// Verify the refresh token's validity (e.g., check expiration, signature).
	if !a.jwtService.VerifyRefreshToken(refreshToken) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	// Generate a new access token.
	accessToken, err := a.jwtService.GenerateAccessToken(userID) // Generate an access token using your utility function
	if err != nil {
		msg := "Failed to generate access token"
		a.log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg})
	}

	// Optionally, generate a new refresh token.
	newRefreshToken, err := a.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		a.log.Error().Err(err).Msg("failed to generate new refresh token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate new refresh token"})
	}

	// Update by replacing the old refresh token with the new one in the database.
	err = a.authService.UpdateUserRefreshToken(c.Context(), userID, newRefreshToken)
	if err != nil {
		a.log.Error().Err(err).Msg("failed to save new refresh token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save new refresh token"})
	}

	return c.JSON(fiber.Map{"accessToken": accessToken, "refreshToken": newRefreshToken})
}

func (a *AuthenticationController) CancelRefreshToken(c *fiber.Ctx) error {
	refreshToken := c.FormValue("refreshToken")

	// Delete the refresh token from the database.
	err := a.authService.DeleteRefreshToken(c.Context(), refreshToken) // Implement this function in your database layer.
	if err != nil {
		msg := "Failed to cancel refresh token"
		a.log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (a *AuthenticationController) Login(c *fiber.Ctx) error {
	loginRequest := new(models.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		a.log.Error().Err(err).Msg("Failed to parse login request")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	user, err := a.userService.GetUserByUsername(c.Context(), loginRequest.Username)
	if err != nil {
		a.log.Error().Err(err).Str("username", loginRequest.Username).Msg("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid credentials"))
	}

	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		a.log.Error().Str("username", loginRequest.Username).Msg("Invalid password")
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid credentials"))
	}

	accessToken, err := a.jwtService.GenerateAccessToken(user.ID.String())
	if err != nil {
		msg := "Failed to generate access token"
		a.log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(msg))
	}

	refreshToken, err := a.jwtService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		msg := "Failed to generate refresh token"
		a.log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(msg))
	}

	err = a.authService.SaveRefreshToken(c.Context(), user.ID.String(), refreshToken)
	if err != nil {
		msg := "Failed to save refresh token"
		a.log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(msg))
	}
	var data map[string]interface{} = make(map[string]interface{})
	data["accessToken"] = accessToken
	data["refreshToken"] = refreshToken

	msg := "Login successful"

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(data, msg))
}

func (a *AuthenticationController) Register(c *fiber.Ctx) error {
	// get data from dto
	registerRequest := new(models.RegisterRequest)
	emailRegisterRequest := new(models.RegisterRequestEmail)

	failedRegErrMsg := "Failed to register user"

	err1 := c.BodyParser(registerRequest)
	if err1 != nil {
		a.log.Error().Err(err1).Msg("Failed to parse register request")
	}
	err2 := c.BodyParser(emailRegisterRequest)
	if err2 != nil {
		a.log.Error().Err(err2).Msg("Failed to parse register request")
	}
	if err1 != nil && err2 != nil {
		a.log.Error().Err(err1).Err(err2).Msg("Failed to identify register request type")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	// user defaults
	user := models.GetUserDefaultsFromHeaders(utils.GetHeaderMap(c))
	var err error

	//register user using email
	err = a.registerUsingEmail(c, *emailRegisterRequest, err2, user, failedRegErrMsg)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to register email-registration user")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(failedRegErrMsg))
	}
	//register user using phoneNumber
	err = a.registrationUsingPhoneNumber(c, *registerRequest, err1, user, failedRegErrMsg)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to register phoneNumber-registration user")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(failedRegErrMsg))
	}

	// This code should never be reached given the previous error checks
	return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Unexpected error"))
}

func (a *AuthenticationController) registerUsingEmail(c *fiber.Ctx, emailRegisterRequest models.RegisterRequestEmail, err2 error, user *models.User, failedRegErrMsg string) error {

	var err error

	//register user using Email
	if err2 == nil {
		user.Email = emailRegisterRequest.Email
		user.Password, err = utils.HashPassword(emailRegisterRequest.Password)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to hash password")
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(failedRegErrMsg))
		}

		createdUser, err := a.userService.CreateUser(c.Context(), user)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to create user")
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(failedRegErrMsg))
		}
		err = a.createSettingsForUser(c, createdUser)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to create user settings")
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("failed to create user settings"))
		}

		// Return the created user without sensitive information
		return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(
			createdUser.Sanitize(),
			"User registered successfully",
		))
	}

	return nil

}

func (a *AuthenticationController) registrationUsingPhoneNumber(c *fiber.Ctx, registerRequest models.RegisterRequest, err1 error, user *models.User, failedRegErrMsg string) error {
	var err error
	if err1 == nil {
		user.PhoneNumber = registerRequest.PhoneNumber
		user.Password, err = utils.HashPassword(registerRequest.Password)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to hash password")
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(failedRegErrMsg))
		}

		createdUser, err := a.userService.CreateUser(c.Context(), user)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to create user")
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("failed to create user account"))
		}
		err = a.createSettingsForUser(c, createdUser)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to create user settings")
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("failed to create user settings"))
		}

		// Return the created user without sensitive information
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
			"user":    createdUser.Sanitize(),
		})
	}
	return nil
}

func (a *AuthenticationController) createSettingsForUser(c *fiber.Ctx, createdUser *models.User) error {
	settings := models.GetSettingsDefaultsFromHeaders(utils.GetHeaderMap(c))
	_, err := a.settingsService.Create(c.Context(), settings)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create user's settings")
		// Handle settings creation error, perhaps delete the user that was created.
		// Rollback user creation if settings creation fails.
		deleteErr := a.userService.DeleteUser(c.Context(), createdUser.ID.String())
		if deleteErr != nil {
			a.log.Error().Err(deleteErr).Msg("Failed to delete user after settings creation error")
		}
		return err
	}
	return nil
}
