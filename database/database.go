package database

import (
	"github.com/Yash294/TODO/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {

	var err error

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database.")
	}

	DB.AutoMigrate(&models.Task{})
	DB.AutoMigrate(&models.User{})

	return nil

}
