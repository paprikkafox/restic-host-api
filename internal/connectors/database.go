package connectors

import (
	"restic-host-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB(path string) *gorm.DB {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to database!")
	}
	database.AutoMigrate(&models.Job{})

	return database
}
