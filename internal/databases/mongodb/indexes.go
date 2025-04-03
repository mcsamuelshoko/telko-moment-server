package mongodb

import (
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateInitialIndexes creates initial indexes for all collections.
func CreateInitialIndexes(db *mongo.Database, log *zerolog.Logger) error {
	if err := createIndexesForUsers(db, log); err != nil {
		log.Error().Str("indexes", "CreateInitialIndexes").Err(err).Msg("Failed to create indexes for users collection")
		return err
	}
	if err := createIndexesForAuthentication(db, log); err != nil {
		log.Error().Str("indexes", "CreateInitialIndexes").Err(err).Msg("Failed to create indexes for authentication collection")
		return err
	}
	if err := createIndexesForSettings(db, log); err != nil {
		log.Error().Str("indexes", "CreateInitialIndexes").Err(err).Msg("Failed to create indexes for settings collection")
		return err
	}
	return nil
}

func createIndexesForUsers(db *mongo.Database, log *zerolog.Logger) error {
	const loggerFunctionName = "createIndexesForUsers"
	u := models.User{}
	err := u.CreateUniqueIndexes(db)
	if err != nil {
		log.Error().Str("indexes", loggerFunctionName).Err(err).Msg("Failed to create indexes for users collection")
		return err
	}
	log.Info().Str("indexes", loggerFunctionName).Msg("Created indexes for users collection")
	return nil
}

func createIndexesForAuthentication(db *mongo.Database, log *zerolog.Logger) error {
	const loggerFunctionName = "createIndexesForAuthentication"
	a := models.Authentication{}
	err := a.CreateUniqueIndexes(db)
	if err != nil {
		log.Error().Str("indexes", loggerFunctionName).Err(err).Msg("Failed to create indexes for authentication collection")
		return err
	}
	log.Info().Str("indexes", loggerFunctionName).Msg("Created indexes for authentication collection")
	return nil
}

func createIndexesForSettings(db *mongo.Database, log *zerolog.Logger) error {
	const loggerFunctionName = "createIndexesForSettings"
	s := models.Settings{}
	err := s.CreateUniqueIndexes(db)
	if err != nil {
		log.Error().Str("indexes", loggerFunctionName).Err(err).Msg("Failed to create indexes for Settings collection")
		return err
	}
	log.Info().Str("indexes", loggerFunctionName).Msg("Created indexes for Settings collection")
	return nil
}

// Add more index creation functions for other collections (e.g., Settings, Authentication)
