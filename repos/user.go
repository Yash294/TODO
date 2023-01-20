package repos

import (
	"errors"
	"fmt"

	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/models"
)

func CreateUser(username string, password string) error {

	user := models.User{Username: username, Password: password}
	result := database.DB.Where("username = ?", username).Find(&user)

	if result.RowsAffected != 0 {
		// username already taken
		return errors.New("this username has already been taken")
	}

	result = database.DB.Create(&user)

	if result.Error != nil {
		return errors.New("failed to create requested user")
	}

	return nil
}

func ChangePassword(username string, currentPassword string, newPassword string) error {

	user := models.User{Username: username, Password: currentPassword}
	result := database.DB.Where("username = ? AND password = ?", username, currentPassword).Find(&user)

	if result.RowsAffected == 0 {
		panic("Username/Password is incorrect.")
	}

	result = database.DB.Model(&user).Update("password", newPassword)

	if result.Error != nil {
		panic("Failed to update user password.")
	}

	return nil
}

func GetLogin(username string, password string) error {
	return nil
}

func GetAllUsernames() ([]string, error) {
	var usernames []string
	result := database.DB.Model(models.User{}).Select("username").Find(&usernames)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return usernames, nil
}
