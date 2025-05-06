package main

import (
	"context"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

// logger Interface Name
const iName = "telko-moment-server"

func main() {
	//logger Key Name
	const kName = "main"

	// initialize logger
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load configuration
	log.Info().Interface(kName, iName).Msg("Loading configuration")
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Str("app", "MAIN").Msg("failed to load config file")
	}
	log.Info().Interface(kName, iName).Msg("Success: Loaded configuration")

	// Programmatically set Swagger info
	docs.SwaggerInfo.Title = "Tellme Messenger API"
	docs.SwaggerInfo.Description = "Codenamed telko-moment-server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost" + cfg.Server.Port // Dynamic for local dev
	docs.SwaggerInfo.BasePath = "/v1"                     // Adjust if needed
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Initialize MongoDB connection with timeout
	log.Info().Interface(kName, iName).Msg("initializing MongoDB connection")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Str("uri", cfg.MongoDB.URI).Msg("Failed to connect to MongoDB")
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Error().Err(err).Interface(kName, iName).Msg("Failed to disconnect from MongoDB") // Error, not Fatal
		}
	}()
	log.Info().Interface(kName, iName).Msg("Success: Established MongoDB connection")

	db := client.Database(cfg.MongoDB.Database)

	// Create initial indexes
	// ::: Indexes
	err = internalmongodb.CreateInitialIndexes(db, &log)
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Msg("Failed to create initial indexes")
		return
	}

	// Initialize dependencies
	encryptionSvc, err := pkgservices.NewAESEncryptionService(cfg.Encryption.AESKey, &log)
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Msg("Failed to create EncryptionService")
		return
	}

	jwtSvc, err := pkgservices.NewJWTService(&log, cfg.Jwt)
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Msg("Failed to create JWTService")
		return
	}

	keyHashSvc, err := pkgservices.NewHMACSearchKeyService(&log, cfg.Hashing.HMACSecretKey)
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Msg("Failed to create HMACSearchKeyService")
		return
	}

	// ::: Authorization
	// Initialize MongoDB adapter for auth
	authznAdapter, err := mongodbadapter.NewAdapter(cfg.MongoDB.URI + "/" + cfg.MongoDB.Database) // authorization collection is "casbin_rules"
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Msg("failed to initialize MongoDB adapter")
	}
	authznModelFilePath := "configs/casbin/abac_model.conf"
	authznSvc, err := services.NewCasbinAuthorizationService(&log, authznModelFilePath, authznAdapter)
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Msg("failed to initialize AuthorizationService")
	}

	// Load policies
	err = authznSvc.LoadPolicies()
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Msg("failed to load authorization policies")
	}

	// ::: Settings
	settingsRepo := mongodb.NewSettingsRepository(&log, db)
	settingsSvc := services.NewSettingsService(settingsRepo)
	settingsCtrl := controllers.NewSettingsController(&log, settingsSvc, authznSvc)

	// ::: Users
	userRepo := mongodb.NewUserRepository(&log, db, encryptionSvc, keyHashSvc)
	userSvc := services.NewUserService(&log, userRepo)
	userCtrl := controllers.NewUserController(&log, userSvc, settingsSvc, authznSvc)

	// ::: Authentication
	authctRepo := mongodb.NewAuthenticationRepository(&log, db, encryptionSvc, keyHashSvc)
	authctSvc := services.NewAuthenticationService(&log, authctRepo)
	authctCtrl := controllers.NewAuthController(&log, userSvc, authctSvc, settingsSvc, jwtSvc)

	// ::: Messages
	msgRepo := mongodb.NewMessageRepository(&log, db)
	msgSvc := services.NewMessageService(msgRepo)
	msgCtrl := controllers.NewMessageController(&log, msgSvc)

	// ::: Middleware
	authctMdw := middleware.NewJWTAuthMiddleware(&log, jwtSvc)
	authCtxMdw := middleware.NewAuthContextMiddleware(&log, userRepo)

	// Setup Fiber app
	app := fiber.New()

	// :::: add middleware
	// Logging remote IP and Port
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n\n",
	}))
	// Idempotency
	//app.Use(idempotency.New())	//TODO use later, for now disable any interfering on requests middleware for testing & development speed

	// ::: health & Metrics
	//app.Use(healthcheck.New())	// Provide a minimal config
	//app.Get("/metrics", monitor.New())	// Initialize default config (Assign the middleware to /metrics)

	// Setup routes
	routesHandler := handlers.NewRoutesHandler(&log, authctMdw, authCtxMdw, userCtrl, settingsCtrl, authctCtrl, msgCtrl)
	routesHandler.SetupRoutes(app) // layered

	// handle swagger routes
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start server with graceful shutdown
	log.Info().Str("port", cfg.Server.Port).Msg("Starting server")
	func() {
		if err := app.Listen(`:` + cfg.Server.Port); err != nil {
			log.Fatal().Err(err).Interface(kName, iName).Msg("Failed to listen")
		}
	}()

	// Wait for interrupt signal (e.g., Ctrl+C)
	err = app.ShutdownWithContext(context.Background())
	if err != nil {
		log.Fatal().Err(err).Interface(kName, iName).Msg("Failed to shutdown")
		return
	}
	log.Info().Interface(kName, iName).Msg("Server stopped")

}
