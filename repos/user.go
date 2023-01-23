package repos

import (
	"fmt"
	"errors"
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/models"
)

func CreateUser(username string, password string) error {

	var user models.User
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

	var user models.User
	result := database.DB.Where("username = ? AND password = ?", username, currentPassword).Find(&user)

	if result.RowsAffected == 0 {
		return errors.New("username/password is incorrect")
	}

	result = database.DB.Model(&user).Update("password", newPassword)

	if result.Error != nil {
		return errors.New("failed to update user password")
	}

	return nil
}

func Login(username string, password string) (bool, error) {
	result := database.DB.Model(models.User{}).Where("username = ? AND password = ?", username, password)

	if result.Error != nil {
		return false, errors.New("failed to retrieve login info")
	}

	// have to fix next
	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func CheckIfUsernameExists(username string) (bool, error) {

	fmt.Println(username)
	var query string
	result := database.DB.Model(models.User{}).Select("username").Where("username = ?", username).Find(&query)

	fmt.Println("QUERY: ", query)

	if result.Error != nil {
		return false, errors.New("failed to retrieve username")
	}

	if username == query {
		return false, nil
	}

	return true, nil
}
