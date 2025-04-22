package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/services"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"github.com/rs/zerolog"
	"strings"
)

// Use the same ContextKey type if defined elsewhere, or define it here
// Define the key used to store the userID string from the JWT claim in the context.
// The next middleware (AuthContextMiddleware) will pick this up.
const (
	UserIDStrContextKey ContextKey = "userID_str"
)

type JWTAuthMiddleware struct {
	iName      string
	jwtService services.IJWTService
	log        *zerolog.Logger // logger
}

// NewJWTAuthMiddleware creates the middleware instance.
func NewJWTAuthMiddleware(log *zerolog.Logger, jwtService services.IJWTService) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		iName:      "JWTAuthMiddleware",
		jwtService: jwtService,
		log:        log,
	}
}

// Authenticate verifies the JWT from the Authorization header and adds the userID to the context.
func (jam *JWTAuthMiddleware) Authenticate() fiber.Handler {
	const kName = "Authenticate"

	jam.log.Debug().Interface(kName, jam.iName).Msg("authenticating user")

	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return func(c *fiber.Ctx) error {
		jam.log.Debug().Interface(kName, jam.iName).Interface("params", c.AllParams()).Msg("authenticating request")

		// 1. Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			jam.log.Debug().Interface(kName, jam.iName).Msg("Authorization header missing")
			//http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Authorization header required"))
		}

		// 2. Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			jam.log.Debug().Interface(kName, jam.iName).Str("header", authHeader).Msg("Authorization header format must be Bearer {token}")
			//http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Authorization header requires Bearer-Token"))
		}
		tokenString := parts[1]

		// 3. Verify the access token using the JWTService
		token, err := jam.jwtService.VerifyAccessToken(tokenString)
		if err != nil {
			// Log the specific JWT validation error
			jam.log.Info().Interface(kName, jam.iName).Err(err).Msg("Invalid or expired access token")
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid or expired access token"))
		}

		// 4. Check if the token is valid and extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 5. Extract the UserID (subject 'sub' claim)
			userIDStr, ok := claims["sub"].(string)
			if !ok || userIDStr == "" {
				jam.log.Debug().Interface(kName, jam.iName).Interface("claims", claims).Msg("Invalid token: 'sub' claim is missing or not a string")
				jam.log.Error().Interface(kName, jam.iName).Err(err).Msg("Invalid access token claims")
				return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid token"))
			}

			// 6. Add the extracted userID string to the request context
			//ctx := context.WithValue(c.Context(), UserIDStrContextKey, userIDStr)
			c.Locals(UserIDStrContextKey, userIDStr) // Store the token in the context

			// 7. Call the next handler in the chain with the new context
			//next.ServeHTTP(w, r.WithContext(ctx))

			//jam.log.Debug().Interface(kName, jam.iName).Interface(kName, jam.iName).
			//	//Interface("claims", claims).
			//	Interface("params", c.AllParams()).
			//	Interface("context", c.Context().UserValue(UserIDStrContextKey)).
			//	Msg("dumped claims and context-key")

		} else {
			jam.log.Warn().Interface(kName, jam.iName).Interface("Authenticate", "JWTAuthMiddleware").Bool("tokenValid", token.Valid).Msg("Token claims invalid or token is not valid")
			//http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Token is invalid"))
		}
		return c.Next()
	}
}
