package main

import (
	"context"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/configs"
	"github.com/mcsamuelshoko/telko-moment-server/docs"
	"github.com/mcsamuelshoko/telko-moment-server/internal/controllers"
	internalmongodb "github.com/mcsamuelshoko/telko-moment-server/internal/databases/mongodb"
	"github.com/mcsamuelshoko/telko-moment-server/internal/handlers"
	"github.com/mcsamuelshoko/telko-moment-server/internal/handlers/middleware"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository/mongodb"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	pkgservices "github.com/mcsamuelshoko/telko-moment-server/pkg/services"
	"github.com/rs/zerolog"
	"github.com/swaggo/fiber-swagger" // fiber-swagger middleware
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {

	// initialize logger
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load configuration
	log.Info().Str("app", "MAIN").Msg("Loading configuration")
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Str("app", "MAIN").Msg("failed to load config file")
	}
	log.Info().Str("app", "MAIN").Msg("Success: Loaded configuration")
	//b, err := json.Marshal(cfg)
	//if err != nil {
	//	log.Fatal().Err(err).Str("app", "MAIN").Msg("failed to marshal config file")
	//}
	// After loading the config
	//log.Debug().
	//	Str("mongodb_uri", cfg.MongoDB.URI).
	//	Str("mongodb_db", cfg.MongoDB.Database).
	//	Str("jwt_refresh_secret", cfg.Jwt.RefreshTokenSecret).
	//	Str("enc_key", cfg.Encryption.AESKey).
	//Str("configs", string(b)).
	//Msg("Loaded configuration")

	// Programmatically set Swagger info
	docs.SwaggerInfo.Title = "Tellme Messenger API"
	docs.SwaggerInfo.Description = "Codenamed telko-moment-server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost" + cfg.Server.Port // Dynamic for local dev
	docs.SwaggerInfo.BasePath = "/v1"                     // Adjust if needed
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Initialize MongoDB connection with timeout
	log.Info().Str("app", "MAIN").Msg("initializing MongoDB connection")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatal().Err(err).Str("app", "MAIN").Str("uri", cfg.MongoDB.URI).Msg("Failed to connect to MongoDB")
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Error().Err(err).Str("app", "MAIN").Msg("Failed to disconnect from MongoDB") // Error, not Fatal
		}
	}()
	log.Info().Str("app", "MAIN").Msg("Success: Established MongoDB connection")

	db := client.Database(cfg.MongoDB.Database)

	// Create initial indexes
	err = internalmongodb.CreateInitialIndexes(db, &log)
	if err != nil {
		log.Fatal().Err(err).Str("app", "MAIN").Msg("Failed to create initial indexes")
		return
	}

	// Initialize dependencies
	encryptionSvc, err := pkgservices.NewAESEncryptionService(cfg.Encryption.AESKey, &log)
	if err != nil {
		log.Fatal().Err(err).Str("app", "MAIN").Msg("Failed to create EncryptionService")
		return
	}

	jwtSvc, err := pkgservices.NewJWTService(&log, cfg.Jwt)
	if err != nil {
		log.Fatal().Err(err).Str("app", "MAIN").Msg("Failed to create JWTService")
		return
	}

	keyHashSvc, err := pkgservices.NewHMACSearchKeyService(&log, cfg.Hashing.HMACSecretKey)
	if err != nil {
		log.Fatal().Err(err).Interface("main", "Server").Msg("Failed to create HMACSearchKeyService")
		return
	}

	settingsRepo := mongodb.NewSettingsRepository(&log, db)
	settingsSvc := services.NewSettingsService(settingsRepo)
	settingsCtrl := controllers.NewSettingsController(&log, settingsSvc)

	userRepo := mongodb.NewUserRepository(&log, db, encryptionSvc, keyHashSvc)
	userSvc := services.NewUserService(&log, userRepo)
	userCtrl := controllers.NewUserController(&log, userSvc, settingsSvc)

	authctRepo := mongodb.NewAuthenticationRepository(&log, db, encryptionSvc, keyHashSvc)
	authctSvc := services.NewAuthenticationService(&log, authctRepo)
	authctCtrl := controllers.NewAuthController(&log, userSvc, authctSvc, settingsSvc, jwtSvc)

	// Initialize MongoDB adapter for auth
	authznAdapter, err := mongodbadapter.NewAdapter(cfg.MongoDB.URI + "/" + cfg.MongoDB.Database) // authorization collection is "casbin_rules"
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize MongoDB adapter")
	}
	authznModelFilePath := "configs/casbin/abac_model.conf"
	authznSvc, err := services.NewAuthorizationService(&log, authznModelFilePath, authznAdapter)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize AuthorizationService")
	}
	// Load policies
	err = authznSvc.LoadPolicies()
	if err != nil {
		log.Fatal().Err(err).Interface("main", "server").Msg("failed to load authorization policies")
	}

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authznSvc, userRepo, settingsRepo)

	// Setup Fiber app
	app := fiber.New()

	// Setup routes
	routesHandler := handlers.NewRoutesHandler(&log, authMiddleware, userCtrl, settingsCtrl, authctCtrl)
	routesHandler.SetupRoutes(app) // layered

	// handle swagger routes
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start server with graceful shutdown
	log.Info().Str("port", cfg.Server.Port).Msg("Starting server")
	func() {
		if err := app.Listen(`:` + cfg.Server.Port); err != nil {
			log.Fatal().Err(err).Msg("Failed to listen")
		}
	}()

	// Wait for interrupt signal (e.g., Ctrl+C)
	err = app.ShutdownWithContext(context.Background())
	if err != nil {
		log.Fatal().Err(err).Str("app", "MAIN").Msg("Failed to shutdown")
		return
	}
	log.Info().Str("app", "MAIN").Msg("Server stopped")

}
