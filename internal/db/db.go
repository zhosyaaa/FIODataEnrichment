package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var (
	DB *gorm.DB
)

func InitDatabase() {
	dbConnString := os.Getenv("DB_CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		return
	}
	DB = db
	log.Info().Msg("Connected to the database")
}
