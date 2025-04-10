package controllers

import (
	"errors"
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
	iName           string
	authService     services.IAuthenticationService
	log             *zerolog.Logger
	userService     services.IUserService
	settingsService services.ISettingsService
	jwtService      pkgservices.IJWTService
}

func NewAuthController(log *zerolog.Logger, userSvc services.IUserService, authSvc services.IAuthenticationService, settingsSvc services.ISettingsService, jwtSvc pkgservices.IJWTService) IAuthenticationController {
	return &AuthenticationController{
		iName:           "AuthenticationController",
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
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(msg))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(fiber.Map{"refreshToken": refreshToken}, "Created Refresh Token"))
}

func (a *AuthenticationController) UpdateRefreshToken(c *fiber.Ctx) error {
	// used by logger
	const funcName = "UpdateRefreshToken"

	updateRequest := new(models.UpdateTokenRequest)
	err := c.BodyParser(updateRequest)
	if err != nil {
		a.log.Error().Interface(funcName, a.iName).Err(err).Msg("Failed to parse update-token request")
	}

	userID, err := a.authService.GetUserIDFromRefreshToken(c.Context(), updateRequest.RefreshToken) // Implement this function in your database layer.
	if err != nil {
		a.log.Error().Interface(funcName, a.iName).Err(err).Msg("Invalid or expired refresh token")
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid or expired refresh token"))
	}

	// Verify the refresh token's validity (e.g., check expiration, signature).
	if !a.jwtService.VerifyRefreshToken(updateRequest.RefreshToken) {
		a.log.Error().Interface(funcName, a.iName).Msg("Invalid refresh token")
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid refresh token"))
	}

	// Generate a new access token.
	accessToken, err := a.jwtService.GenerateAccessToken(userID) // Generate an access token using your utility function
	if err != nil {
		msg := "Failed to generate access token"
		a.log.Error().Interface(funcName, a.iName).Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(msg))
	}

	// Optionally, generate a new refresh token.
	newRefreshToken, err := a.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		a.log.Error().Interface(funcName, a.iName).Err(err).Msg("failed to generate new refresh token")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to generate new refresh token"))
	}

	// Update by replacing the old refresh token with the new one in the database.
	err = a.authService.UpdateUserRefreshToken(c.Context(), userID, newRefreshToken)
	if err != nil {
		a.log.Error().Interface(funcName, a.iName).Err(err).Msg("failed to save new refresh token")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to save new refresh token"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.Map{"accessToken": accessToken, "refreshToken": newRefreshToken}, "Token refreshed"))
}

func (a *AuthenticationController) CancelRefreshToken(c *fiber.Ctx) error {
	// used by logger
	const funcName = "CancelRefreshToken"

	logoutRequest := new(models.LogoutRequest)
	err := c.BodyParser(logoutRequest)
	if err != nil {
		a.log.Error().Interface(funcName, a.iName).Err(err).Msg("Failed to parse logout request")
	}

	// Delete the refresh token from the database.
	err = a.authService.DeleteRefreshToken(c.Context(), logoutRequest.RefreshToken) // Implement this function in your database layer.
	if err != nil {
		msg := "Failed to cancel refresh token"
		a.log.Error().Interface(funcName, a.iName).Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(msg))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "Token cancelled"))
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
	var regUserResponse fiber.Map
	var responseStatus int

	//register user using email
	regUserResponse, err, responseStatus = a.registerUsingEmail(c, *emailRegisterRequest, err2, user, failedRegErrMsg)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to register email-registration user")
		return c.Status(responseStatus).JSON(utils.ErrorResponse(err.Error()))
	}
	if regUserResponse != nil {
		a.log.Info().Msg("User registered successfully with Email")
		return c.Status(fiber.StatusCreated).JSON(regUserResponse)
	}
	//register user using phoneNumber
	regUserResponse, err, responseStatus = a.registrationUsingPhoneNumber(c, *registerRequest, err1, user, failedRegErrMsg)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to register phoneNumber-registration user")
		return c.Status(responseStatus).JSON(utils.ErrorResponse(err.Error()))
	}
	if regUserResponse != nil {
		a.log.Info().Msg("User registered successfully with Phone Number")
		return c.Status(fiber.StatusCreated).JSON(regUserResponse)
	}

	// This code should never be reached given the previous error checks
	return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Unexpected error"))
}

