package util

import (
	"fmt"
	"github.com/Yash294/TODO/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {

	config, err := LoadDBConfig(".")
    if err != nil {
        panic("Cannot load config")
    }

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", config.DBHost, config.DBUser, config.DBPass, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database.")
	}

	DB.AutoMigrate(&models.Task{})
	DB.AutoMigrate(&models.User{})

	return nil

}
