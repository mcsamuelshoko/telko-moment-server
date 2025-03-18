package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/api"
	"github.com/mcsamuelshoko/telko-moment-server/configs"
	"github.com/mcsamuelshoko/telko-moment-server/internal/handlers"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository/mongodb"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file")
	}

	// Initialize MongoDB connection
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongodb URI: " + cfg.MongoDB.URI)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to disconnect from mongodb")
		}
	}(client, context.Background())

	db := client.Database(cfg.MongoDB.Database)

	// Initialize dependencies
	userRepo := mongodb.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Setup Fiber app
	app := fiber.New()

	// Setup routes
	api.SetupRoutes(app, userHandler)

	// Start server
	err = app.Listen(cfg.Server.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
}

//
//import (
//	"github.com/gofiber/fiber/v2"
//)
//
//func main() {
//	app := fiber.New()
//
//	app.Get("/", func(c *fiber.Ctx) error {
//		return c.SendString("Hello, World!")
//	})
//
//	app.Listen(":8080")
//}