func (a *AuthenticationController) registerUsingEmail(c *fiber.Ctx, emailRegisterRequest models.RegisterRequestEmail, err2 error, user *models.User, failedRegErrMsg string) (fiber.Map, error, int) {

	var err error

	//check if email is valid
	if !(utils.IsValidEmail(emailRegisterRequest.Email)) {
		a.log.Error().Msg("Email format is invalid")
		return nil, errors.New("invalid email"), fiber.StatusBadRequest
	}
	//check if password is strong
	if !(utils.IsStrongPassword(emailRegisterRequest.Password)) {
		a.log.Error().Msg("Password is weak")
		return nil, errors.New("password is weak"), fiber.StatusBadRequest
	}

	//check if user exists
	_, err = a.userService.GetUserByEmail(c.Context(), emailRegisterRequest.Email)
	if err != nil {
		a.log.Info().Msg("User does not exist, & can be registered")
	} else {
		a.log.Error().Err(err).Msg("User already exists")
		return nil, err, fiber.StatusConflict
	}

	//register user using Email
	if err2 == nil {
		user.Email = emailRegisterRequest.Email
		user.Password, err = utils.HashPassword(emailRegisterRequest.Password)
		user.Username = emailRegisterRequest.Email
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to hash password")
			return nil, errors.New("server had an error"), fiber.StatusInternalServerError
		}

		createdUser, err := a.userService.CreateUser(c.Context(), user)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to create user")
			return nil, errors.New(failedRegErrMsg), fiber.StatusInternalServerError
		}

		err = a.createSettingsForUser(c, createdUser)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to create user settings")
			return nil, errors.New(failedRegErrMsg), fiber.StatusInternalServerError

		}

		// Return the created user without sensitive information
		return utils.SuccessResponse(
			createdUser.Sanitize(),
			"User registered successfully",
		), nil, fiber.StatusCreated
	}

	a.log.Error().Err(err2).Msg("Unexpected Error, Failed to register user with Email")
	return nil, errors.New("unexpected email registration error"), fiber.StatusInternalServerError

}

func (a *AuthenticationController) registrationUsingPhoneNumber(c *fiber.Ctx, registerRequest models.RegisterRequest, err1 error, user *models.User, failedRegErrMsg string) (fiber.Map, error, int) {
	var err error

	//check if user exists
	_, err = a.userService.GetUserByPhoneNumber(c.Context(), registerRequest.PhoneNumber)
	if err != nil {
		a.log.Info().Msg("User does not exist, & can be registered")
	} else {
		a.log.Error().Err(err).Msg("User already exists")
		return nil, err, fiber.StatusConflict
	}

	//check if PhoneNumber is valid
	if !(utils.IsValidPhoneNumber(registerRequest.PhoneNumber)) {
		a.log.Error().Msg("Phone number format is invalid")
		return nil, errors.New("invalid phone number"), fiber.StatusBadRequest
	}
	//check if password is strong
	if !(utils.IsStrongPassword(registerRequest.Password)) {
		a.log.Error().Msg("Password is weak")
		return nil, errors.New("password is weak, please make it min-chars=8 and include a [Number], & [special character], & [small letter], & [uppercase letter]"), fiber.StatusBadRequest
	}

	if err1 == nil {
		user.PhoneNumber = registerRequest.PhoneNumber
		user.Password, err = utils.HashPassword(registerRequest.Password)
		user.Username = registerRequest.PhoneNumber
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to hash password")
			return nil, errors.New("server had an error"), fiber.StatusInternalServerError
		}

		createdUser, err := a.userService.CreateUser(c.Context(), user)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to create user")
			return nil, errors.New(failedRegErrMsg), fiber.StatusInternalServerError
		}
		err = a.createSettingsForUser(c, createdUser)
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to create user settings")
			return nil, errors.New(failedRegErrMsg), fiber.StatusInternalServerError
		}

		// Return the created user without sensitive information
		return utils.SuccessResponse(
			createdUser.Sanitize(),
			"User registered successfully",
		), nil, fiber.StatusCreated
	}

	a.log.Error().Err(err1).Msg("Unexpected Error, Failed to register user with Phone Number")
	return nil, errors.New("unexpected phone number registration error"), fiber.StatusInternalServerError
}

func (a *AuthenticationController) createSettingsForUser(c *fiber.Ctx, createdUser *models.User) error {
	settings := models.GetSettingsDefaultsFromHeaders(utils.GetHeaderMap(c))
	settings.UserId = createdUser.ID
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
