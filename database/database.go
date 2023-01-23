package database

import (
	"fmt"
	"github.com/Yash294/TODO/util"
	"github.com/Yash294/TODO/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {

	config, err := util.LoadConfig(".")
    if err != nil {
        panic("Cannot load config")
    }

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.Host, config.DBUser, config.DBPass, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database.")
	}

	DB.AutoMigrate(&models.Task{})
	DB.AutoMigrate(&models.User{})

	return nil

}
