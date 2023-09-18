package db

import (
	"fmt"
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
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	DB = db
}
